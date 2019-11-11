# salt

A simple command-line utility written in [Go](https://golang.org) that supports
the box (_public key cryptography_) and secretbox (_secret key cryptography_)
as described by [NaCL](https://en.wikipedia.org/wiki/NaCl_(software)) and
part of the Golang standard library [box](https://godoc.org/golang.org/x/crypto/nacl/box)
and [secretbox](https://godoc.org/golang.org/x/crypto/nacl/secretbox).

`salt` supports reading the encryption key securely from a prompt or passed
via the `-k` option. You may use other UNIX tools such as `base64` to encode
the resulting cipher text for portability.

## Installation

```#!bash
$ go get -u github.com/prologic/salt
```

## Usage

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

## License

MIT
