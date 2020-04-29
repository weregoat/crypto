package ch11

import (
	"bytes"
	"crypto/aes"
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
	"testing"
)

func TestOracle_New(t *testing.T) {
	tests := map[int]int{
		0:  0,
		-1: 0,
		16: 16,
		31: 31,
	}
	for size, expected := range tests {
		o := New(size)
		if o.BlockSize != expected {
			t.Errorf("expecting blocksize of %d, got %d", expected, size)
		}
	}
}

func TestOracle_Encrypt(t *testing.T) {
	blockSize := aes.BlockSize
tests:
	for i := 0; i < 50; i++ {
		pSize := blockSize * (1 + (i % 3)) // 16,32,48
		text, err := util.RandomBytes(pSize)
		if err != nil {
			t.Error(err)
			break tests
		}
		//t.Logf("PlainText: % x", text)
		o := New(blockSize)
		err = o.Encrypt(text)
		if err != nil {
			t.Error(err)
			break tests
		}
		//t.Logf("Mode: %s", o.Mode)
		switch o.Mode {
		case ModeCBC:
			cipherText, err := cbc.Encrypt(o.PlainText, o.Key, o.IV)
			//t.Logf("Expected ciphertext: % x", cipherText)
			//t.Logf("Oracle ciphertext: % x", o.CipherText)
			if err != nil {
				t.Error(err)
				break tests
			}
			if !bytes.Equal(o.CipherText, cipherText) {
				t.Fail()
				break tests
			}
			plainText, err := cbc.Decrypt(o.CipherText, o.Key, o.IV)
			//t.Logf("Expected plaintext: % x", plainText)
			//t.Logf("Oracle plaintext: % x", o.PlainText)
			if err != nil {
				t.Log(err)
				//break tests
			}
			if !bytes.Equal(o.PlainText, pkcs7.RemovePadding(plainText)) {
				t.Fail()
				break tests
			}

		case ModeECB:
			cipherText, err := ecb.Encrypt(o.PlainText, o.Key)
			//t.Logf("Expected ciphertext: % x", cipherText)
			//t.Logf("Oracle ciphertext: % x", o.CipherText)
			if err != nil {
				t.Error(err)
				break tests
			}
			if !bytes.Equal(o.CipherText, cipherText) {
				t.Fail()
				break tests
			}
			plainText, err := ecb.Decrypt(o.CipherText, o.Key)
			//t.Logf("Expected plaintext: % x", plainText)
			//t.Logf("Oracle plaintext: % x", o.PlainText)
			if err != nil {
				t.Error(err)
				break tests
			}
			if !bytes.Equal(o.PlainText, pkcs7.RemovePadding(plainText)) {
				t.Fail()
				break tests
			}
		}

	}
}
