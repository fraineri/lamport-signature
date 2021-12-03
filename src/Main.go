package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
)

type Block [32]byte
type Key struct {
	zero [256]Block
	one  [256]Block
}

func main() {
	// rand.Seed(time.Now().Unix())

	// KEY GENERATION
	privateKey, publicKey := generateKeys()
	err := saveKeys(privateKey, publicKey)
	if err != nil {
		return
	}

	// MESSAGE SIGNING
	var message string = "Hello world!"
	var messageHash Block = sha256.Sum256([]byte(message))
	fmt.Printf("Message Hash: %x \n", messageHash)

	signature := generateSignature(messageHash, privateKey)
	err = saveSignature(signature)
	if err != nil {
		return
	}

	// SIGNATURE VERIFICATION
	// <-- PUT BREAKPOINT TO MODIGY MESSAGE FOR TESTING
	var signatureIdx int = 0
	for _, hex := range sha256.Sum256([]byte(message)) {
		var bitCompare int = 128
		for i := 0; i < 8; i++ {
			bitWise := int(hex) & bitCompare
			if int(bitWise) == bitCompare {
				if signature[signatureIdx] != privateKey.one[signatureIdx] {
					fmt.Println("INVALID MESSAGE")
					return
				}
			} else {
				if signature[signatureIdx] != privateKey.zero[signatureIdx] {
					fmt.Println("INVALID MESSAGE")
					return
				}
			}
			bitCompare /= 2
			signatureIdx += 1
		}
	}
	fmt.Println("VALID MESSAGE")
}

func generateKeys() (Key, Key) {
	var privateKey Key
	for i := 0; i < 256; i++ {
		rand.Read(privateKey.zero[i][:])
		rand.Read(privateKey.one[i][:])
	}

	var publicKey Key
	for i, num := range privateKey.zero {
		publicKey.zero[i] = sha256.Sum256(num[:])
	}

	for i, num := range privateKey.one {
		publicKey.one[i] = sha256.Sum256(num[:])
	}

	return privateKey, publicKey
}

func saveKeys(privateKey Key, publicKey Key) error {
	// CREATE FILE AND SAVE PRIVATE KEY
	file, err := createFile("private.keys")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = saveKey(file, privateKey)
	if err != nil {
		fmt.Println(err)
		file.Close()
		return err
	}
	file.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// CREATE FILE AND SAVE PUBLIC KEY
	file, err = createFile("public.keys")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = saveKey(file, publicKey)
	if err != nil {
		fmt.Println(err)
		file.Close()
		return err
	}
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func createFile(fileName string) (os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return *file, err
	}

	return *file, err
}

func saveKey(file os.File, key Key) error {
	var line string
	for i := 0; i < 256; i++ {
		line = hex.EncodeToString(key.zero[i][:]) + "," + hex.EncodeToString(key.one[i][:]) + "\n"
		_, err := file.WriteString(line)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func generateSignature(messageHash Block, privateKey Key) [256]Block {
	var signature [256]Block
	var signatureIdx int = 0
	for _, hex := range messageHash {
		var bitCompare int = 128
		for i := 0; i < 8; i++ {
			bitWise := int(hex) & bitCompare
			if int(bitWise) == bitCompare {
				signature[signatureIdx] = privateKey.one[signatureIdx]
			} else {
				signature[signatureIdx] = privateKey.zero[signatureIdx]
			}
			bitCompare /= 2
			signatureIdx += 1
		}
	}

	return signature
}

func saveSignature(signature [256]Block) error {
	file, err := createFile("signature.keys")
	if err != nil {
		fmt.Println(err)
		return err
	}

	var line string
	for i := 0; i < 256; i++ {
		line = hex.EncodeToString(signature[i][:]) + "\n"
		_, err := file.WriteString(line)

		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
