package etcdv3

import (
	"context"
	"crypto/sha1"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"sync"
	"time"
)

type registry struct {
	sync.RWMutex
	client 	  *Client
	leases    map[string]clientv3.LeaseID
	register  map[string]int
	timeout   time.Duration
}

func (e *registry) Register(key , value string, ttl int64) error {
	e.Lock()
	leaseID, ok := e.leases[key]
	e.Unlock()

	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(),  e.timeout)
		defer cancel()

		rsp, err := e.client.Get(ctx, key, clientv3.WithSerializable())
		if err != nil {
			return err
		}

		for _, kv := range rsp.Kvs {
			if kv.Lease > 0 {
				leaseID = clientv3.LeaseID(kv.Lease)

				hash := sha1.New()
				h, err := hash.Write(kv.Value)
				if err != nil {
					continue
				}

				e.Lock()
				e.leases[key] = leaseID
				e.register[key] = h
				e.Unlock()
				break
			}
		}
	}

	var leaseNotFound bool
	if leaseID > 0 {
		if _, err := e.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
			if err != rpctypes.ErrLeaseNotFound {
				return err
			}
			leaseNotFound = true
		}
	}

	hash := sha1.New()
	h, err := hash.Write([]byte(value))
	if err != nil {
		return err
	}

	v, ok := e.register[key]
	if ok && v == h && !leaseNotFound {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	var lgr *clientv3.LeaseGrantResponse
	if ttl > 0 {
		lgr, err = e.client.Grant(ctx, ttl)
		if err != nil {
			return err
		}
	}
	var putOpts []clientv3.OpOption
	if lgr != nil {
		putOpts = append(putOpts, clientv3.WithLease(lgr.ID))
	}
	if _, err = e.client.Put(ctx, key, value, putOpts...); err != nil {
		return err
	}

	e.Lock()
	e.register[key] = h
	if lgr != nil {
		e.leases[key] = lgr.ID
	}
	e.Unlock()

	return nil
}

func (e *registry) Deregister(key string) error {
	e.Lock()
	delete(e.register, key)
	delete(e.leases, key)
	e.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	if _, err := e.client.Delete(ctx, key); err != nil {
		return err
	}
	return nil
}