package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"time"
	"fmt"
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"crypto/elliptic"
	"strings"
)

type Transaction struct {
	TXID      []byte     //交易id
	TXInputs  []TXIput   //可以有多个输入
	TXOutputs []TXOutput //可以多个输出
	TimeStamp uint64
}
type TXIput struct {
	Txid  []byte //这个input所引用output所在交易ID
	Index int64  //在上一笔中交易索引
	//Scriptsig string//付款人对当前交易的签名
	Scriptsig []byte //付款人对当前交易的签名
	PubKey    []byte//付款人的公钥
}
type TXOutput struct {
	ScriptPubKeyHash []byte  //收款人的公钥
	Value      float64 //转账金额
}
//由于没有办法直接将地址赋值给TXoutput，所以需要提供一个output的方法
func newTXOutput(address string,amount float64)TXOutput {
	output:=TXOutput{Value:amount}
	pubKeyHash:=getPubKeyHashFromAddress(address)
	output.ScriptPubKeyHash=pubKeyHash
	return output
}

//获取交易ID哈希运算
func (tx *Transaction) setHash() {
	//gob编码tx得到字节流，做sha256，赋值给TXID
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Fatal(err)
		return
	}
	hash := sha256.Sum256(buf.Bytes())
	//交易节流做哈希 值作为id
	tx.TXID = hash[:]
}

//创建挖矿交易
var reward = 12.5

func NewCoinbaseTx(miner /*挖矿人*/ , data string) *Transaction {
	//没有输入,只有一个输出 挖矿奖励
	//需要辨识度，没有input所以不需要签名，给特点
	//挖矿交易不需要签名，所以挖矿字段可以写任意值，只有矿工有权利写，
	//钟本聪写的创世语 现在矿池写自己名字
	input := TXIput{
		Txid:      nil,
		Index:     -1,
		Scriptsig: nil,
		PubKey:[]byte(data),
	}
	//output := TXOutput{Value: reward, ScriptPubk: miner}
	output:=newTXOutput(miner,reward)
	timeStamp := time.Now().Unix()
	tx := Transaction{
		TXID:      nil,
		TXInputs:  []TXIput{input},
		TXOutputs: []TXOutput{output},
		TimeStamp: uint64(timeStamp),
	}
	tx.setHash()
	return &tx
}

//创建普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//钱包就是在这里使用的，from=》钱包里面找到对应的wallet-》私钥-》签名
	wm:=NewWalletManager()
	if wm==nil{fmt.Println("打开钱包失败!")
		return nil}
		//钱包里面找到对应的wallet
		wallet,ok:=wm.Wallets[from]
		if !ok{
			fmt.Println("没有找到付款人对应的私钥")
			return nil
		}
		fmt.Println("找到付款人的私钥和公钥，准备创建交易...")
		pubkey:=wallet.PubKey
	pubKeyHash:=getPubKeyHashFromPubKey(pubkey)

	//from的utxo集合
	var spentUTXO = make(map[string][]int64)
	//utxo的sum
	var retValue float64
	//遍历账本找到可以使用的utxo 以及其中的钱
	spentUTXO, retValue = bc.FindNeedUTXO(pubKeyHash, amount)
	if retValue < amount {
		fmt.Println("金额不足，交易失败！")
		return nil
	}
	//正常交易
	var inputs []TXIput
	var outputs []TXOutput
	// 4. 拼接inputs
	// > 遍历utxo集合，每一个output都要转换为一个input(3)
	for txid, indexArray := range spentUTXO {
		//遍历下标, 注意value才是我们消耗的output的下标
		for _, i := range indexArray {
			input := TXIput{[]byte(txid), i, nil,pubkey}
			inputs = append(inputs, input)
		} //for
	} //for
	// 5. 拼接outputs
	// > 创建一个属于to的output
	//创建给收款人的output
	output1 := newTXOutput(to, amount)
	outputs = append(outputs, output1)
	// > 如果总金额大于需要转账的金额，进行找零：给from创建一个output
	if retValue > amount {
		output2 := newTXOutput(from, retValue - amount)
		outputs = append(outputs, output2)
	}
	timeStamp := time.Now().Unix()
	//6 设置哈希 返回
	tx := Transaction{nil, inputs, outputs, uint64(timeStamp)}
	tx.setHash()
	return &tx

}

//判断一笔交易是否为挖矿交易
func (tx *Transaction) isCoinbaseTx() bool {
	inputs := tx.TXInputs
	//input个数为1，id为nil，索引为-1
	if len(inputs) == 1 && inputs[0].Txid == nil && inputs[0].Index == -1 {
		return true
	}
	return false
}

