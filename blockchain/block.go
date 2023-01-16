package blockchain

import(
	"bytes"
	"log"
	"encoding/gob"
)

type Block struct{
	Hash		[]byte
	Data		[]byte
	PrevHash	[]byte
	Nonce 		int
}

func CreateBlock(data string, prevHash []byte) *Block{
	/* Create a block object & stores address to return */
	block := &Block{[]byte{}, []byte(data), prevHash, 0}

	// Running the proof of work algorithm to derive hash
	pow := NewProof(block)
	nonce, hash := pow.Run()

	// Making updates to created block 
	block.Hash = hash[:]
	block.Nonce = nonce

	/* JUNK CODE:
	 // Hash is derived and updated as it is a pointer method 
	 block.DeriveHash() 
	*/
	return block
}

func Genesis() *Block {
	/* Creates a genesis block */
	return CreateBlock("Genesis", []byte{})
}

func (b *Block) Serialize() []byte {
	/* Serialize the block(byte bormat) to 
		store in DB(ledger) */
	
	var res bytes.Buffer
	// 1. Creating new encoder
	encoder := gob.NewEncoder(&res)

	// 2. Encoding the block
	err := encoder.Encode(b)
	Handle(err)

	// 3. Returns the byte portion of resulted encoding of block
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	/* Deserializing data from DB(ledger) to Block */
	
	// 1. Creating an empty block
	var block Block

	// 2. Creating decoder
	decoder := gob.NewDecoder(bytes.NewReader(data))

	// 3. Decoding and storing in block var
	err := decoder.Decode(&block)
	Handle(err)

	// 4. returning address of deserialized block
	return &block
}

func Handle(err error) {
	/* Used to handle errors */
	if err != nil{
		log.Panic(err)
	}
}

/* JUNK CODE:
func (b *Block) DeriveHash() {
	// Uses Data and prevhash to create new hash for block 
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}
*/