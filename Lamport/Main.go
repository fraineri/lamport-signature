package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Block [32]byte
type Key struct {
	zero [256]Block
	one  [256]Block
}

func main() {
	rand.Seed(time.Now().Unix())

	privateKey, publicKey := generateKeys()
	err := saveKeys(privateKey, publicKey)
	if err != nil {
		return
	}
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
