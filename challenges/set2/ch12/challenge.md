#Byte-at-a-time ECB decryption (Simple)
Copy your oracle function to a new function that encrypts buffers under ECB mode using a consistent but unknown key (for instance, assign a single random key, once, to a global variable).

Now take that same function and have it append to the plaintext, BEFORE ENCRYPTING, the following string:
```
Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK
```
---
###Spoiler alert.
Do not decode this string now. Don't do it.

---

Base64 decode the string before appending it. Do not base64 decode the string by hand; make your code do it. The point is that you don't know its contents.

What you have now is a function that produces:
```
AES-128-ECB(your-string || unknown-string, random-key)
```
It turns out: you can decrypt "unknown-string" with repeated calls to the oracle function!

Here's roughly how:

1. Feed identical bytes of your-string to the function 1 at a time --- start with 1 byte ("A"), then "AA", then "AAA" and so on. Discover the block size of the cipher. You know it, but do this step anyway.
2. Detect that the function is using ECB. You already know, but do this step anyways.
3. Knowing the block size, craft an input block that is exactly 1 byte short (for instance, if the block size is 8 bytes, make "AAAAAAA"). Think about what the oracle function is going to put in that last byte position.
4. Make a dictionary of every possible last byte by feeding different strings to the oracle; for instance, "AAAAAAAA", "AAAAAAAB", "AAAAAAAC", remembering the first block of each invocation.
5. Match the output of the one-byte-short input to one of the entries in your dictionary. You've now discovered the first byte of unknown-string.
6. Repeat for the next byte.
---
###Congratulations.
This is the first challenge we've given you whose solution will break real crypto. Lots of people know that when you encrypt something in ECB mode, you can see penguins through it. Not so many of them can decrypt the contents of those ciphertexts, and now you can. If our experience is any guideline, this attack will get you code execution in security tests about once a year.

## Goat notes
I should have seen it coming from the previous challenge: [Wikipedia](https://en.wikipedia.org/wiki/Chosen-plaintext_attack#Different_forms).

Okay, I now understand how it's supposed to work. And all my work today for moving the oracle code to a separate package has been wasted.

Something I did knew, but that this exercise reminded me: the attacks are 
somewhat specific to the implementation, not the cypher or mode in general.
In other words, this attack is against the particular implementation in the oracle (combining the plaintexts) **through** a weakness in the ECB mode. 

In fact the oracle in challenge #13, for example, is not vulnerable to this attack (although is still vulnerable; or so I suppose, as I haven't solved it yet).

The way I figure it to work is a bit like this:
First. If we submit **one** block X bytes short, the oracle will fill it up with the fist X bytes of its plaintext.
We can also submit a series of full size blocks, each with a different 16th byte and put the resulting ciphertext block in a table so we can know, by the comparing the above with the ones in the table, the byte it was that it was pushed from the plaintext.
 
Once we discover the first plaintext byte this become part of the **chosen** plaintext (we drop the leftmost byte to make space) for the lookup table and we force the oracle to give us the first two bytes of the plaintext (of which one is known) so we can compare the ciphertext with the table (which was created by submitting to the oracle chosen texts composed by 14 chosen bytes + the first byte of the plaintext and varying byte (from 0-256, 0-127, or a more restricted choice).
Etcetera until we have the first block of plaintext.
At that point we have the plaintext of the first block so we can shift it left one byte, build a table with all the ciphertexts of first plaintext block (minus the leftmost) and all the possible 16th bytes. 
Force the oracle to shift the plaintext to the left again one byte and we can compare the **second** block of the ciphertext with the ones in the table to find the first byte of the **second** block of the plaintext. And so forth...

Gosh, that was messy to explain. I hope my future self is brighter than this one and can make it out.

Notice that we could as well, do this two bytes at a time, just that the table will have 256*256=65536 entries, and 16777216 for three bytes etc...