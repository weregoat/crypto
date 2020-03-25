package util

import (
	"testing"
)

/*
Write a function to compute the edit distance/Hamming distance between two strings. The Hamming distance is just the number of differing bits. The distance between:
this is a test
and
wokka wokka!!!
is 37. Make sure your code agrees before you proceed.

[...]
No, that's not a mistake.
We get more tech support questions for this challenge than any of the other ones. We promise, there aren't any blatant errors in this text. In particular: the "wokka wokka!!!" edit distance really is 37.

 */
func TestHammingDistance(t *testing.T) {
	a := []byte("this is a test")
	b := []byte("wokka wokka!!!")
	expected := 37
	result, err := HammingDistance(a, b)
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Errorf("expecting a distance of %d, but got %d", expected, result)
	}

	a = []byte{1,2,3}
	b = []byte{1,2}
	result, err = HammingDistance(a, b)
	if err == nil {
		t.Errorf("expecting error from different slice lenghts, but got nothing")
	}

}


func BenchmarkHammingDistance(ben *testing.B) {
	a := []byte("this is a test")
	b := []byte("wokka wokka!!!")
	for i := 0; i < ben.N; i++ {
		HammingDistance(a, b)
	}
}