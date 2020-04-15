#PKCS#7 padding validation
Write a function that takes a plaintext, determines if it has valid PKCS#7 padding, and strips the padding off.

The string:

```
"ICE ICE BABY\x04\x04\x04\x04"
```
... has valid padding, and produces the result "ICE ICE BABY".

The string:
```
"ICE ICE BABY\x05\x05\x05\x05"
```
... does not have valid padding, nor does:

```
"ICE ICE BABY\x01\x02\x03\x04"
```
If you are writing in a language with exceptions, like Python or Ruby, make your function throw an exception on bad padding.

Crypto nerds know where we're going with this. Bear with us.

## Goat notes
I am not sure I like the error part. I find it too harsh, in encryption/decryption cases. If the text was not padded properly, or with PKCS#7, I can still decrypt it and figure out the garbage at the end later.

I think I will prefer a separate, specific function that check if the plaintext is padded or not according to PKCS#7.
