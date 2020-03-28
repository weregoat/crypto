package util

// Remove PKCS#7 padding
// https://en.wikipedia.org/wiki/Padding_(cryptography)#PKCS#5_and_PKCS#7
func RemovePadding(plainText []byte) []byte {
	// By default the last byte position is the end of the slice
	last := len(plainText) - 1
	// A slice smaller than 1 makes little sense
	if last < 0 {
		return plainText
	}
	// Pick the last byte value, it should the the number of bytes in the padding
	count := int(plainText[last])
	// It cannot be padded if the slice length is less than the padding length
	if len(plainText) < count {
		return plainText
	}
	// If the text is padded is should start at len - endByte
	padStart := len(plainText) - count

	// Pick count bytes from the end and they all should have count value
	for i := padStart; i < len(plainText); i++ {
		if int(plainText[i]) != count {
			return plainText
		}
	}
	return plainText[:padStart]
}
