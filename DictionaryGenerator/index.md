# Dictionary Generator

## Design

Dictionary Generator is designed to create a tire tree(dictionary tree) that mapping every english word to
 another word or null. The mapping policy is below: if the word is a stop word, it will be mapped to null,
else it will be mapped to its stem.

## Reference

> Go-stem https://github.com/agonopol/go-stem

> stopwords.txt http://www.onjava.com/2003/01/15/examples/EnglishStopWords.txt

## Implement

By assessment, there is only 524 stop words. I can easily push them all in a tire tree in memory. Final output
 interface is a function map word to word as `Design` says.

## Test

I did some unit test, and get 85% coverage.