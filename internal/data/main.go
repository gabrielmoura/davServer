package data

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/gabrielmoura/davServer/config"
)

var dB *badger.DB

func readPublicFiles(db *badger.DB) ([]PubFile, error) {
	var files []PubFile
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("files"))
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
func readUsers(db *badger.DB) ([]User, error) {
	var user []User
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("users"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			user, err = deserialize(val)
			return err
		})
	})
	return user, err
}

func writePublicFiles(db *badger.DB, files []PubFile) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := serializeFiles(files)
		if err != nil {
			return err
		}
		return txn.Set([]byte("files"), data)
	})
}
func writeUsers(db *badger.DB, user []User) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := serialize(user)
		if err != nil {
			return err
		}
		return txn.Set([]byte("users"), data)
	})
}
func serialize(user []User) ([]byte, error) {
	return json.Marshal(user)
}
func serializeFiles(files []PubFile) ([]byte, error) {
	return json.Marshal(files)
}

func deserialize(data []byte) ([]User, error) {
	var user []User
	err := json.Unmarshal(data, &user)
	return user, err
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
