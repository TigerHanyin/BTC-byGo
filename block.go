package main

import (
	"bytes"
	"crypto/sha256"
	"time"
)

//定义区块结构
//第一阶段：先实现基础字段：前区块哈希 哈希 数据
//补充字段：Version 时间戳 难度值
type Block struct {
	Version uint64

	PreHash    []byte
	MerkleRoot []byte
	TimeStamp  uint64
	Bits       uint64
	Nonce      uint64
	Hash       []byte
	//数据
	Data []byte
}

//创建一个区块的方法
func NewBlock(data string, prevHash []byte) *Block {
	b := Block{
		Version:    0,
		PreHash:    prevHash,
		MerkleRoot: nil,
		TimeStamp:  uint64(time.Now().Unix()),
		Bits:       0, //随意写
		Nonce:      0, //同上
		Hash:       nil,
		Data:       []byte(data),
	}
	//计算哈希值
	b.setHash()

	//todo
	return &b
}

////计算哈希值方法
func (b *Block) setHash() {
	tmp := [][]byte{
		uintToByte(b.Version),
		b.MerkleRoot,
		uintToByte(b.TimeStamp),
		uintToByte( b.Bits),
		uintToByte(b.Nonce) ,
		b.PreHash,
		b.Hash,
		b.Data,
	}
	//使用join 方法 将二维切片转化为1 维切片
	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	b.Hash = hash[:]
}
