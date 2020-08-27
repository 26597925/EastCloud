package buntdb

import (
	"context"
	"github.com/tidwall/buntdb"
	"os"
	"path/filepath"
	"time"
)

type Store struct {
	db *buntdb.DB
}

func NewStore(path string) (*Store, error) {
	if path != ":memory:" {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	var val string
	err := s.db.View(func(tx *buntdb.Tx) error {
		var err error
		val, err = tx.Get(key)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}

		return nil
	})

	return val, err
}

func (s *Store) Set(ctx context.Context, key string, value string, expire time.Duration) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		opts := &buntdb.SetOptions{Expires: expire != 0, TTL: expire}
		_, _, err := tx.Set(key, value, opts)
		return err
	})
}

func (s *Store) Delete(ctx context.Context, key string) error {
	return s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		return nil
	})
}

func (s *Store) Keys() ([]string, error) {
	keys := make([]string, 0)
	err := s.db.View(func(tx *buntdb.Tx) error {

		err := tx.Ascend("", func(key, value string) bool {
			//fmt.Printf("key: %s, value: %s\n", key, value)
			keys = append(keys, key)
			return true
		})
		return err
	})

	return keys, err
}

func (s *Store) Close() error{
	return s.db.Close()
}

func (s *Store) Type() string{
	return "buntdb"
}