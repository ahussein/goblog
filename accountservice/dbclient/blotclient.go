package dbclient

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/ahussein/goblog/accountservice/model"
	"github.com/boltdb/bolt"
)

type BoltClient struct {
	boltDB *bolt.DB
}

func (bc *BoltClient) OpenDB() {
	var err error
	bc.boltDB, err = bolt.Open("accounts.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (bc *BoltClient) QueryAccount(accountId string) (model.Account, error) {
	account := model.Account{}

	err := bc.boltDB.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("AccountBucket"))
		accountBytes := b.Get([]byte(accountId))
		if accountBytes == nil {
			return fmt.Errorf("No account found for ID: %s", accountId)
		}
		json.Unmarshal(accountBytes, &account)
		return nil
	})
	return account, err
}

func (bc *BoltClient) Seed() {
	bc.initializeBucket()
	bc.seedAccounts()
}

func (bc *BoltClient) initializeBucket() {
	bc.boltDB.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucket([]byte("AccountBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
}

func (bc *BoltClient) seedAccounts() {
	total := 100
	for i := 0; i < total; i++ {
		key := strconv.Itoa(1000 + i)
		acc := model.Account{
			Id:   key,
			Name: "Person_" + strconv.Itoa(i),
		}

		jsonBytes, _ := json.Marshal(acc)

		bc.boltDB.Update(func(t *bolt.Tx) error {
			b := t.Bucket([]byte("AccountBucket"))
			err := b.Put([]byte(key), jsonBytes)
			return err
		})
	}
	fmt.Printf("Seeded %v fake accounts....\n", total)
}
