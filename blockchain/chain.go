package blockchain

import (
	"sync"
)

/*
[블록체인 기본 원리]
#1_block's Hash = (Data+"")의 해쉬
#2_block's Hash = (Data + #1_block의 해쉬)의 해쉬.......

[이 프로젝트에서 사용할 패턴 : Singleton]
- 우리의 어플리케이션 내에서 언제든지 "blockchain"의 인스턴스 단 하나만을 사용하여 작업하는 것 > init single instance only
*/



type blockChain struct {
	NewestHash 	string 	`json:"newestHash"`
	Height 		int		`json:"height"`
}

// Singleton 패턴으로 이번 프로젝트에서 사용할 "blockchain" 인스턴스
// 오직 blockchain 패키지 내부에서만 이 변수에 접근하도록 함 (Go에서 Public한 변수를 만들려면, 변수 이름이 대문자로 시작해야함)
var b *blockChain
var once sync.Once

func (b *blockChain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height)	// 블록을 만들고
	b.NewestHash = block.Hash							// 체인을 업데이트
	b.Height = block.Height
}

func BlockChain() *blockChain {
	// 싱글톤 패턴으로 블록체인을 다룬다. -> b 라는 블록체인 인스턴스를 하나 생성해놓고 요 객체 하나로 모든 블록체인에 대한 내용을 다루는 것 (블록체인 인스턴스는 한개만 생성됨)
	// 우리가 이번 프로젝트에서 사용할 블록체인 인스턴스의 포인터가 존재하는지 확인
	if b == nil {
		// b 가 비어있으면, "blockChain 구조의 포인터"를 변수 "b"에 초기화
		// !! 만약 이 블록체인 어플리케이션이 병렬적으로 실행된다면, "첫 실행"에서만 초기화 되어야하는 이 b 변수가 여러번 초기화될 수 있다. -> sync 패키지 사용
		once.Do(func() {
			// 첫번째 블록
			b = &blockChain{"", 0}
			b.AddBlock("Genesis Block")
		})
		//b = &blockChain{}
		//b.blocks = append(b.blocks, createBlock("Genesis Block")) // 이 기능이 통틀어 한번만 호출되도록 함
	}
	return b
}