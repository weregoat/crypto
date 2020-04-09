
#[Byte-at-a-time ECB decryption (Harder)](https://cryptopals.com/sets/2/challenges/14)
Take your oracle function from #12. Now generate a random count of random bytes and prepend this string to every plaintext. You are now doing:
```
AES-128-ECB(random-prefix || attacker-controlled || target-bytes, random-key)
```
Same goal: decrypt the target-bytes.

Stop and think for a second.
---
What's harder than challenge #12 about doing this? How would you overcome that obstacle? The hint is: you're using all the tools you already have; no crazy math is required.

Think "STIMULUS" and "RESPONSE".

## Goat notes
It seems to me that the way to solve this is to lead it back to a case #12.
I am thinking that all I need to know is the length of the prefix (which I assume is constant as the key and not varying with every encryption).
I can guess the block length the same way as with #12, so I will skip that.

So, if prepare a plaintext of two identical blocks, and then prepend to it an increasing number of 'A' I can guess the length of the prefix the moment I can detect two identical, consecutive, cipherblocks roughly at the beginning of the ciphertext.
Then, as I see it, it's just a matter of proceeding as with #12 with the new chosen plaintext.

It seems too easy. I need to think a bit and see if I am missing something obvious.