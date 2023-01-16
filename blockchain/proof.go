package blockchain

import(
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
	"math"
	"fmt"
	"crypto/sha256"
)

// Take the data from the block

// create a counter (nonce) which starts at 0

// create hash of data + counter

// check the resulting hash if it meets a set of requirements

/* Requirements:
	1. first few bytes must contain 0s
*/

const Difficulty = 18

type ProofOfWork struct{
	Block *Block
	Target *big.Int
}


func NewProof(b *Block) *ProofOfWork {
	// This steps manages the requirement of 0s
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - Difficulty)) // Lsh = left shift i.e 256 - difficulty set by us 

	// create the proof of work
	pow := &ProofOfWork{b, target}
	// and return the address of created Proof of work
	return pow
}

func (pow *ProofOfWork) InitData(nonce int) []byte {

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	// running for nearl infinity
	// This runs for verifying untill matches the req hash (nonce keep changing)
	for nonce < math.MaxInt64 {
		// Creating new sets each time
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:]) // converting to bytes to check

		// Verifying
		if intHash.Cmp(pow.Target) == -1 {
			// breaks loop when hashes match
			break
		} else {
			// else continues with new nonce
			nonce++
		}
	}

	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	/* Here we try to recreate blocks hash and match with the created block's */
	var intHash big.Int

	// recreating the parameters
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:]) // conv. to byte format

	// Comparing and returning validity(bool)
	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	/* Used to convert nonce(int) to byte format */
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}