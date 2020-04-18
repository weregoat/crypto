package ch14

import (
	"encoding/base64"
	ecb "gitlab.com/weregoat/crypto/ecb/aes"
	"testing"
)

func TestNew(t *testing.T) {
	text := "secret"
	secret := base64.StdEncoding.EncodeToString([]byte(text))
	oracle, err := New(secret)
	if err != nil {
		t.Error(err)
	}
	if string(oracle.Secret) != text {
		t.Errorf("expecting the oracle secret to be %+q, but got %+q", text, oracle.Secret)
	}

}

func TestOracle_Encrypt(t *testing.T) {
	text := "secret"
	secret := base64.StdEncoding.EncodeToString([]byte(text))
	o, _ := New(secret)
	chosenText := []byte("foobar")
	cipherText := o.Encrypt(chosenText)
	a := append(o.Prefix, chosenText...)
	b := append(a, o.Secret...)
	manualCipherText, err := ecb.Encrypt(b, o.Key)
	if err != nil {
		t.Error(err)
	}
	if len(manualCipherText) != len(cipherText) {
		t.Errorf("oracle and manual ciphertext lengths don't match %d <> %d", len(manualCipherText), len(cipherText))
	}
	for i := 0; i < len(cipherText); i++ {
		if manualCipherText[i] != cipherText[i] {
			t.Errorf("expecting byte %d of the ciphertext to be %x, but got %x", i, manualCipherText[i], cipherText[i])
		}
	}
}
