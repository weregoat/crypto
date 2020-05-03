package ch17

import (
	"testing"
)

func TestAttack(t *testing.T) {
	/*
	As the attack is implemented now, it doesn't always correctly decrypt the last
	block.
	 */
	for i:=0; i < 1000; i++ {
		oracle, err := NewOracle()
		if err != nil {
			t.Error(err)
			break
		}
		cipherText, iv := oracle.Encrypt()
		plainText := string(attackOracle(oracle, cipherText, iv))
		// The first block *should* always succeed
		if plainText != oracle.Plaintext {
			t.Errorf("Expecting %+q, but got %+q",
				oracle.Plaintext,
				plainText,
			)
			t.Fail()
			break
		}
	}
}
