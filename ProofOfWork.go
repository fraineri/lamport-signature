package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	randSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randSource)
	var nonce uint32 = r.Uint32()
	var tries uint32 = 0

	var difficulty uint8 = 7
	var sb strings.Builder
	for i := 0; i < int(difficulty); i++ {
		sb.WriteString("0")
	}
	expected := sb.String()

	var solution [32]byte
	for {
		solution = sha256.Sum256([]byte(strconv.FormatUint(uint64(nonce), 10)))
		hexSolution := hex.EncodeToString(solution[:])
		if strings.HasPrefix(hexSolution, expected) {
			fmt.Println("~~SOLUTION~~")
			fmt.Printf("TIME: %v \n", time.Since(start).Seconds())
			fmt.Printf("DIFFICULTY: %v \n", difficulty)
			fmt.Printf("HASH: %x \n", solution)
			fmt.Printf("NONCE: %v \n", nonce)
			fmt.Printf("TRIES: %v \n", tries)
			break
		}
		nonce = r.Uint32()
		tries += 1
	}

}
