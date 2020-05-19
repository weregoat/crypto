package aes

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

// https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38a.pdf
var testVectors = []struct {
	Key       string
	IV        string
	PlainText string
	CipherText string
}{
	{ // Block #1
		Key:       "2b7e151628aed2a6abf7158809cf4f3c",
		IV:        "f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff",
		PlainText: "6bc1bee22e409f96e93d7e117393172a",
		CipherText: "874d6191b620e3261bef6864990db6ce",
	},
	{ // Block #2
		Key:       "2b7e151628aed2a6abf7158809cf4f3c",
		IV: "f0f1f2f3f4f5f6f7f8f9fafbfcfdff00",
		PlainText: "ae2d8a571e03ac9c9eb76fac45af8e51",
		CipherText: "9806f66b7970fdff8617187bb9fffdff",
	},
	{ // Block #3
		Key:       "2b7e151628aed2a6abf7158809cf4f3c",
		IV: "f0f1f2f3f4f5f6f7f8f9fafbfcfdff01",
		PlainText: "30c81c46a35ce411e5fbc1191a0a52ef",
		CipherText: "5ae4df3edbd5d35e5b4f09020db03eab",
	},
	{ // Block #4
		Key:       "2b7e151628aed2a6abf7158809cf4f3c",
		IV: "f0f1f2f3f4f5f6f7f8f9fafbfcfdff02",
		PlainText: "f69f2445df4f9b17ad2b417be66c3710",
		CipherText: "1e031dda2fbe03d1792170a0f3009cee",
	},
}

func TestEncryptBlock(t *testing.T) {

	for _,test := range testVectors {
		dst := make([]byte, blockSize)
		key, err := hex.DecodeString(test.Key)
		if err != nil {
			t.Error(err)
			continue
		}
		iv, err := hex.DecodeString(test.IV)
		if err != nil {
			t.Error(err)
			continue
		}
		plainText, err := hex.DecodeString(test.PlainText)
		if err != nil {
			t.Error(err)
			continue
		}
		cipherText, err := hex.DecodeString(test.CipherText)
		if err != nil {
			t.Error(err)
			continue
		}
		processBlock(key, iv, plainText, dst)
		if !bytes.Equal(dst, cipherText) {
			t.Errorf(
				"Expected ciphertext to be %x, but got %x",
				cipherText,
				dst,
				)
		}
	}

}

func TestDecryptBlock(t *testing.T) {

	for _,test := range testVectors {
		dst := make([]byte, blockSize)
		key, err := hex.DecodeString(test.Key)
		if err != nil {
			t.Error(err)
			continue
		}
		iv, err := hex.DecodeString(test.IV)
		if err != nil {
			t.Error(err)
			continue
		}
		plainText, err := hex.DecodeString(test.PlainText)
		if err != nil {
			t.Error(err)
			continue
		}
		cipherText, err := hex.DecodeString(test.CipherText)
		if err != nil {
			t.Error(err)
			continue
		}
		processBlock(key, iv, cipherText, dst)
		if !bytes.Equal(dst, plainText) {
			t.Errorf(
				"Expected plaintext to be %x, but got %x",
				plainText,
				dst,
			)
		}
	}

}

func TestDecrypt(t *testing.T) {
	cipherText, err := base64.StdEncoding.DecodeString("L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ==")
	solution := "Yo, VIP Let's kick it Ice, Ice, baby Ice, Ice, baby "
	if err != nil {
		t.Error(err)
	}

	plainText, err := Decrypt([]byte("YELLOW SUBMARINE"), cipherText)
	if err != nil {
		t.Error(err)
	}
	if plainText != solution {
		t.Errorf("expecting %+q, got %+q", solution, plainText)

	}
}

func TestEncrypt(t *testing.T) {
	solution, err := base64.StdEncoding.DecodeString("L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ==")
	if err != nil {
		t.Error(err)
	}
	cipherText, err := Encrypt([]byte("YELLOW SUBMARINE"), "Yo, VIP Let's kick it Ice, Ice, baby Ice, Ice, baby ")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(cipherText, solution) {
		t.Errorf("expecting %+q, got %+q", solution, cipherText)

	}
}