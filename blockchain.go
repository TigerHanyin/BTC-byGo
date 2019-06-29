package main
//提供计算
//定义区块链结构 用数组模拟区块链
type BlockChain struct {
	Blocks []*Block//区块链
}
//创世语
const genesisInfo ="The Time 03/Jan/2009 Chancellor on brink of second bailput for banks"
//提供一个创建区块链的方法
func NewBlockChain ()*BlockChain{
	//创建创世链，同时添加一个创世块
	genesisBlock:=NewBlock(genesisInfo,nil)

	bc:=BlockChain{
		Blocks:[]*Block{genesisBlock},
	}
	return &bc
}
func (bc *BlockChain)AddBlock(data string){
	lastBlock:=bc.Blocks[len(bc.Blocks)-1]
	//最后一个区块的哈希值是新区块的前哈细致
	preHash:=lastBlock.Hash
	//创建block
	newBlock:=NewBlock(data,preHash)
	//添加bc中
	bc.Blocks=append(bc.Blocks,newBlock)
}
