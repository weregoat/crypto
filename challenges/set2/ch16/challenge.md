#CBC bitflipping attacks
Generate a random AES key.

Combine your padding code and CBC code to write two functions.

The first function should take an arbitrary input string, prepend the string:

```
"comment1=cooking%20MCs;userdata="
```

.. and append the string:

```
";comment2=%20like%20a%20pound%20of%20bacon"
```

The function should quote out the ";" and "=" characters.

The function should then pad out the input to the 16-byte AES block length and encrypt it under the random AES key.

The second function should decrypt the string and look for the characters ";admin=true;" (or, equivalently, decrypt, split the string on ";", convert each resulting string into 2-tuples, and look for the "admin" tuple).

Return true or false based on whether the string exists.

If you've written the first function properly, it should not be possible to provide user input to it that will generate the string the second function is looking for. We'll have to break the crypto to do that.

Instead, modify the ciphertext (without knowledge of the AES key) to accomplish this.

You're relying on the fact that in CBC mode, a 1-bit error in a ciphertext block:

Completely scrambles the block the error occurs in
Produces the identical 1-bit error(/edit) in the next ciphertext block.

Stop and think for a second.
---
Before you implement this attack, answer this question: why does CBC mode have this property?

## Goat notes
Again, I wish the text was more clear.
So, I need to lay out my understanding of the challenge a moment.
We control part of the plaintext. 
The issue is that whatever we pass as this plaintext, ";" and "=" are quoted out, so we cannot just pass ";admin=true" and then pass the resulting ciphertext to the second function.

The 1-bit error in the challenge is confusing, because we are not going to just change one bit (although one may think [bit-flipping attack](https://en.wikipedia.org/wiki/Bit-flipping_attack) is just that).

So, IMHO, the attack to this specific implementation should work this way:

* We have two or more blocks of plaintext (P1, P2) that we control.
Sadly the first function sanitises P1 and P2, so we cannot just inject the string
we want.

* We have two blocks of ciphertext (C1, C2) that we also control to a small extent resulting from the encryption of P1 and P2 (we can change them, but not decipher them as we have no key)
* We know that CBC XOR C1 with the result of the decryption of C2 (E2) to get P2 (the 1-bit error/edit part). 
* Given the above, having C1 and P2, we can get E2: C1^P2=E2
* Having E2 we could inject a C1' resulting from E2^P2 (which is derived in the oracle from C2 and doesn't change if we change C1) so that when XOR with E2 will result in a version of P2 (P2'?) with text we want.
* Of course P1' will not be same as P1, but, I assume we don't care.

Easier to follow with [CBC scheme](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Cipher_block_chaining_(CBC\)) in front. 

It would also be possible to just flip targeted bytes of the ciphertext to alter the resulting plaintext only in part, but this, I think, is simpler to put in place (and, in a way, more general as it works regardless of the way the plaintext is sanitised).

