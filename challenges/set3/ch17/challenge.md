#The CBC padding oracle
This is the best-known attack on modern block-cipher cryptography.

Combine your padding code and your CBC code to write two functions.

The first function should select at random one of the following 10 strings:
```
MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=
MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=
MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==
MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==
MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl
MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==
MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==
MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=
MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=
MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93
```
... generate a random AES key (which it should save for all future encryptions), pad the string out to the 16-byte AES block size and CBC-encrypt it under that key, providing the caller the ciphertext and IV.

The second function should consume the ciphertext produced by the first function, decrypt it, check its padding, and return true or false depending on whether the padding is valid.

##What you're doing here.

---
This pair of functions approximates AES-CBC encryption as its deployed serverside in web applications; the second function models the server's consumption of an encrypted session token, as if it was a cookie.

---

It turns out that it's possible to decrypt the ciphertexts provided by the first function.

The decryption here depends on a side-channel leak by the decryption function. The leak is the error message that the padding is valid or not.

You can find 100 web pages on how this attack works, so I won't re-explain it. What I'll say is this:

The fundamental insight behind this attack is that the byte 01h is valid padding, and occur in 1/256 trials of "randomized" plaintexts produced by decrypting a tampered ciphertext.

02h in isolation is not valid padding.

02h 02h is valid padding, but is much less likely to occur randomly than 01h.

03h 03h 03h is even less likely.

So you can assume that if you corrupt a decryption AND it had valid padding, you know what that padding byte is.

It is easy to get tripped up on the fact that CBC plaintexts are "padded". Padding oracles have nothing to do with the actual padding on a CBC plaintext. It's an attack that targets a specific bit of code that handles decryption. You can mount a padding oracle on any CBC block, whether it's padded or not.

## Goat notes
Here is how I understand it to work:
* We know how PKCS#7 padding works. 0x01 is attached for one byte padding 0x02, 0x02 for two etc...
* We control the ciphertext we submit (we even have the IV so we can manipulate the first block too).
* The oracle will tell us if the decrypted ciphertext is padded correctly. 
* We know, from previous attacks, that we can XOR the cyphertext so that we
  can, to a certain extent, control the resulting plaintext.
  
Now, the basic idea is that the oracle would tell us if the last byte of the resulting plaintext is 0x01, or the last two are 0x02, 0x02 etc.
  
Given the way CBC works, and with (Ciphertext, Decrypted, Plaintext) is like this:
C1 => D1^IV => P1
C2 => D2^C1 => P2

If we change any byte of C1, that will result in a change in the same byte of P2 (regardless of the decryption, because it's done on C2).

The idea is to modify C1 in a way that we end up with a plaintext P2' we know; since the oracle tells us if the last byte is 0x01 (or 0x02 and 0x02, or 0x03 0x03 0x03 etc) we exploit that.

Something like this:
We start by changing the last byte (15) of C1.

We submit the ciphertexts (I1 and C2) to the oracle.
If it returns true it's possibly because the last byte of the resulting plaintext (P2') is 0x01 (or that the last two bytes are 0x02 and 0x02, or the last three 0x03 0x03 0x03; even less probable).
If it returns false, well the last byte of P2' is not 0x01

We keep trying all the 256 possible values for I1[15] until we get P2'[15] as 0x01.

We can then calculate D2[15] with I1[15]^P2'[15] (because P2^C1 => D2)

We now have D2[15] and C1[15], so we can get the P2[15] by XORing them: D2[15]^C1[15] => P2[15]

The next step is to try to get the last **two** bytes of P2' to be 0x02 and 0x02 and so forth.

We know D2[15] so we can calculate I1[15] to result in P2'[15] == 0x02.
i.e. D2[15]^0x02 => C1'[15]

And we repeat the procedure to extract a P2'[14] == 0x02 by manipulating I1[14] and checking with the oracle.

We keep going until the have P2, and then we can do the same for P1, with IV instead of C1 and C1 as C2...

