package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/ddhyun93/seancoin/db"
	"github.com/ddhyun93/seancoin/utils"
)

type Block struct {
	Data     	string	`json:"data"`
	Hash		string	`json:"hash"`
	PrevHash 	string	`json:"prev_hash,omitempty"`
	Height	 	int		`json:"height"`
}

func (b *Block) toBytes() []byte {
	var blockBuffer bytes.Buffer				// 비어있는 변수 만듬
	encoder := gob.NewEncoder(&blockBuffer)		// blockBuffer 에 인코딩 한 값을 담을 수 있게 하여 이걸 encoder로 명명한 뒤
	utils.HandleErr(encoder.Encode(b))			// block 데이터를 인코딩하고 에러 검증
	return blockBuffer.Bytes()
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, b.toBytes())
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data: 		data,
		Hash: 		"",
		PrevHash: 	prevHash,
		Height: 	height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)	// int 인 Height 를 string 으로 변환
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}