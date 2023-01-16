package blockchain

import(
	"fmt"
	"github.com/dgraph-io/badger"
)

const(
	// dbPath = "./tmp/blocks"
)


type BlockChain struct{
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	/* This is to regain iteratiblity on our chain */
	CurrentHash	[]byte
	Database	*badger.DB
}

func InitBlockChain() *BlockChain {
	/* Creates a blockchain & 1st makes a genesis block*/

	var lastHash []byte

	// 1. Setting up options for DB
	// opts := badger.DefaultOptions
	// opts.Dir = dbPath		 // path to store keys
	// opts.ValueDir = dbPath	 // path to store values

	// 2. Opens database session
	db, err := badger.Open(badger.DefaultOptions("tmp/blocks"))
	Handle(err)
	
	// 3. Initiating Blockchain in DB or Getting lastHash of existing blockchain in DB
	err = db.Update(func(txn *badger.Txn) error {
		// Incase of no existing Blockchain (lastHash not found)
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound{
			fmt.Println("No existing BlockChain found.")

			// 1. Iniates a new blockchain with Genesis block
			genesis := Genesis()
			fmt.Println("Genesis proved !!")

			// 2. Storing genesis block in DB (key:Hash of Genesis block, val: serialized block)
			err := txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			// 3. Storing lastHash value in DB
			err = txn.Set([]byte("lh"), genesis.Hash) // this error will be handled outside closure

			// 4. Updating lasthash on memory
			lastHash = genesis.Hash 

			return err
		} else /* when DB and blockhain already exists */{
			/* We just note the lastHash*/

			item, err:= txn.Get([]byte("lh")) // Retrieving data from DB
			Handle(err)

			lastHash, err = item.ValueCopy(nil) // Extrating the value form Key-Value pair
			return err
		}
	})
	Handle(err)

	// 4. Creating the blockchain obj & returning the reference
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}	

func (chain *BlockChain) AddBlock(data string){
	/* This adds the block to Chain(DB/Ledger)*/
	var lastHash []byte

	// 1. Making a read in DB to find LastHash
	err := chain.Database.View(func(txn *badger.Txn) error {
		// read in db returns (key, value) pair
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		// Updating the lastHash var(on memory)
		lastHash, err = item.ValueCopy(nil)

		return err
	})

	// 2. Creating new block
	newBlock := CreateBlock(data, lastHash)

	// 3. Updating DB with new block and new lastHash
	err = chain.Database.Update(func(txn *badger.Txn) error {
		// i. Adding new block (hash, serialized block)
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)

		// ii. Updating LH key in DB
		err = txn.Set([]byte("lh"), newBlock.Hash)

		// iii. Updating lastHash on memory(of main chain)
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	/* This creates the iterator for chain */

	// Converting blockchain struct -> blockchainiterator struct
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	/* This helps to iterate through blockchain 
	   (works in reverse order: last block to genesis block) 
	*/

	var block *Block // empty block obj

	// 1. Querrying DB for Blocks
	err := iter.Database.View(func(txn *badger.Txn) error {
		// i. Returns (key, value)pair -> (Hash, serialized block)
		item, err := txn.Get(iter.CurrentHash) 
		Handle(err)

		// ii. Accessing the encoded block
		encodedBlock, err := item.ValueCopy(nil)

		// iii. Deserializing the encoded block
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	// 2. Updating iterator
	iter.CurrentHash = block.PrevHash

	return block
}