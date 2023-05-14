package repository

import "go.etcd.io/bbolt"

type Repository struct {
	Path    string
	Pointer *bbolt.DB
}
