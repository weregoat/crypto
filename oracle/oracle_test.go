package oracle

import (
	"bytes"
	"crypto/aes"
	"gitlab.com/weregoat/crypto/util"
	"testing"
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
)

func TestOracle_New(t *testing.T) {
	tests := map[int]int {
		0:0,
		-1:0,
		16:16,
		31:31,
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
	pSize := blockSize
	for i := 0; i < 50; i++ {
		text, err := util.RandomBytes(pSize)
		if err != nil {
			t.Error(err)
		}
		t.Logf("PlainText: %q", util.EncodeToBase64(text))
		o := New(blockSize)
		err = o.Encrypt(text)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Mode: %s", o.Mode)
		switch o.Mode {
		case ModeCBC:
			cipherText, err := cbc.Encrypt(o.PlainText, o.Key, o.IV)
			if err != nil {
				t.Logf("Expected ciphertext: %q", util.EncodeToBase64(cipherText))
				t.Logf("Oracle ciphertext: %q", util.EncodeToBase64(o.CipherText))
				t.Error(err)
			}
			if !bytes.Equal(o.CipherText, cipherText) {
				t.Fail()
			}
			plainText, err := cbc.Decrypt(o.CipherText, o.Key, o.IV)

			if err != nil {
				t.Error(err)
			}
			if ! bytes.Equal(o.PlainText, plainText) {
				t.Logf("Expected plaintext: %q", util.EncodeToBase64(plainText))
				t.Logf("Oracle plaintext: %q", util.EncodeToBase64(o.PlainText))
				t.Fail()
			}

		case ModeECB:
			cipherText, err := ecb.Encrypt(o.PlainText, o.Key)
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(o.CipherText, cipherText) {
				t.Logf("Expected ciphertext: %q", util.EncodeToBase64(cipherText))
				t.Logf("Oracle ciphertext: %q", util.EncodeToBase64(o.CipherText))
				t.Fail()
			}
			plainText, err := ecb.Decrypt(o.CipherText, o.Key)
			if err != nil {
				t.Error(err)
			}
			if ! bytes.Equal(o.PlainText, plainText) {
				t.Logf("Expected plaintext: %q", util.EncodeToBase64(plainText))
				t.Logf("Oracle plaintext: %q", util.EncodeToBase64(o.PlainText))
				t.Fail()
			}
		}

	}
}
