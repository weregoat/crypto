package ch16

import (
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	"gitlab.com/weregoat/crypto/util"
	"strings"
	"testing"
)

func TestOracle(t *testing.T) {
	IWant := ";admin=true;"
	for i:=0; i < 2000; i++ {
		pre, err := util.RandomBytes(util.RandomInt(0, 10))
		if err != nil {
			t.Error(err)
		}
		post, err := util.RandomBytes(util.RandomInt(0, 10))
		if err != nil {
			t.Error(err)
		}

		o, err := New()
		if err != nil {
			t.Error(err)
		}
		plaintext := strings.Join(
			[]string{
				string(pre),
				IWant,
				string(post),
			},
			"",
		)
		oracleCipherText := o.Encrypt(plaintext)
		if o.IsAdmin(oracleCipherText) == true {
			t.Logf("Expecting failure on injecting plaintext, but got admin access with plaintext %+q", plaintext)
			t.Fail()
		}
		fakeCipher, err := cbc.Encrypt([]byte(plaintext), o.Key, o.IV)
		if err != nil {
			t.Error(err)
		}
		if o.IsAdmin(fakeCipher) == false {
			t.Log("ciphertext:")
			for _, block := range util.Split(fakeCipher, cbc.BlockSize) {
				t.Logf("%x", block)
			}
			t.Log("---")
			p, _ := cbc.Decrypt(fakeCipher, o.Key, o.IV)
			t.Logf("plaintext: %+q", p)
			t.Log("expecting the ciphertext to result in admin rights, but got nothing")
			t.Fail()
		}
	}
}
