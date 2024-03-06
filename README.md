# Go-Blockchain

**Description:**
This project is a simple command-line interface (CLI) application for managing a basic blockchain. It allows users to add blocks to the blockchain and view the entire chain. Badger DB is used as the underlying database for storing the blockchain data efficiently.

**How to Use:**

1. **Installation:**
   - Clone the repository from [Github](https://github.com/denilbhatt0814/go-blockchain/).
   - Ensure you have Go 1.19 or above installed on your machine.

2. **Running the Application:**
   - Navigate to the project directory in your terminal.
   - Run the command `go run main.go` to start the application.

3. **Commands:**
   - `addblock <data>`: Adds a new block to the blockchain with the provided data.
   - `print`: Prints the entire blockchain.

4. **Badger DB:**
   - Badger DB is a fast, embeddable, and persistent key-value (KV) database written in pure Go. It is used in this project to store and manage the blockchain data efficiently.

5. **Example Usage:**
   - To add a block: `addblock "Block Data Here"`
   - To print the blockchain: `print`

6. **Important Notes:**
   - Make sure to provide valid data when adding a block.
   - The blockchain will automatically close the database connection upon exiting the application.

7. **Additional Information:**
   - This project is a basic demonstration of blockchain concepts using Go and Badger DB.
   - Feel free to explore and modify the code as needed.

**Thank you for using the Blockchain CLI!**
