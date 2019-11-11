# salt

A simple command-line utility written in [Go](https://golang.org) that supports
the box (_public key cryptography_) and secretbox (_secret key cryptography_)
as described by [NaCL](https://en.wikipedia.org/wiki/NaCl_(software)) and
part of the Golang standard library [box](https://godoc.org/golang.org/x/crypto/nacl/box)
and [secretbox](https://godoc.org/golang.org/x/crypto/nacl/secretbox).

## Installation

```#!bash
$ go get -u github.com/prologic/salt
```

## Usage

### Secret Box (Symetric Crypto)

Encryption:

```#!bash
$ echo 'hello world' | ./salt -k foobar -e
�+TQI?�B�E�LvM��4�)3��C��*C63�ۂ�~=%��'
```

Decryption:

```#!bash
$ echo 'hello world' | ./salt -k foobar -e > hello.enc
$ ./salt -k foobar -d < hello.enc
hello world
```

### Box (Asymetric Crypto)

Generate Keys:

```#!sh
$ ./salt -g
Private Key: 9B07A5664EDEE31FB1B5F456CA552BC4CD8474345A6185A8B9CB51EC7B0CC65E
Public Key: 96EBFD5C8092D5CBC420990159079F72DAF5B88AE54C25432D027CCB9E7AD85E

$ ./salt -g
Private Key: DA50675E53991113CB271664099E87E0A9251E31189A4B5BA8B55F15674799E7
Public Key: 279E5C728B4F0D851DA4A98E1DCF725834CC6B38B3E00C1A14DFE1B926730835
```

Encrypt and Sign:

```#!sh
$ ./salt -k 9B07A5664EDEE31FB1B5F456CA552BC4CD8474345A6185A8B9CB51EC7B0CC65E -p 279E5C728B4F0D851DA4A98E1DCF725834CC6B38B3E00C1A14DFE1B926730835 -e > msg.enc
hello world
^D
```

Decrypt and Verify:

```#!sh
$ ./salt -k DA50675E53991113CB271664099E87E0A9251E31189A4B5BA8B55F15674799E7 -p 96EBFD5C8092D5CBC420990159079F72DAF5B88AE54C25432D027CCB9E7AD85E -d < msg.enc
hello world
```

## License

MIT
