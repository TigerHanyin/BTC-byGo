package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"time"
)

type Transaction struct {
	TXID []byte //交易id
	TXInputs []TXIput//可以有多个输入
	TXOutputs []TXOutput//可以多个输出
	TimeStamp uint64
}
type TXIput struct {
	Txid []byte//这个input所引用output所在交易ID
	Index int64//在上一笔中交易索引
	Scriptsig string//付款人对当前交易的签名
}
type TXOutput struct {
	ScriptPubk string//收款人的公钥。先理解为地址
	Value float64 //转账金额
}
//获取交易ID哈希运算
func (tx *Transaction)setHash(){
	//gob编码tx得到字节流，做sha256，赋值给TXID
	var buf bytes.Buffer
	encoder:=gob.NewEncoder(&buf)
	err:=encoder.Encode(tx)
	if err!=nil{
		log.Fatal(err)
		return
	}
	hash:=sha256.Sum256(buf.Bytes())
	//交易节流做哈希 值作为id
	tx.TXID=hash[:]
}
//创建挖矿交易
var reward =12.5
func NewCoinbaseTx(miner /*挖矿人*/ ,data string) *Transaction{
	//没有输入,只有一个输出 挖矿奖励
	//需要辨识度，没有input所以不需要签名，给特点
	//挖矿交易不需要签名，所以挖矿字段可以写任意值，只有矿工有权利写，
	//钟本聪写的创世语 现在矿池写自己名字
	input:=TXIput{
		Txid:nil,
		Index:-1,
		Scriptsig:data,
	}
	output:=TXOutput{Value:reward,ScriptPubk:miner}
	timeStamp:=time.Now().Unix()
	tx:=Transaction{
		TXID:nil,
		TXInputs:[]TXIput{input},
		TXOutputs:[]TXOutput{output},
		TimeStamp:uint64( timeStamp),
	}
	return &tx
}
