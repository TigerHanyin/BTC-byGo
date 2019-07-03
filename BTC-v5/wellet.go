package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcutil/base58"
	"bytes"
)

type wallet struct {
	PriKey *ecdsa.PrivateKey
	//公钥原型定义
	// type PublicKey struct {
	// 	elliptic.Curve
	// 	X, Y *big.Int
	// }

	// 公钥, X,Y类型一致，长度一致，我们将X和Y拼接成字节流，赋值给pubKey字段，用于传输
	// 验证时，将X，Y截取出来（类似r,s),再创建一条曲线，就可以还原公钥，进一步进行校验
	PubKey []byte
}
func newWalletKeyPair()*wallet{
	//曲线
	curve:=elliptic.P256()
	//创建私钥
	priKey,err:= ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{
		fmt.Println("ecdsa.GenetateKey err:",err)
		return nil
	}
	//获取公钥
	pubKeyRaw:=priKey.PublicKey
	//将公钥x,y拼接到一起
	pubKey:=append(pubKeyRaw.X.Bytes(),pubKeyRaw.Y.Bytes()...)
	//创建wallet结构返回
	wallet:=wallet{priKey,pubKey}
	return &wallet
}
//给定公钥 -》公钥哈希
func getPubKeyHashFromPubKey(pubKey []byte)[]byte{
	//公钥
	hash1:=sha256.Sum256(pubKey)

	//hash160处理
	hasher:=ripemd160.New()
	hasher.Write(hash1[:])

	//公钥哈希，锁定output时就用这个值
	pubKeyHash:=hasher.Sum(nil)
	return pubKeyHash
}
//根据私钥生成地址
func (w*wallet)getAddress()string{
	pubKeyHash:=getPubKeyHashFromPubKey(w.PubKey)
	//凭借version和公钥哈希，得到21字节的数据
	payload:=append([]byte{byte(0x00)},pubKeyHash...)
	//生成4字节的校验码
	checksum:=checkSum(payload)
	//25字节数据
	payload=append(payload,checksum...)
	address:=base58.Encode(payload)
	return address
}
//地址-->公钥哈希
func getPubKeyHashFromAddress(address string)[]byte{
	//base58解码
	decodeInfo:=base58.Decode(address)
	if len(decodeInfo)!=25{
		fmt.Println("getPubKeyHashFromAddress,传入地址无效")
		return nil
	}
	////需要校验一下地址

	//截取
	pubKeyHash:=decodeInfo[1:len(decodeInfo)-4]
	return pubKeyHash
}
//得到4字节的校验码
func checkSum(payload []byte) []byte {
	first:=sha256.Sum256(payload)
	second:=sha256.Sum256(first[:])
	//4字节的checksum
	checksum:=second[:4]
	return checksum
}
func isValidAddress(addr string)bool{
	//解码 得到25字节数据
	decodeInfo:=base58.Decode(addr)
	if len(decodeInfo)!=25{
		fmt.Println("isValidAddress ，传入地址长度无效")
		return false
	}
	//截取前21 字节的payload 截取后字节的checksum
	payload:=decodeInfo[:len(decodeInfo)-4]//21字节
	checksum1:=decodeInfo[len(decodeInfo)-4:]//4字节
	// 对palyload计算，得到checksum2，与checksum1对比，true校验成功，反之失败
	checksum2:=checkSum(payload)
	return bytes.Equal(checksum1,checksum2)
}
