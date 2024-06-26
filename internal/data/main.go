package data

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/gabrielmoura/davServer/config"
)

var dB *badger.DB

const IndexPubName = "files"

func readPublicFiles(db *badger.DB) ([]PubFile, error) {
	var files []PubFile
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(IndexPubName))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			files, err = deserializeFiles(val)
			return err
		})
	})
	return files, err
}

func writePublicFiles(db *badger.DB, files []PubFile) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := serializeFiles(files)
		if err != nil {
			return err
		}
		return txn.Set([]byte(IndexPubName), data)
	})
}

func serializeFiles(files []PubFile) ([]byte, error) {
	return json.Marshal(files)
}

func deserializeFiles(data []byte) ([]PubFile, error) {
	var files []PubFile
	err := json.Unmarshal(data, &files)
	return files, err
}

func InitDB() error {
	opts := badger.DefaultOptions(config.Conf.DBDir)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	dB = db
	return nil
}
