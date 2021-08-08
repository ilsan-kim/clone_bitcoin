package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

/*
[블록체인 기본 원리]
#1_block's Hash = (Data+"")의 해쉬
#2_block's Hash = (Data + #1_block의 해쉬)의 해쉬.......

[이 프로젝트에서 사용할 패턴 : Singleton]
- 우리의 어플리케이션 내에서 언제든지 "blockchain"의 인스턴스 단 하나만을 사용하여 작업하는 것 > init single instance only
*/

type Block struct {
	Data     	string	`json:"data"`
	Hash		string	`json:"hash"`
	PrevHash 	string	`json:"prev_hash,omitempty"`
	Height	 	int		`json:"height"`
}

type blockChain struct {
	blocks []*Block
}

// Singleton 패턴으로 이번 프로젝트에서 사용할 "blockchain" 인스턴스
// 오직 blockchain 패키지 내부에서만 이 변수에 접근하도록 함 (Go에서 Public한 변수를 만들려면, 변수 이름이 대문자로 시작해야함)
var b *blockChain
var once sync.Once
var ErrNotFound = errors.New("block not found")

func (b *Block) calcHash() {
	// 데이터와 이전 해쉬를 기반으로 내 해쉬값을 만들어냄 > Go 에 내장된 SHA256 패키지 사용 // string을 []bytes로 바꿔줘야 해싱 가능
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	// 빗코나 이더 표준인 16진수로 변환
	b.Hash = fmt.Sprintf("%x", hash)

}

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)
	// 블록체인의 길이가 0이면, 이는 첫번째 블록으로 직전 블록의 해쉬가 존재하지 않음
	if totalBlocks == 0 {
		return ""
	}
	// 블록체인의 길이가 0이 아니면 직전 블록의 해쉬가 존재함
	return GetBlockChain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockChain().blocks) + 1}
	newBlock.calcHash()
	return &newBlock
}

func (b *blockChain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockChain() *blockChain {
	// 우리가 이번 프로젝트에서 사용할 블록체인 인스턴스의 포인터가 존재하는지 확인
	if b == nil {
		// b 가 비어있으면, "blockChain 구조의 포인터"를 변수 "b"에 초기화
		// !! 만약 이 블록체인 어플리케이션이 병렬적으로 실행된다면, "첫 실행"에서만 초기화 되어야하는 이 b 변수가 여러번 초기화될 수 있다. -> sync 패키지 사용
		once.Do(func() {
			b = &blockChain{}
			b.AddBlock("Genesis Block")
		})
		//b = &blockChain{}
		//b.blocks = append(b.blocks, createBlock("Genesis Block")) // 이 기능이 통틀어 한번만 호출되도록 함
	}
	return b
}

func (b *blockChain) AllBlocks() []*Block {
	return b.blocks
}

func (b * blockChain) GetBlock(height int) (*Block, error) {
	if height > len(b.blocks) {
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}

