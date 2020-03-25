package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

var (
	enc bool
	dec bool
	gen bool
	key string
	pub string

	errDecryptionError = errors.New("error: decryption failed")
)

func init() {
	flag.BoolVar(&enc, "e", false, "encrypt contents of message")
	flag.BoolVar(&dec, "d", false, "decrypt contents of message")
	flag.BoolVar(&gen, "g", false, "generate public/private key pair")
	flag.StringVar(&key, "k", "", "secret key | private key")
	flag.StringVar(&pub, "p", "", "recipient/sender key | public key")

	flag.Parse()
}

func genkeys() (*[32]byte, *[32]byte, error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return pub, priv, nil
}

func encrypt(key, pub, message []byte) ([]byte, error) {
	var (
		secretKey [32]byte
		publicKey [32]byte
		nonce     [24]byte
	)

	copy(secretKey[:], key)

	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	if len(pub) == 32 {
		copy(publicKey[:], pub)

		return box.Seal(nonce[:], message, &nonce, &publicKey, &secretKey), nil
	} else if len(pub) > 0 {
		return nil, errors.New("Recipient public key must be 32 bytes long")
	}

	return secretbox.Seal(nonce[:], message, &nonce, &secretKey), nil
}

func decrypt(key, sender, ciphertext []byte) ([]byte, error) {
	var (
		decryptNonce [24]byte
		secretKey    [32]byte
		publicKey    [32]byte
	)

	copy(secretKey[:], key)

	copy(decryptNonce[:], ciphertext[:24])

	if len(sender) == 32 {
		copy(publicKey[:], sender)

		plaintext, ok := box.Open(nil, ciphertext[24:], &decryptNonce, &publicKey, &secretKey)
		if !ok {
			return nil, errDecryptionError
		}
		return plaintext, nil
	} else if len(sender) > 0 {
		return nil, errors.New("Sender public key must be 32 bytes long")
	}

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

	if !enc && !dec && !gen {
		fmt.Printf("Usage: %s [options] [message]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	if gen {
		pub, priv, err := genkeys()
		if err != nil {
			log.Fatalf("error generating public/private key pair: %s", err)
		}
		fmt.Printf("Private Key: %X\n", *priv)
		fmt.Printf("Public Key: %X\n", *pub)
		os.Exit(0)
	}

	if len(flag.Args()) == 1 {
		input = []byte(flag.Arg(0))
	} else {
		input, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading message from stdin: %s", err)
		}
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		log.Fatalf("error decoding secret/private key: %s", err)
	}

	pubBytes, err := hex.DecodeString(pub)
	if err != nil {
		log.Fatalf("error decoding secret/private key: %s", err)
	}

	if enc {
		ciphertext, err := encrypt(keyBytes, pubBytes, input)
		if err != nil {
			log.Fatalf("error encrypting source: %s", err)
		}

		os.Stdout.Write(ciphertext)
	} else if dec {
		plaintext, err := decrypt(keyBytes, pubBytes, input)
		if err != nil {
			log.Fatalf("error encrypting source: %s", err)
		}

		os.Stdout.Write(plaintext)
	}
}
