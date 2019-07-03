package main

import (
	"github.com/bolt"
	"errors"
	"fmt"
)

//提供计算
//定义区块链结构 用数组模拟区块链
type BlockChain struct {
	db   *bolt.DB //存储数据
	tail []byte   //最后一个区块的哈希值
}

//创世语
const genesisInfo = "The Time 03/Jan/2009 Chancellor on brink of second bailput for banks"
const blockchainDBFile = "blockchain.db"
const bucketBlock = "bucketBlock"           //装block的桶
const lastBlockHashKey = "lastBlockHashKey" //用于访问bolt数据库，得到最好一个区块的哈希值

//创建区块，从无到有：这个函数仅执行一次
func CreateBlockChain() error {
	// 1. 区块链不存在，创建
	db, err := bolt.Open(blockchainDBFile, 0600, nil)
	if err != nil {
		return err
	}

	//不要db.Close，后续要使用这个句柄
	defer db.Close()

	// 2. 开始创建
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		//如果bucket为空，说明不存在
		if bucket == nil {
			//创建bucket
			bucket, err := tx.CreateBucket([]byte(bucketBlock))
			if err != nil {
				return err
			}
			//1 创建挖矿交易
			coinbase := NewCoinbaseTx("中本聪", genesisInfo)
			//2 拼装 txs
			txs := []*Transaction{coinbase}
			//3 创建创世快
			//func NewBlock(txs []*Transaction, prevHash []byte) *Block {
			genesisBlock := NewBlock(txs, nil)
			//key是区块的哈希值，value是block的字节流
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize()) //将block序列化
			//更新最后区块哈希值到数据库
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
		}
		return nil
	})
	return err //nil
}

//获取区块链实例，用于后续操作, 每一次有业务时都会调用
func GetBlockChainInstance() (*BlockChain, error) {

	var lastHash []byte //内存中最后一个区块的哈希值

	//两个功能：
	// 1. 如果区块链不存在，则创建，同时返回blockchain的示例
	db, err := bolt.Open(blockchainDBFile, 0600, nil) //rwx  0100 => 4

	if err != nil {
		return nil, err
	}

	//不要db.Close，后续要使用这个句柄

	// 2. 如果区块链存在，则直接返回blockchain示例
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketBlock))

		//如果bucket为空，说明不存在
		if bucket == nil {
			return errors.New("bucket不应为nil")
		} else {
			//直接读取特定的key，得到最后一个区块的哈希值
			lastHash = bucket.Get([]byte(lastBlockHashKey))
		}

		return nil
	})

	//5. 拼成BlockChain然后返回
	bc := BlockChain{db, lastHash}
	return &bc, nil
}

//提供一个向区块链中添加区块的方法
//func (bc *BlockChain) AddBlock(data string) error {
//
//	lashBlockHash := bc.tail //区块链中最后一个区块的哈希值
//
//	//1. 创建区块
//	newBlock := NewBlock(data, lashBlockHash)
//
//	//2. 写入数据库
//	err := bc.db.Update(func(tx *bolt.Tx) error {
//		bucket := tx.Bucket([]byte(bucketBlock))
//		if bucket == nil {
//			return errors.New("AddBlock时Bucket不应为空")
//		}
//
//		//key是新区块的哈希值， value是这个区块的字节流
//		bucket.Put(newBlock.Hash, newBlock.Serialize())
//		bucket.Put([]byte(lastBlockHashKey), newBlock.Hash)
//
//		//更新bc的tail，这样后续的AddBlock才会基于我们newBlock追加
//		bc.tail = newBlock.Hash
//		return nil
//	})
//
//	return err
//}
//====================迭代器 遍历区块链=================
//定义迭代器
type Iterator struct {
	db          *bolt.DB
	currentHash []byte //不断移动的哈希值，由于访问所有区块
}

func (bc *BlockChain) NewIterator() *Iterator {
	it := Iterator{
		db:          bc.db,
		currentHash: bc.tail,
	}
	return &it
}

//给Iterator绑定一个方法：Next
//1. 返回当前所指向的区块
//2. 向左移动（指向前一个区块）
func (it *Iterator) Next() (block *Block) {
	//读取Bucket当前哈希block
	err := it.db.View(func(tx *bolt.Tx) error {
		//读取bucket
		bucket := tx.Bucket([]byte(bucketBlock))
		if bucket == nil {
			return errors.New("Iterator Next时bucket 不应为Nil")
		}

		fmt.Printf("currentHash:%x\n", it.currentHash)
		blockTmpInfo /**block 字节流*/ := bucket.Get(it.currentHash) //一定要注意，是currentHash
		block = Deserialize(blockTmpInfo)
		it.currentHash = block.PreHash //游标左移
		return nil
	})
	//哈希游标向左移动
	if err != nil {
		fmt.Println("iterator next err:", err)
		return nil
	}
	return
}

//获取指定地址金额，实现遍历账本通用的函数
func (bc *BlockChain) FindMyUTXO(addr string) []TXOutput {
	var utxos []TXOutput
	//定义一个存放已经消耗过的
	spentUtxos:=make(map[string][]int)
	it := bc.NewIterator()
	for {
		//遍历区块
		block := it.Next()
		//遍历交易

		for _, tx := range block.Transactions {
		LABLE:
			//遍历output，判断锁定脚本是否为目标地址
			for outputIndex, output := range tx.TXOutputs {

				if output.ScriptPubk == addr {
					//开始过滤
					currentTxid:=string(tx.TXID)
					//花费中查看
					indexArray:=spentUtxos[currentTxid]
					//若果不为零 说明这个交易ID在篮子中有数据 一定有某个output被使用了
					if len(indexArray)!=0{
						for _,spendIndex/*0,1 */:=range indexArray{
							if outputIndex/*当前的*/==spendIndex{
								continue LABLE
							}//if
						}//for

					}//if

					utxos = append(utxos, output)
				}//if
			}//for
			///=======================遍历inputs===============
			for _,input:=range tx.TXInputs{
				if input.Scriptsig==addr{
					//map [交易ID][]int
					spentKey:=string(input.Txid)
					spentUtxos[spentKey]=append(spentUtxos[spentKey],int(input.Index))
				}//if
			}//for
		}
		//退出
		if len(block.PreHash)==0{
			break
		}

	}
	return utxos

}
