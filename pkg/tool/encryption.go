package tool

import (
	"crypto/md5"
	"fmt"
	"io"
)

func Md5ByString(s string) (string, error) {
	h := md5.New()
	_, err := io.WriteString(h, s)

	if err != nil {
		return "", err
	}

	arr := h.Sum(nil)                  // 将字节数组转换成16进制字符串
	return fmt.Sprintf("%x", arr), nil // %x表示将字节数组转换成16进制字符串

}

func Md5ByBytes(b []byte) string {

	arr := md5.Sum(b) // 将字节数组转换成16进制字符串

	return fmt.Sprintf("%x", arr)
}

func main() {
	// Example usage
	hash, err := Md5ByString("123456")
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}
