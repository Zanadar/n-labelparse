# Name Compression

DNS messages have a lot of names in them and they can get repetitive! To reduce the size of a given message, the DNS RFC defines a special encoding compressed with "pointers" is used (detailed [below](#name-encoding)).

As an exercise, we would like you to implement DNS label parsing such that it can interpret and expand the names.

**Write a program, which when given a binary file and a byte offset, returns the decompressed name at that offset.**

A binary file, [response.bin](response.bin), is included.

Using the following offsets on the file should produce these outputs:

```shell
> ./labelparse response.bin 20
IS.AWE.SOME
> ./labelparse response.bin 40
NS1.IS.AWE.SOME
> ./labelparse response.bin 64
SOME
```

The program can be implemented in the language of your choice; it should not use DNS parsing libraries. It should demonstrate best practices.

Please do not spend more than four hours on this exercise.

## name encoding

Uncompressed, each dot-delimited part of a name -- a `label` -- is encoded with a length byte and then the number of octets; finally, there is a null byte for the implicit root element. For example, `ns1` -- the first label -- has a length byte of 3, while it's followed by bytes 110 (ASCII "n") , 115 ("s"), and 49 ("1"). The whole name `ns1.com.` is encoded as `[3, 110, 115, 49, 3, 99, 111, 109, 0]` or `[3, "n", "s", "1", 3, "c", "o", "m", 0]`.

To reduce message size, DNS uses two-byte label "pointers" to compress names. This means that in fact, when parsing a name, the first byte of a label can either be its length _or_ the start of a pointer to another label in the message. If it's a pointer, the first two bits are set, and the remaining bits are the byte offset of the label in the message; note that pointers can point to names containing more pointers, but there will never be any labels _after_ an expanded pointer.

That means that while the first occurrence of `ns1.com.` at byte 6 might be encoded as `[3, "n", "s", "1", 3, "c", "o", "m", 0]`, the _next_ occurrence will be a pointer to the first occurrence -- i.e. it'll be a pointer to the sixth byte -- or, in bytes, `[192, 6]` or, in bits, `1100000000000110`.  If `foo.ns1.com.` were also in the message, it could be encoded as `[3, "f", "o", "o", 192, 6]` and so on.
