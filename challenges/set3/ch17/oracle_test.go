package ch17

import (
	"bytes"
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"testing"
)

// Test with padded plaintexts
func TestOKPadding(t *testing.T) {
	tests := []struct {
		PlainText string
		Padded  bool
	}{
		{"BBBBBBBBBBBBBBB\x01", true},
		{"BBBBBBBBBBBBBB\x02\x02", true},
		{"BBBBBBBBBBBBB\x03\x03\x03", true},
		{"BBBBBBBBBBBB\x04\x04\x04\x04", true},
		{"BBBBBBBBBBBBBBB\x02", false},
		{"BBBBBBBBBBBBBB\x01\x01", true},
		{"BBBBBBBBBBBBBB\x01\x02", false},
		{"BBBBBBBBBBBBBB\x03\x03", false},
		{"BBBBBBBBBBBBBBB\x00", false},
	}
	prefix := bytes.Repeat([]byte{'A'}, cbc.BlockSize)
	for _,test := range tests {
		text := append(prefix, []byte(test.PlainText)...)
		o, err := NewOracle()
		if err != nil {
			t.Error(err)
			break
		}
		c, err := cbc.Encrypt(text, o.Key, o.IV)
		if o.CheckPadding(c) != test.Padded {
			t.Errorf("expecting oracle to return %t for text in %+q", test.Padded, text)
		}
	}
}



func TestOracle_Encrypt(t *testing.T) {
	for i:=0; i < 50; i++ {
		o, err := NewOracle()
		if err != nil {
			t.Error(err)
			break
		}
		c, iv := o.Encrypt()
		plainText, err := cbc.Decrypt(c, o.Key, iv)
		if string(pkcs7.RemovePadding(plainText)) != o.Plaintext {
			t.Errorf("expecting ciphertext to decrypt to %+q, got %+q",
				o.Plaintext,
				pkcs7.RemovePadding(plainText),
				)
		}
	}
}
