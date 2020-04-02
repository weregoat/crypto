package aes

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestAES128(t *testing.T) {
	// Test vectors from https://csrc.nist.gov/projects/cryptographic-algorithm-validation-program/block-ciphers#AES

	//	Notice everything is Hex encoded; key too.
	//	And, no, I didn't implement them all, just a few samples to verify my
	//	implementation is not completely broken.

	testVectors := []struct {
		Key        string
		IV         string
		PlainText  string
		CipherText string
	}{
		// CBCGFSbox128.rsp
		// 0
		{
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"f34481ec3cc627bacd5dc3fb08f273e6",
			"0336763e966d92595a567cc9ce537f5e",
		},
		// 5
		{
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"b26aeb1874e47ca8358ff22378f09144",
			"459264f4798f6a78bacb89c15ed3d601",
		},
		// CBCKeySbox128.rsp
		// 0
		{
			"10a58869d74be5a374cf867cfb473859",
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"6d251e6944b051e04eaa6fb4dbf78465",
		},
		// 12
		{
			"6c002b682483e0cabcc731c253be5674",
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"3580d19cff44f1014a7c966a69059de5",
		},
		// 20
		{
			"febd9a24d8b65c1c787d50a4ed3619a9",
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"f4a70d8af877f9b02b4c40df57d45b17",
		},
		// CBCVarKey128.rsp
		// 0
		{
			"80000000000000000000000000000000",
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"0edd33d3c621e546455bd8ba1418bec8",
		},
		// 23
		{
			"ffffff00000000000000000000000000",
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"de11722d893e9f9121c381becc1da59a",
		},
		// 127
		{
			"ffffffffffffffffffffffffffffffff",
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"a1f6258c877d5fcd8964484538bfc92c",
		},
		// CBCVarTxt128.rsp
		// 0
		{
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"80000000000000000000000000000000",
			"3ad78e726c1ec02b7ebfe92b23d9ec34",
		},
		// 44
		{
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"fffffffffff800000000000000000000",
			"85062c2c909f15d9269b6c18ce99c4f0",
		},
		// 127
		{
			"00000000000000000000000000000000",
			"00000000000000000000000000000000",
			"ffffffffffffffffffffffffffffffff",
			"3f5b8cc9ea855a0afa7347d23e8d664e",
		},
	}

	for _, tv := range testVectors {
		key, err := hex.DecodeString(tv.Key)
		if err != nil {
			t.Error(err)
		}
		iv, err := hex.DecodeString(tv.IV)
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
		encrypt, err := Encrypt(plainText, iv, key)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(encrypt, cipherText) {
			t.Errorf("expecting cipherText to be %q, got %q",
				base64.StdEncoding.EncodeToString(cipherText),
				encrypt,
			)
		}
		decrypt, err := Decrypt(cipherText, iv, key)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(decrypt, plainText) {
			t.Errorf("expecting plainText to be %q, got %q",
				base64.StdEncoding.EncodeToString(plainText),
				base64.StdEncoding.EncodeToString(decrypt),
			)
		}
	}
}

func TestMultiBlocks(t *testing.T) {
	plainText := "Ehrsam, Meyer, Smith and Tuchman invented the Cipher Block Chaining (CBC) mode of operation in 1976."
	iv := make([]byte, 16)
	key := []byte("YELLOW SUBMARINE")
	cipherText, err := Encrypt([]byte(plainText), iv, key)
	if err != nil {
		t.Errorf("%q", cipherText)
	}
	p, err := Decrypt(cipherText, iv, key)
	if err != nil {
		t.Error(err)
	}
	if string(p) != plainText {
		t.Errorf("expecting %q, got %q", plainText, p)
	}
}
