package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func uintToByte(num uint64) []byte {
	//todo
	var buffer bytes.Buffer
	//使用二进制编码
	//func Write(w io.Writer, order ByteOrder, data interface{}) error
	err := binary.Write(&buffer, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println("binary.Write err:", err)
		return nil
	}
	return buffer.Bytes()
}
//判断文件是否存在
func isFileExist(filename string)bool {
	//
	_,err:= os.Stat(filename)
	if err!=nil{
		return false
	}
	return true
}
