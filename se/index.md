# MiniSearchEngine

## CreateIndex

The invert index include two parts index file and the invert file. The invert file is a sequence of <rawFId, rawFOff> pairs following
 a 8 bytes integer `length` that means the number of pairs. The sequence means the word appears in raw file #rawFId, at offset #rawFOff.
The index file is a sequence of <hashCode, invOff> pairs, which means the appear records of the word with the hashCode is at offset
#invOff of its corresponding invert file.

The invert files will be save to `./invert` with a number between 0-255 as name, e.g. `./invert/1.inv`. The number `n` represents
 that all words' hash in this index has a first byte `n`. For example, a word `TTT`'s hash code is `0x15 0x88 0x96 0x23`, it will
 be in index file `./index/21.inv` and invert file `./invert/21.inv`.

> By default, the program is working in one process, this is because it is designed to process big data. If a raw file is 4GB large,
its invert index can hardly be entirely put in main memory. If use multi-process in this case, the main memory will be lack and the
 operation system will continue swapping page in and out, which will loss performance. But if you only want to process some little
 files, you can run multiple processes of this program using its command line parameters.

## Search

A query word will be hashed and stop-words handled before search it from the index file. To improve the search speed, there's a cache
, implemented by Splay tree in main memory, which can contains 1024 word's result cache.

## Implement

The file IO is all buffered to reduce system call cost and using memory-mapped IO. Sequences of <hashCode, invOff> is from big to little.
 The search of index is cached and use binary search if not cached.

## Feature

This program can run in 4G RAM machine to create index of about **1,280,000,000,000** input words, and has a well performance for about
**512,000,000,000** input words.