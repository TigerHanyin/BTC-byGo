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
	if !isValidAddress(addr){
		fmt.Println("传入地址无效，无效地址为：",addr)
		return
	}
	bc,_:=GetBlockChainInstance()
	defer bc.db.Close()
	//得到查询地址的公钥哈希
	pubKeyHash:=getPubKeyHashFromAddress(addr)
	//获取所有相关的utxo集合
	utxinfos:=bc.FindMyUTXO(pubKeyHash)
	total:=0.0
	for _,utxo:=range utxinfos{
		total+=utxo.TXOutput.Value
	}
	fmt.Printf("%s的金额为%f\n",addr,total)
}
func (cli *CLI)send(from, to string, amount float64, miner, data string){
	bc,err:=GetBlockChainInstance()
	if err!=nil{
		fmt.Println("send err:",err)
		return
	}
	defer bc.db.Close()
	//创建挖矿交易
	coinbaseTx:=NewCoinbaseTx(miner,data)
	//创建txs数组，将有效交易添加进来
	txs:=[]*Transaction{coinbaseTx}
	//创建普通交易
	tx:=NewTransaction(from,to,amount,bc)
	if tx!=nil{
		fmt.Println("找到一笔有效的转账交易" )
		txs=append(txs,tx)
	}else{
		fmt.Println("注意，找到一笔无效的转账交易，不添加到区块")
	}
	//调用Addblock
	err=bc.AddBlock(txs)
	if err!=nil{
		fmt.Println("添加区块失败，转账失败")
	}
	fmt.Println("添加区块成功，转账成功")
}
func(cli *CLI)createWallet(){
	wm:=NewWalletManager()
	if wm==nil{
		fmt.Println("createWallet 失败")
		return
	}
	address:=wm.createWallet()
	if len(address)==0{
		fmt.Println("创建钱包失败！")
		return
	}
	fmt.Println("新钱包地址为：",address)

}
func (cli *CLI)listAddress(){
	wm:=NewWalletManager()
	if wm==nil{
		fmt.Println("newwalletManager 失败")
		return
	}
	addresses:=wm.listAddresses()
	for _,address:=range addresses{
		fmt.Printf("%s\n",address)
	}

}
func (cli *CLI)printTx(){
	bc,err:=GetBlockChainInstance()
	if err!=nil{
		fmt.Println("printTx err:",err)
		return
	}
	defer bc.db.Close()
	it:=bc.NewIterator()
	for {
		block:=it.Next()
		fmt.Println("\n..........区块分割线...........")
		for _,tx:=range block.Transactions{
			//直接打印交易
			fmt.Println(tx)
		}
		if len(block.PreHash)==0{
			break
		}
	}

}