//trim修剪, 签名和校验时都会使用
func (tx *Transaction)trimmedCopy()*Transaction  {
	var inputs []TXIput
	//var outputs []TXOutput
	//创建一个交易副本，每一个input的pubKey和Sig都设置为空。
	for _,input:=range tx.TXInputs{

		input:=TXIput{
			input.Txid,
			input.Index,
			nil,
			nil,
		}
		inputs=append(inputs,input)
	}
	//outputs=tx.TXOutputs
	txCopy:=Transaction{tx.TXID,inputs,tx.TXOutputs,tx.TimeStamp}
	return &txCopy
}

//实现具体签名动作（copy，设置为空，签名动作）
//参数1：私钥
//参数2：inputs所引用的output所在交易的集合:
// > key :交易id
// > value：交易本身
func (tx *Transaction)sign(priKey *ecdsa.PrivateKey,prevTxs map[string]*Transaction)bool{

	fmt.Println("具体对交易签名。。。。")
	if tx.isCoinbaseTx(){
		fmt.Println("找到挖矿交易，无需签名")
		return true
	}//if
	//1. 获取交易copy，pubKey，ScriptPubKey字段置空
	txCopy:=tx.trimmedCopy()
	//2. 遍历交易的inputs for, 注意，不要遍历tx本身，而是遍历txCopy
	for i,input:=range txCopy.TXInputs{
		fmt.Printf("开始对input[%d]进行签名...",i)
		prevTx:=prevTxs[string(input.Txid)]
		if prevTx==nil{
			return false
		}
		//input引用的output
		output:=prevTx.TXOutputs[input.Index]
		// > 获取引用的output的公钥哈希
		//for range是input是副本，不会影响到变量的结构
		// input.PubKey = output.ScriptPubKeyHash
		txCopy.TXInputs[i].PubKey=output.ScriptPubKeyHash
		//对copy交易进行签名 需要得到交易的哈希值
		txCopy.setHash()
		//将input的pubKey字段设置为nil，还原数据，防止干扰后面的input的签名
		txCopy.TXInputs[i].PubKey=nil
		hashData:=txCopy.TXID//我们去签名的具体数据
		//开始签名
		r,s,err:=ecdsa.Sign(rand.Reader,priKey,hashData)
		if err!=nil{
			fmt.Printf("签名失败")
			return false
		}
		signature:=append(r.Bytes(),s.Bytes()...)
		//将数字签名赋值给原始tx
		tx.TXInputs[i].Scriptsig=signature
	}//for
	fmt.Printf("交易签名成功")
	return true
}
//具体校验
func (tx *Transaction)verify(prevTxs map[string]*Transaction)bool{
	//1 获取交易副本txCopy
	txCopy:=tx.trimmedCopy()
	//2 遍历交易 inputs
	for i,input:=range tx.TXInputs{
		prevTx:=prevTxs[string(input.Txid)]
		if prevTx==nil{
			return false
		}
		//还原数据（得到引用Output的公钥哈希 ） 获取交易的哈希值
		output:=prevTx.TXInputs[input.Index]
		txCopy.TXInputs[i].PubKey=output.Scriptsig
		txCopy.setHash()
		//清零环境 设置为Nil
		txCopy.TXInputs[i].PubKey=nil
		//具体还原的签名数据哈希值
		hashData:=txCopy.TXID
		//签名
		signature:=input.Scriptsig
		//公钥的字节流
		pubKey:=input.PubKey
		//开始校验
		var r,s,x,y big.Int
		//r,s 从signature 截取出来
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])
		//x,y从pubkey截取出来 还原为公钥本身
		x.SetBytes(pubKey[:len(pubKey)/2])
		y.SetBytes(pubKey[len(pubKey)/2:])
		curve:=elliptic.P256()
		pubKeyRaw:=ecdsa.PublicKey{curve,&x,&y}
		//进行校验
		res:=ecdsa.Verify(&pubKeyRaw,hashData,&r,&s)
		if !res{
			fmt.Println("发现校验失败的input")
			return false
		}
	}
	//4 通过tx.scriptsig ,tx.Pubkey进行校验
	fmt.Println("交易检验成功")
	return true
}
func (tx *Transaction)String()string{
	var lines []string
	lines=append(lines,fmt.Sprintf("-----Trasaction %x:",tx.TXID))
	for i, input := range tx.TXInputs {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Index))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Scriptsig))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.TXOutputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %f", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.ScriptPubKeyHash))
	}
return strings.Join(lines,"\n")
}