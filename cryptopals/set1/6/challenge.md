#Break repeating-key XOR
###It is officially on, now.

This challenge isn't conceptually hard, but it involves actual error-prone coding. The other challenges in this set are there to bring you up to speed. This one is there to qualify you. If you can do this one, you're probably just fine up to Set 6.

---

[There's a file here](https://cryptopals.com/static/challenge-data/6.txt). It's been base64'd after being encrypted with repeating-key XOR.

Decrypt it.

Here's how:

1. Let KEYSIZE be the guessed length of the key; try values from 2 to (say) 40.
2. Write a function to compute the edit distance/Hamming distance between two strings. The Hamming distance is just the number of differing bits. The distance between:
```
this is a test
```
and
```
wokka wokka!!!
```
is 37. _Make sure your code agrees before you proceed_.
3. For each KEYSIZE, take the first KEYSIZE worth of bytes, and the second KEYSIZE worth of bytes, and find the edit distance between them. Normalize this result by dividing by KEYSIZE.
4. The KEYSIZE with the smallest normalized edit distance is probably the key. You could proceed perhaps with the smallest 2-3 KEYSIZE values. Or take 4 KEYSIZE blocks instead of 2 and average the distances.
5. Now that you probably know the KEYSIZE: break the ciphertext into blocks of KEYSIZE length.
6. Now transpose the blocks: make a block that is the first byte of every block, and a block that is the second byte of every block, and so on.
7. Solve each block as if it was single-character XOR. You already have code to do this.
8. For each block, the single-byte XOR key that produces the best looking histogram is the repeating-key XOR key byte for that block. Put them together and you have the key.

This code is going to turn out to be surprisingly useful later on. Breaking repeating-key XOR ("Vigenere") statistically is obviously an academic exercise, a "Crypto 101" thing. But more people "know how" to break it than can actually break it, and a similar technique breaks something much more important.

---
####No, that's not a mistake.
We get more tech support questions for this challenge than any of the other ones. We promise, there aren't any blatant errors in this text. In particular: the "wokka wokka!!!" edit distance really is 37.

---
###Goat notes
This was very confusing, because, step 3, as I read it, is wrong. Just picking the first two blocks doesn't work; you need to calculate the Hamming distance for more blocks (I didn't try the second algorithm yet). How many blocks? I have no idea. 

Nor I have a good idea why this method works; it's not in the suggested methods in the [Wikipedia article](https://en.wikipedia.org/wiki/Vigenère_cipher#Cryptanalysis).
I found [this](https://crypto.stackexchange.com/questions/8115/repeating-key-xor-and-hamming-distance) may be an explanation.
Also, [this person](https://trustedsignal.blogspot.com/2015/06/xord-play-normalized-hamming-distance.html) ran some tests, and the result seems to be that is not such a general solution (I didn't really check the numbers).
Another issue, imagine we were using UTF-16 or some other encoding, if the above explanation is correct this method would not work very well...
Other links with similar explanations:
* https://carterbancroft.com/breaking-repeating-key-xor-theory/


I do wonder how they came out with it. 