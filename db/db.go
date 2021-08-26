package db

import (
	"github.com/boltdb/bolt"
	"github.com/ddhyun93/seancoin/utils"
)

const (
	dbName = "blockchain.db"
	dataBucket = "data"
	blocksBucket = "blocks"
	checkpoint = "checkpoint"
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
			_, err := tx.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = tx.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			return err
		})
		utils.HandleErr(err)
	}
	return db
}

func Close() {
	DB().Close()
}

func SaveBlock(hash string, data []byte) { 		// key: 블록의 Hash value: 블록의 바이트 값
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))	// 버켓 인스턴스 불러오고
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

func SaveBlockChain(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

func Checkpoint() []byte {
	// 블록체인읨 체크포인트를 바이트로 리턴함
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

func Block(hash string) []byte {
	// 해쉬값에 맞는 Block 데이터를 전달함
	var data []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})
	return data
}