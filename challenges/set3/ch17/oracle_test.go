package ch17

import (
	"encoding/base64"
	cbc "gitlab.com/weregoat/crypto/cbc/aes"
	"gitlab.com/weregoat/crypto/pkcs7"
	"gitlab.com/weregoat/crypto/util"
	"testing"
)

// Test with padded plaintexts
func TestOKPadding(t *testing.T) {
	var texts []string
	n := util.RandomInt(1,20)
	for i:=0; i < n; i++ {
		r, err := util.RandomBytes(util.RandomInt(0,60))
		if err != nil {
			t.Error(err)
			t.Fail()
			break
		}
		p := pkcs7.Pad(r, cbc.BlockSize)
		t.Logf("generated padded text: %+q", p)
		e := base64.StdEncoding.EncodeToString(p)
		texts = append(texts, e)
	}
	o, err := NewOracle(texts...)
	if err != nil {
		t.Error(err)
	}
	for j:=0; j < 100; j++ {
		c, _ := o.Encrypt()
		if err != nil {
			t.Error(err)
			t.Fail()
			break
		}
		result := o.CheckPadding(c)
		if err != nil {
			t.Error(err)
			t.Fail()
			break
		}
		if result != true {
			t.Errorf("expecting oracle to return good padding, got wrong one")
		}
	}
}

