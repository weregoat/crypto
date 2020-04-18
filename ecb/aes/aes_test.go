package aes

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
	"testing"
)

func TestDecryptAES128(t *testing.T) {
	key := []byte("YELLOW SUBMARINE")
	cipher, err := aes.NewCipher(key)
	if err != nil {
		t.Error(err)
	}
	for j := 0; j < 5000; j++ {
		plainText, err := util.RandomBytes(util.RandomInt(0, 64))
		if err != nil {
			t.Error(err)
		}
		/* The plaintext should not be already padded.
		 But it's random generated, so it might end up in a way that
		 is padded according to PKCS#7 (e.g. multiple of 16 bytes and
		 last byte \x01.
		In this case I am adding a byte to break it.
		*/
		if pkcs7.IsPadded(plainText) {
			plainText = append(plainText, 'Z')
		}
		padded := make([]byte, len(plainText))
		if len(plainText)%aes.BlockSize != 0 {
			padded = pkcs7.Pad(plainText, aes.BlockSize)
		} else {
			copy(padded, plainText)
		}
		var cipherText = make([]byte, len(padded))
		i := 0
		for {
			s := i * aes.BlockSize
			if s >= len(padded) {
				break
			}
			e := (i + 1) * aes.BlockSize
			cipher.Encrypt(cipherText[s:e], padded[s:e])
			i++
		}
		decrypted, err := Decrypt(cipherText, key)
		if err != nil {
			t.Error(err)
		}
		if string(decrypted) != string(plainText) {
			t.Errorf("expecting %q, got %q", plainText, decrypted)
		}

		shortKey := key[:12]
		_, err = Decrypt(cipherText, shortKey)
		if err == nil {
			t.Errorf("expecting error because of short key, but got nothing")
		}
		if len(plainText) > 0 {
			shortText := cipherText[:12]
			_, err = Decrypt(shortText, key)
			if err == nil {
				t.Errorf("expecting error because of short cipherText, but got nothing")
			}
		}
	}
}

func TestAES128(t *testing.T) {
	// Test vectors from https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/block-ciphers#AES
	/*
		Notice everything is Hex encoded; key too.
		And, no, I didn't implement them all, just a few samples to verify my
		implementation is not completely broken.
	*/

	testVectors := []struct {
		Key        string
		PlainText  string
		CipherText string
	}{
		// ECBGFSbox128.rsp
		// 0
		{
			"00000000000000000000000000000000",
			"f34481ec3cc627bacd5dc3fb08f273e6",
			"0336763e966d92595a567cc9ce537f5e",
		},
		// 6
		{
			"00000000000000000000000000000000",
			"58c8e00b2631686d54eab84b91f0aca1",
			"08a4e2efec8a8e3312ca7460b9040bbf",
		},
		// EBCKeySbox128.rsp
		// 0
		{
			"10a58869d74be5a374cf867cfb473859",
			"00000000000000000000000000000000",
			"6d251e6944b051e04eaa6fb4dbf78465",
		},
		// 3
		{
			"b6364ac4e1de1e285eaf144a2415f7a0",
			"00000000000000000000000000000000",
			"5d9b05578fc944b3cf1ccf0e746cd581",
		},
		// 20
		{
			"febd9a24d8b65c1c787d50a4ed3619a9",
			"00000000000000000000000000000000",
			"f4a70d8af877f9b02b4c40df57d45b17",
		},
		// ECBVarKey128.rsp
		// 0
		{
			"80000000000000000000000000000000",
			"00000000000000000000000000000000",
			"0edd33d3c621e546455bd8ba1418bec8",
		},
		// 7
		{
			"ff000000000000000000000000000000",
			"00000000000000000000000000000000",
			"b1d758256b28fd850ad4944208cf1155",
		},
		// 127
		{
			"ffffffffffffffffffffffffffffffff",
			"00000000000000000000000000000000",
			"a1f6258c877d5fcd8964484538bfc92c",
		},
		// ECBVarTxt128.rsp
		// 0
		{
			"00000000000000000000000000000000",
			"80000000000000000000000000000000",
			"3ad78e726c1ec02b7ebfe92b23d9ec34",
		},
		// 32
		{
			"00000000000000000000000000000000",
			"ffffffff800000000000000000000000",
			"171a0e1b2dd424f0e089af2c4c10f32f",
		},
		// 126
		{
			"00000000000000000000000000000000",
			"fffffffffffffffffffffffffffffffe",
			"5c005e72c1418c44f569f2ea33ba54f3",
		},
	}

	for _, tv := range testVectors {
		key, err := hex.DecodeString(tv.Key)
		if err != nil {
			t.Error(err)
		}
		cipherText, err := hex.DecodeString(tv.CipherText)
		if err != nil {
			t.Error(err)
		}
		plainText, err := hex.DecodeString(tv.PlainText)
		if err != nil {
			t.Error(err)
		}
		encrypt, err := Encrypt(plainText, key)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(encrypt, cipherText) {
			t.Errorf("expecting cipherText to be %q, got %q", cipherText, encrypt)
		}
		decrypt, err := Decrypt(cipherText, key)
		if err != nil {
			t.Log(err)
		}
		if !bytes.Equal(decrypt, plainText) {
			t.Errorf("expecting plainText to be %q, got %q", plainText, decrypt)
		}
	}

}

func TestEncryptDecript(t *testing.T) {
	plainText := "The simplest of the encryption modes is the Electronic Codebook (ECB) mode (named after conventional physical codebooks)."
	key := []byte("YELLOW SUBMARINE")
	cipherText, err := Encrypt([]byte(plainText), key)
	if err != nil {
		t.Errorf("%q", cipherText)
	}
	p, err := Decrypt(cipherText, key)
	if err != nil {
		t.Error(err)
	}
	if string(p) != plainText {
		t.Errorf("expecting %q, got %q", plainText, p)
	}
}

func TestEmptyPlainText(t *testing.T) {
	plainText := ""
	key := []byte("YELLOW SUBMARINE")
	cipherText, err := Encrypt([]byte(plainText), key)
	if err != nil {
		t.Error(err)
	}
	p, err := Decrypt(cipherText, key)
	if err != nil {
		t.Error(err)
	}
	if string(p) != plainText {
		t.Errorf("expecting %x, got %x", plainText, p)
	}
}
