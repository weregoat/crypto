#Detect single-character XOR
One of the 60-character strings in this [file](https://cryptopals.com/static/challenge-data/4.txt) has been encrypted by single-character XOR.

Find it.

(Your code from #3 should help.)

## Goat notes
The code in #3 helps **only** if you include the space character in the frequency map.
Otherwise is `Ieeacdm*GI-y*fcao*k*ze\u007fdn*el*hkied"` and not `Cooking MC's like a pound of bacon`.

I tried to find a frequency list that included punctuation, but it was quite hard.
I found one [here](https://web.archive.org/web/20040603075055/http://www.data-compression.com/english.html), through the Wikipedia article on [Frequency Analysis](https://en.wikipedia.org/wiki/Frequency_analysis).

In the end, I wrote a simple script to extract the frequencies from a corpus of 
English classics from [Project Gutenberg](https://www.gutenberg.org/browse/scores/top#books-last30)
Interestingly enough the order is not really the same as on Wikipedia, but it should be enough for my scope (and including some punctuation may help with these challenges).

I guess that, even if I had a Twitter account, I would not be allowed make jokes about "ETAOIN SHRDLU".