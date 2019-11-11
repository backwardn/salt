package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	enc bool
	dec bool
	key string

	errDecryptionError = errors.New("error: decryption failed")
)

func init() {
	flag.BoolVar(&enc, "e", false, "encrypt contents of message")
	flag.BoolVar(&dec, "d", false, "decrypt contents of message")
	flag.StringVar(&key, "k", "", "secret key")

	flag.Parse()
}

func encrypt(key, message []byte) ([]byte, error) {
	var (
		secretKey [32]byte
		nonce     [24]byte
	)

	copy(secretKey[:], key)

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	ciphertext := secretbox.Seal(nonce[:], message, &nonce, &secretKey)
	return ciphertext, nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	var (
		decryptNonce [24]byte
		secretKey    [32]byte
	)

	copy(secretKey[:], key)

	copy(decryptNonce[:], ciphertext[:24])
	plaintext, ok := secretbox.Open(nil, ciphertext[24:], &decryptNonce, &secretKey)
	if !ok {
		return nil, errDecryptionError
	}
	return plaintext, nil
}

func main() {
	var (
		input []byte
		err   error
	)

	if !enc && !dec {
		fmt.Printf("Usage: %s [options] [message]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(flag.Args()) == 1 {
		input = []byte(flag.Arg(0))
	} else {
		input, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading message from stdin: %s", err)
		}
	}

	if enc {
		ciphertext, err := encrypt([]byte(key), input)
		if err != nil {
			log.Fatalf("error encrypting source: %s", err)
		}

		os.Stdout.Write(ciphertext)
	} else if dec {
		plaintext, err := decrypt([]byte(key), input)
		if err != nil {
			log.Fatalf("error encrypting source: %s", err)
		}

		os.Stdout.Write(plaintext)
	}
}
