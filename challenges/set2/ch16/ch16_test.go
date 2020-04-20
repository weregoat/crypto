package ch16

import (
	"testing"
)

func TestAttack(t *testing.T) {
	for i := 0; i < 500; i++ {
		o, err := New()
		if err != nil {
			t.Error(err)
		}

		poison := PoisonCipherText(o)
		if o.IsAdmin(poison) != true {
			t.Errorf("the attack didn't succeed")
		}
	}
}
