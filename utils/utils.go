package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

/**
* @des 时间转换函数
* @param atime string 要转换的时间戳（秒）
* @return string
 */
func StrTime(atime int64) string {
	var byTime = []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1}
	var unit = []string{"年", "天", "小时", "分钟", "秒钟"}
	now := time.Now().Unix()
	ct := now - atime
	if ct < 0 {
		return "已结束"
	}
	var res string
	for i := 0; i < len(byTime); i++ {
		if ct < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(ct / byTime[i]))
		ct = ct % byTime[i];
		if temp > 0 {
			var tempStr string
			tempStr = strconv.FormatFloat(temp, 'f', -1, 64)
			res = MergeString(tempStr, unit[i]) //此处调用了一个我自己封装的字符串拼接的函数（你也可以自己实现）
		}
		break //我想要的形式是精确到最大单位，即："2天前"这种形式，如果想要"2天12小时36分钟48秒前"这种形式，把此处break去掉，然后把字符串拼接调整下即可（别问我怎么调整，这如果都不会我也是无语）
	}
	return res
}

/**
* @des 拼接字符串
* @param args ...string 要被拼接的字符串序列
* @return string
 */
func MergeString(args ...string) string {
	buffer := bytes.Buffer{}
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
	}
	return buffer.String()
}

func CreateCaptcha() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}
func Md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func FormatTokens(tokens float64, p int) string {
	if tokens > 0 {
		str := strconv.FormatFloat(tokens/1000000000000000000, 'f', p, 64)
		return str
	} else {
		return "0"
	}
}

func GetRealAebalanceBigInt(amount float64) *big.Int {
	newFloat := big.NewFloat(amount)
	basefloat := big.NewFloat(1000000000000000000)
	float1 := big.NewFloat(1)
	float1.Mul(newFloat, basefloat)
	resultAmount := new(big.Int)
	float1.Int(resultAmount)
	return resultAmount
}

func GetRealAebalanceFloat64(amount float64) float64 {
	return amount * 1000000000000000000
}
func GetAEFloat64(amount float64) float64 {
	return amount / 1000000000000000000
}
func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}
func AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

//补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func IsEmail(email string) bool {
	// 识别电子邮件地址
	isorno, _ := regexp.MatchString(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`, email)

	if isorno {
		return true
	} else {
		return false
	}
}

//检测是不是手机访问
func IsMobile(userAgent string) bool {
	mobileRe, _ := regexp.Compile("(?i:Mobile|iPod|iPhone|Android|Opera Mini|BlackBerry|webOS|UCWEB|Blazer|PSP)")
	if mobileRe.FindString(userAgent) == "" {
		return false
	}
	return true
}
