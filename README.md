# A Golang library for Bitcoin hash puzzles

Bitcoin hash puzzles arise from a combination of the original [Bitcoin]() protocol
and the [Stratum](github.com/DanielKrawisz/go-Stratum) protocol, which is the way that miners communicate with
mining pools, aka nodes. This library provides data types for dealing with hash
puzzles.

dependencies: libsv/go-bc

The goal of a hash puzzle is to create a bitcoin block header, which is an 80-byte
data structure that includes a target such that, when hashed with Hash256 (aka
double SHA256), the hash digest, when read as a little-endian number, should be
less than the given target. This is done by providing a 32 bit nonce of arbitrary
data. A timestamp must also be provided.

The block header establishes that a list of transactions belongs in the block
through the Merkle root.

32 bits is not enough arbitrary data to ensure that a given hash puzzle has a
solution. However, a wider space of solutions can be explored via other means.
The timestamp increments every second and the Merkle root changes as more
transactions are added to the block. Finally, the coinbase input script can have
up to 100 bytes of arbitrary data, which is enough for any hash puzzle over a
32 byte hash digest space.  

Stratum establishes additional properties of the hash puzzle. The mining pool
creates the coinbase transaction other than 96 bits of data that will be in the
input script. 32 of these is a user id that the pool assigns to the miner and
the other 64 can be used by the miner.

[ASIC Boost](http://www.math.rwth-aachen.de/~Timo.Hanke/AsicBoostWhitepaperrev5.pdf)
is a strategy for arriving faster at a given difficulty target by
caching certain data. It has to do with the details of SHA256. Because miners were
messing around with the timestamp in order to use ASIC Boost, which is bad for
Bitcoin's functionality as a timestamp server,
[BIP 320](https://github.com/bitcoin/bips/blob/master/bip-0320.mediawiki) was
designed to address the problem by assigning 16 bits of the version field of the
block header to be used as additional nonce data that the miner provides to the
mining pool. Thus, there 48 total bits of arbitrary data in the block header.
