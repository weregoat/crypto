package util

import (
	"crypto/aes"
	"testing"
)

func TestDecryptAES128ECB(t *testing.T) {
	key := []byte("YELLOW SUBMARINE")
	plainText := "Yellow Submarine"
	cipher, err := aes.NewCipher(key)
	if err != nil {
		t.Error(err)
	}
	cipherText := make([]byte, len(plainText))
	cipher.Encrypt(cipherText, []byte(plainText))
	t.Logf("%q encrypted to %q with key %q", plainText, EncodeToBase64(cipherText), key)
	decrypted, err := DecryptAES128ECB(cipherText, key)
	if err != nil {
		t.Error(err)
	}

	if string(decrypted) != plainText {
		t.Errorf("expecting %q, got %q", plainText, decrypted)
	}

	shortKey := key[:12]
	_, err = DecryptAES128ECB(cipherText, shortKey)
	if err == nil {
		t.Errorf("expecting error because of short key, but got nothing")
	}

	shortText := cipherText[:12]
	_, err = DecryptAES128ECB(shortText, key)
	if err == nil {
		t.Errorf("expecting error because of short cipherText, but got nothing")
	}
}
