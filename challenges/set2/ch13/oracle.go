package ch13

import (
	"fmt"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
	"strings"
)

type Oracle struct {
	key       []byte
	Error     error
	plainText []byte
}

var debug = false

func New() (Oracle, error) {
	o := Oracle{}
	key, err := util.RandomBytes(ecb.BlockSize)
	if err != nil {
		o.Error = err
		return o, err
	}
	o.key = key
	return o, nil
}

func (o Oracle) Encrypt(email string) []byte {
	data := encode(email)
	if debug {
		fmt.Printf("oracle got: %+q\n", data)
	}
	cipherText, err := ecb.Encrypt([]byte(data), o.key)
	if err != nil {
		o.Error = err
	}
	return cipherText
}

func (o Oracle) Decrypt(src string) map[string]string {
	data, err := ecb.Decrypt([]byte(src), o.key)
	if err != nil {
		o.Error = err
	}
	if debug {
		fmt.Printf("oracle returns: %+q\n", data)
	}
	return parse(string(pkcs7.RemovePadding(data)))
}

func parse(src string) map[string]string {
	var data = make(map[string]string)
	for _, i := range strings.Split(src, "&") {
		parts := strings.Split(i, "=")
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}
	return data
}

func encode(email string) string {
	uid := 10
	role := "user"
	stripped := strings.ReplaceAll(strings.ReplaceAll(email, "&", ""), "=", "")
	return fmt.Sprintf("email=%s&uid=%d&role=%s", stripped, uid, role)
}
