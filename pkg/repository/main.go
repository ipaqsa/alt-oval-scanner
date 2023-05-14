package repository

import (
	"go.etcd.io/bbolt"
	"log"
)

func NewRepository(path string) *Repository {
	return &Repository{Path: path}
}

func (r *Repository) Open() error {
	log.Printf("open oval db: %s\n", r.Path)
	db, err := bbolt.Open(r.Path, 0666, nil)
	if err != nil {
		return err
	}
	r.Pointer = db
	return nil
}

func (r *Repository) Close() error {
	return r.Pointer.Close()
}

func (r *Repository) CreateBucket(bucket string) error {
	err := r.Pointer.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *Repository) Save(bucket string, key, data []byte) error {
	err := r.Pointer.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, data)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *Repository) Get(bucket string, key []byte) ([]byte, error) {
	var data []byte
	err := r.Pointer.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data = b.Get(key)
		return nil
	})
	return data, err
}
