# imohash
Fast hash for large files

imohash is a fast, constant time hash for files with the following properties:

- files of different length will have different hashes
- files are sampled, providing fast fixed-time performance
- the underlying hash is murmur3

# Background

imohash fell out of a need to do file synchronization and deduplication overa fairly slow network. Managing media (photos and video) over wi-fi between a NAS and multiple family computers is a typical example. To check whether two files are the same without doing a byto-for-byte comparison, a some properties to check include:

1. filename
2. date/time
3. size
4. hash of content

I tend to avoid using numbers 1 & 2 for very much. size is an simple way to tell that files aren't the same, and using a decent hash function in addition to size can shown that two files are most likely the same. But hashing a file still involves reading the file, and this can be very slow when you're talking about large files. imohash attempts to address the problem by working from the following assumption: **hashing a small part of a file may be good enough**.

If you you a library of tens of thousands of photos/video, chances are that there aren't hundreds of files with the same file size. So right away you can show what things *aren't* the same. And of items that are the same size, there's probably a very good chance that comparing or hashing a few KB of the beginning, middle and end will be sufficient to determine which files are likely the same.

# Design

imohash is a small wrapper around murmur3, a fast and well-regarded hashing algorithm.

# Performance

# Credits
The "sparseFingerprints" used in [TMSU](https://github.com/oniony/TMSU) gave me some confidence in this approach to hashing.
