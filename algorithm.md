## Introduction

imohash is a file hashing algorithm optimized for large files. It uses
file size and sampling in hash generation. Because it does not process
the whole file, it is not a general purpose hashing algorithm. But for
applications where a hash sample is sufficient, imohash will provide a
high performance hashing, especially for large files over slow
networks.

## Algorithm

imohash generates a 128-bit hash from a fixed length message or file.
This is done in two phases:

1. hash calculation
2. size injection

### Parameters and mode

imohash takes two parameters, as well as the message length:

* sample size (s)
* sampling threshold (t)
* message length (L)

There are two mode of operation: **sampled** and **full**. Mode is
determined as follows:

```
if (s > 0) && (t > 0) && (L > t) && (t > 2s) 
  mode = sampled
else
  mode = full
```

### Hash calculation

The core hashing routine uses [MurmurHash3](https://code.google.com/p/smhasher/wiki/MurmurHash3) in a 128-bit configuration.
Hashing in *Full* mode is identical to passing the entire
message to Murmhash3.  *Sampled* mode constructs a new message using
three samples from the original:

Message M of length L is an array of bytes, M[0]...M[L-1]. If
L > t, full mode is used and h'=Murmur3(M). Otherwise, samples are selected and concatenated as follows:

```
middle = floor(L/2)
S0 = M[0:s-1]           // samples are s bytes long
S1 = M[middle:middle+s]
S2 = M[L-s:L-1]

h' = Murmur3(concat(S0, S1, S2))
```
### Size injection

Size is inserted into the hash directly. This means that two files
that differ in size are guaranteed to have different hashes.

The message size is converted to a variable-length integer (varint)
using 128-bit encoding. Consult [Google Protobuf documentation](https://developers.google.com/protocol-buffers/docs/encoding#varints) for more
information on the technique.

The result of encoding will be an array **v** of 1 or more bytes. This
array will replace the highest-order bytes of h.

```
h = concat(v, h'[len(v):])
```

h is the final imosum hash.

## Default parameters

The default imohash parameters are:

s = 16384  
t = 131072

t was chosen to delay sampling until file size was outside the range
of "small" files, such as text files that might be hand-edited and
escape both size changes and being detect by sampling. s was chosen to
provide a large enough sample to distiguish files of like size, but
still small enough to provide high performance.

An application should adjust these values as necessary.

## Test Vectors

(Note: these have not been independently verified using another implementation.)

To avoid offset errors in testing, the test messages needs to not repeat
trivially. To this end, MurmurHash3 is used to generate pseudorandom test
data (which also provides a convenient verification of the implementation
of MurmurHash3). The message generation uses the 128-bit variant of
MurmurHash3 to add 16 bytes at a time. M(n) shall be a test data n bytes
long:

```
M(n):
   msg = []
   while len(msg) < n:
       MurmurHash3.Write('A')
       msg = msg + MurmurHash3.Sum()
   return msg[0:n]

// M(16)      ==     035fc2b79a29b17a387df29c46dd9937
// M(1000000) == ... 7b7cd46489e93605eb7894c5d338463f
```

Test vectors for imohash of M length n using sample size s and sample
threshold t.

```
  s       t     M(n)        I
{16384, 131072, 0, "00000000000000000000000000000000"},
{16384, 131072, 1, "016ac6dd306a3e594e711127c5b5a8e4"},
{16384, 131072, 127, "7f0e0eaa0a415e2878303214a087fbab"},
{16384, 131072, 128, "80017902d85ab752459758292a217ed8"},
{16384, 131072, 4095, "ff1fc7e3b9890469ee676df99f6ac7b7"},
{16384, 131072, 4096, "80203d4ea589349dcdf63511a8d49d4a"},
{16384, 131072, 131072, "808008ed886018d5cd37a9b35bbae286"},
{16384, 131073, 131072, "808008ac6ab33ddcaadea68ba5ad4f05"},
{16384, 131072, 500000, "a0c21ea46738cf366d496962000a45f7"},
```




