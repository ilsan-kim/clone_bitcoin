package db

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/ddhyun93/seancoin/utils"
)

const (
	dbName = "blockchain.db"
	daataBucket = "data"
	blocksBucket = "blocks"
)
var db *bolt.DB


/*
bolt db는 단순히 key/value 를 저장하는 db이다.
따라서 정리나 정렬 등의 기능이 없다.

우리는 블록체인을 이 bolt db에 저장할 예정인데, 문제는 순서를 모른다는 것이고, 순서대로 db에 기록할 수도 없다는 것이다.
--> 그래서 "마지막 생성된 블록이 어떤건지"만 기록하는 bucket (테이블개념)을 하나 만들어서 prev prev prev 로 찾아갈 계획이다.
 */

func DB() *bolt.DB {
	if db == nil {
		// if db is not exist
		// initialize db
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		// pointing db
		db = dbPointer
		utils.HandleErr(err)
		err = db.Update(func(tx *bolt.Tx) error {
			// 버킷이 없으면 생성
			_, err := tx.CreateBucketIfNotExists([]byte(daataBucket))
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func SaveBlock(hash string, data []byte) { 		// key: 블록의 Hash value: 블록의 바이트 값
	fmt.Printf("Saving Block %s\nData: %b", hash, data)
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))	// 버켓 인스턴스 불러오고
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockChain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte("checkpoint"), data)
		return err
	})
	utils.HandleErr(err)
}