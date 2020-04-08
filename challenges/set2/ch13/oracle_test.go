package ch13

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := map[string]map[string]string{
		"foo=bar&baz=qux&zap=zazzle":
		{
			"foo": "bar",
			"baz": "qux",
			"zap": "zazzle",
		},
		"foo=bar":
		{
			"foo": "bar",
		},
		"foo&bar":
		{},
	}
	for src, data := range tests {
		for key, value := range parse(src) {
			expected, present := data[key]
			if ! present {
				t.Errorf("expecting key %+q to be present in parsed string %+q, but it was not", key, src)
			}
			if expected != value {
				t.Errorf("expecting value %+q for key %+q in parsed string %+q, but got %+q", expected, key, src, value)
			}
		}
	}

}

func TestOracle(t *testing.T) {
	email := "email@test.com"
	var expected = map[string]string {
		"email":email,
		"role":"user",
		"uid":"10",
	}
	o, err := New()
	if err != nil {
		t.Error(err)
	}
	c := o.Encrypt(email)
	s := o.Decrypt(string(c))
	if len(s) != len(expected) {
		t.Errorf("wrong map size %d <> %d", len(s), len(expected))
	}
	for k,v := range expected {
		i,ok := s[k]
		if !ok {
			t.Errorf("missing key %+q from decrypted map %v", k, s)
		}
		if i != v {
			t.Errorf("expecting value %+q for key %+q, got %+q", v,k,i)
		}
	}
}

func TestAttack(t *testing.T) {
	blockSize := 16
	var targets = []string{"admin", "super", "superuser", "master", "root"} // We want admin roles
	o, err := New()
	if err != nil {
		t.Error(err)
	}
	for _,target := range targets {
		cipherText := CraftCiphertext(o, target, blockSize)
		data := o.Decrypt(string(cipherText))
		role := data["role"]
		if role != target {
			t.Errorf("expecting role to be %+q, got %+q", target, role)
		}
	}
}