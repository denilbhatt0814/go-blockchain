package main

import(
	"fmt"
	"time"
	"os"
	"flag"
	"runtime"
	"strconv"
	"github.com/denilbhatt0814/go-blockchain/blockchain"
)

type CommandLine struct{
	blockchain 	*blockchain.BlockChain
} 

func (cli *CommandLine) printUsage() {
	/* Lets user know to Use the cli */
	fmt.Println("Usage: ")
	fmt.Println(" add -block BLOCK_DATA - add block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	/* Validates arguments of cli */
	
	// in case args are < 2
	if len(os.Args) < 2 {
		cli.printUsage() // Let user know how to use CLI
		runtime.Goexit() // terminates Goroutine unlike os exit
		// It initaites a proper shutdown letting DB to close properly
	}
}

func (cli *CommandLine) addBlock(data string) {
	/* Adds block through cli to main chain */
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) printChain() {
	/* Prints whole bllockchain using iterator */
	
	// 1. Accesing the iterator 
	iter := cli.blockchain.Iterator()

	// 2. Printing all the blocks
	for {
		block := iter.Next()

		fmt.Printf("\n")
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)


		// Checking the validity of block
		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		// 3. Exit when reached Genesis block
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	// Validating recieved args
	cli.validateArgs()

	// Setting Flags for cli
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)

	// to add other options in addBlockCmd to addBlockData
	addBlockData := addBlockCmd.String("block", "", "Block data")

	// Setting condn. and parsing
	switch os.Args[1] {
	// In case user type add
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	// In case user type print
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	// In any other case
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	// If parsed succesfully
	if addBlockCmd.Parsed() {
		// if parsing is empty string
		if *addBlockData == ""{
			cli.printUsage()
			runtime.Goexit()
		}
		//else
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main(){
	defer os.Exit(0) // making sure program 
	// To see how much time whole program takes to run
	t := time.Now()
	defer fmt.Println(time.Since(t))

	// Initializing block chain
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close() // making sure to close DB after use

	cli := CommandLine{chain}
	cli.run()
}