package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// interface 타입을 인자로 받는다. -> 어떤 인자든 받도록 한다.
func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer				// 비어있는 변수 만듬
	encoder := gob.NewEncoder(&aBuffer)		// blockBuffer 에 인코딩 한 값을 담을 수 있게 하여 이걸 encoder로 명명한 뒤
	HandleErr(encoder.Encode(i))			// block 데이터를 인코딩하고 에러 검증
	return aBuffer.Bytes()
}

// 바이트를 받아 블록이나 블록체인의 형태로 디코딩해준다. (블록의 포인터나 블록체인의 포인터)
func FromBytes(i interface{}, data[]byte) {
	encoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(encoder.Decode(i))
}