package security

import (
	"bytes"
	"crypto/md5"
	_ "encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func Md5(msg string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(msg))
	//cipherStr := md5Ctx.Sum(nil)
	//fmt.Print(cipherStr)
	//fmt.Print("\n")
	//fmt.Print(hex.EncodeToString(cipherStr))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func GenAppKey(appId string) string {
	//md5.
	return string(Md5(fmt.Sprintf("%s_%d%d", appId, 234, 546574)))
}

func GenSessionID() string {
	t := time.Now().UnixNano()
	return string(Md5(fmt.Sprintf("%d", t)))
}
