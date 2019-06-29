package main

import "fmt"

func main() {
	//初始化区块链
	bc:=NewBlockChain()
	bc.AddBlock("26号btc暴涨20%")
	bc.AddBlock("27号btc暴涨10%")
	//变量区块数据
	for i,block:=range bc.Blocks{
		fmt.Printf("\n+++++++++ 当前区块高度: %d ++++++++++\n", i)
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevHash : %x\n", block.PreHash)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Bits : %d\n", block.Bits)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("Data : %s\n", block.Data)
	}
}

