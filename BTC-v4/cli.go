package main

import (
	"os"
	"fmt"
)

//处理用户输入命令，完成具体函数的调用
//cli : command line 命令行
type CLI struct {
	//不需要字段
}

//使用说明 ，帮助用户正确使用
const Usage = `
正确使用方法：
	./blockchain create "创建区块链"
	./blockchain addBlock <需要写入的的数据> "添加区块"
	./blockchain print "打印区块链"
`

//负责解析命令的方法
func (cli *CLI) Run() {
	cmds := os.Args
	//用户至少输入两个参数
if len(cmds)<2{
	fmt.Printf("输入参数无效，请检查！")
	fmt.Println(Usage)
	return
}
	switch cmds[1] {
	case "create":
		fmt.Println("创建区块链被调用！")
		cli.createBlockChain()
	case "addBlock":
		if len(cmds)!=3{
			fmt.Println("输入参数无效，请检查！")
			return
		}
		data:=cmds[2]
		cli.addBlock(data)
	case "getBalance":
		if len(cmds)!=3{
			fmt.Println("输入参数无效，请检查！")
			return
		}
		fmt.Println("获取余额被调用！")
		addr:=cmds[2]
		cli.getBalance(addr)

	case "print":
		fmt.Println("打印区块链被调用")
		cli.print()
	default:
		fmt.Println("输入参数无效，请检查！")
		fmt.Println(Usage)
	}




}
