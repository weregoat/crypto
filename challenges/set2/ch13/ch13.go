package ch13

import (
	"gitlab.com/weregoat/crypto/pkcs7"
	"log"
)

/*
	I am going to cheat a bit here.
I assume I know already it's ECB, it's 16byte blocks, that the _role_ is at the end of the encoded string,
and the accepted input (i.e. no "=" or "&").
 */

/* First we get the length the email address should be to have "user" in the last block"
	We cannot use "role=user" because we know the email can't have "=" or "&" in it.
We assume the structure is email=[our input]somethingsomething=user
So, we pass longer and longer mails until we get an extra block (which should be encoding "r"+padding) and we add 4 bytes ("user").
It's true that we are assuming the somethingsomething remains constant, but it's quite easy to adapt to a scenario where
the uid increases.
*/

func GetEmailSuffix(o Oracle, email string, blockSize int) []byte {
	e := []byte(email)
	baselen := len(o.Encrypt(email))
	for i:=0; i <= blockSize; i++ { // If we had to add more than 16 bytes, we got the block wrong
		cipherLen := len(o.Encrypt(string(e)))
		if cipherLen > baselen {
			break
		}
		e = append([]byte{'a'}, e...)
	}
	// Remember we need 3 more bytes (user)
	e = append([]byte("eeee"), e...)
	return e
}

func CraftCiphertext(o Oracle, target string, blocksize int) []byte {
	if len(target) >= blocksize {
		log.Fatal("code only works with admin roles shorter than 16 bytes")
	} // Code only works with one block of ciphertext (it can be expanded, of course)
	// Get the email we need to have the "user" encrypted in the last block
	email := GetEmailSuffix(o, "foo@mail.com", blocksize)
	// Get the cipherText for this email
	userCipher := o.Encrypt(string(email))
	// Remove the last block where "user" is
	noUser := userCipher[0:len(userCipher)-blocksize]
	// "email=" is 6 bytes that the oracle add, so we need to
	// insert the "admin" block where it will be on its own block
	// i.e. after byte 10 (16-6)
	beforeAdmin := email[0:10]
	afterAdmin := email[10:]
	// Pad the admin to fill up a full block
	admin := pkcs7.Pad([]byte(target), blocksize)
	// Put the email address together 10bytes + adminblock + rest
	admin = append(admin, afterAdmin...)
	email = append(beforeAdmin, admin...)
	// Get the poisoned ciphertext
	poisoned := o.Encrypt(string(email))
	// The block of the ciphertext with the admin text should be the second one
	adminBlock := poisoned[16:32]
	// Attach this block to the cipherText without the user part
	cipherText := append(noUser, adminBlock...)
	// We can (or should be able to) use this cipherText in the cookie to get
	// admin privileges...
	return cipherText
}