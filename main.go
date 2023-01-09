package main

import (
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v3"
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("./tmp/badger"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := write(db, "hello", "world"); err != nil {
		log.Fatal("failed to write:", err)
	}
	res, err := read(db, "hello")
	if err != nil {
		log.Fatal("failed to read:", err)
	}
	fmt.Println(string(res))
}

func write(db *badger.DB, key, value string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
}

func read(db *badger.DB, key string) (res string, err error) {
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			// WARN: Do not do this
			// res = val
			//
			// Do this instead.
			// res = append([]byte{}, val...)
			//
			// Or just turn this into a string.
			res = string(val)
			return nil
		})

		return nil
	})

	return
}
