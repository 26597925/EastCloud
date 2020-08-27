package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
)

type Client struct {
	opts *Options
	pool *pool
}

func newClient(options *Options) *Client {
	rc := &Client{
		opts: options,
	}

	rc.pool = newPool(options.PoolSize, options.PoolTTL, options.PoolMaxIdle, options.PoolMaxStreams)
	return rc
}

//path:/service.Struct/Method
func (c *Client) Call(ctx context.Context, path string, args, reply interface{}, opts ...grpc.CallOption) error {
	var grr error
	cc, err := c.pool.getConn(c.opts.address, c.opts.grpcDialOptions...)
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending request: %v", err))
	}
	defer func() {
		// defer execution of release
		c.pool.release(c.opts.address, cc, grr)
	}()

	ch := make(chan error, 1)

	go func() {
		err := cc.Invoke(ctx, path, args, reply, opts...)
		ch <- err
	}()

	select {
	case err := <-ch:
		grr = err
	case <-ctx.Done():
		grr = errors.New(fmt.Sprintf("Timeout: %v",  ctx.Err()))
	}

	return grr
}

func (c *Client) Stream(ctx context.Context, path string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	var grr error
	cc, err := c.pool.getConn(c.opts.address, c.opts.grpcDialOptions...)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error sending request: %v", err))
	}

	defer func() {
		// defer execution of release
		c.pool.release(c.opts.address, cc, grr)
	}()

	desc := &grpc.StreamDesc{
		StreamName:    path,
		ClientStreams: true,
		ServerStreams: true,
	}

	newCtx, cancel := context.WithCancel(ctx)
	st, err := cc.NewStream(newCtx, desc, path, opts...)
	if err != nil {
		cancel()
		cc.Close()
		return nil, err
	}

	return st, nil
}