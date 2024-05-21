package data

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/gabrielmoura/davServer/config"
)

func deserializeI2pCfg(data []byte) (*config.I2PCfg, error) {
	var i2pCfg *config.I2PCfg
	err := json.Unmarshal(data, &i2pCfg)
	return i2pCfg, err
}
func serializeI2pCfg(i2pCfg *config.I2PCfg) ([]byte, error) {
	return json.Marshal(i2pCfg)
}

func writeI2pConfig(db *badger.DB, i2pCfg *config.I2PCfg) error {
	return db.Update(func(txn *badger.Txn) error {
		data, err := serializeI2pCfg(i2pCfg)
		if err != nil {
			return err
		}
		return txn.Set([]byte("i2p"), data)
	})
}
func SaveI2pConfig(i2pCfg *config.I2PCfg) error {
	return writeI2pConfig(dB, i2pCfg)
}
func readI2pConfig(db *badger.DB) (*config.I2PCfg, error) {
	var i2pCfg *config.I2PCfg
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("i2p"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			i2pCfg, err = deserializeI2pCfg(val)
			return err
		})
	})
	return i2pCfg, err
}
