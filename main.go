package main

import (
	"github.com/ddhyun93/seancoin/blockchain"
	"github.com/ddhyun93/seancoin/cli"
	"github.com/ddhyun93/seancoin/db"
)

func main() {
	defer db.Close()
	blockchain.BlockChain()
	cli.Start()
}
