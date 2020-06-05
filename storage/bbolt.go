package storage

import (
	"encoding/json"

	"go.etcd.io/bbolt"
)

//PassageBBoltStorage is an implementation for Passage storage interface.
type PassageBBoltStorage struct {
	db         *bbolt.DB
	bucketName []byte
}

//Set PassageEntry for given Name.
func (p *PassageBBoltStorage) Set(entry *PassageEntry) error {
	return p.db.Update(func(tx *bbolt.Tx) error {
		bs, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		bucket := tx.Bucket(p.bucketName)
		if bucket == nil {
			bucket, err = tx.CreateBucketIfNotExists(p.bucketName)
			if err != nil {
				return err
			}
		}
		err = bucket.Put([]byte(entry.Name), bs)
		if err != nil {
			return err
		}
		return nil
	})
}

//Get name from BBolt storage
func (p *PassageBBoltStorage) Get(name string) (*PassageEntry, error) {
	var entry *PassageEntry
	err := p.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(p.bucketName)
		if err != nil {
			return err
		}
		payload := bucket.Get([]byte(name))
		err = json.Unmarshal(payload, entry)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entry, nil
}

const (
	_DefaultDB     = "passage.db"
	_DefaultBucket = "passwds"
)

//Create a new PassageBBoltStorage
func NewPassageBBoltStorage(dbPath string, bucketName string) (*PassageBBoltStorage, error) {
	if dbPath == "" {
		dbPath = _DefaultDB
	}
	if bucketName == "" {
		bucketName = _DefaultBucket
	}
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &PassageBBoltStorage{
		db,
		[]byte(_DefaultBucket),
	}, nil
}
