package main

import "fmt"

func (cli *CLI)addBlock(data string){
	//fmt.Println("添加区块被调用！")
	//bc,_:=GetBlockChainInstance()
	//err:=bc.AddBlock(data)
	//if err!=nil{
	//	fmt.Println("AddBlock failed :",err)
	//	return
	//}
	//fmt.Println("添加区块成功")
}
func (cli *CLI)createBlockChain(){
	err:=CreateBlockChain()
	if err!=nil{
		fmt.Println("CreateBlockChain failed :",err)
		return
	}
	fmt.Println("创建区块链成功")
}
func (cli *CLI)print(){
	bc,_:=GetBlockChainInstance()
	//调用迭代器，输出blockChain
	it:=bc.NewIterator()
	for{
		//调用Next方法 获取区块 游标左移
		block:=it.Next()
		fmt.Printf("\n++++++++++++++++++++++\n")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevHash : %x\n", block.PreHash)
		fmt.Printf("MerkleRoot : %x\n", block.MerkleRoot)
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Bits : %d\n", block.Bits)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Hash : %x\n", block.Hash)
		//fmt.Printf("Data : %s\n", block.Data)
		fmt.Printf("Data:%s\n",block.Transactions[0].TXInputs[0].Scriptsig)
		pow:=NewProofOfWork(block)
		fmt.Printf("IsVaild:%v\n",pow.IsValid())
		if block.PreHash==nil{
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}
func (cli *CLI)getBalance(addr string){
	bc,_:=GetBlockChainInstance()
	//获取utxo集合
	utxos:=bc.FindMyUTXO(addr)
	total:=0.0
	for _,utxo:=range utxos{
		total+=utxo.Value
	}
	fmt.Printf("%s的金额为%f\n",addr,total)
}