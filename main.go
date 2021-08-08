package main

import (
	"github.com/ddhyun93/seancoin/explorer"
	"github.com/ddhyun93/seancoin/rest"
)

func main() {
	go explorer.Start(3000)		// 얘가 실행되면 무한 루프 상태이기 때문에, restAPI가 실행될 여지가없음. 따라서 고루틴으로 실행
	rest.Start(4000)
}
