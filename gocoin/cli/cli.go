package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"example.com/gocoin/network"
	"example.com/gocoin/wallet"

	"example.com/gocoin/blockchain"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for the address")
	fmt.Println(" createblockchain -address ADDRESS - creates a blockchain")
	fmt.Println(" printchain - prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO amount AMOUNT -mine - send amount from FROM to TO. Then -mine flag is set, mine off of this node")
	fmt.Println(" createwallet - Creates a now Wallet")
	fmt.Println(" listaddresses - Lists the addresses in our wallet file")
	fmt.Println(" reindexutxo - Rebuilds the UTXO set")
	fmt.Println(" startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) reindexUTXO(nodeID string) {
	chain := blockchain.ContinueBlockChain(nodeID)
	defer chain.Database.Close()
	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

func (cli *CommandLine) listAddresses(nodeID string) {
	wallets, _ := wallet.CreateWallets(nodeID)
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CommandLine) StartNode(nodeID, minerAddress string) {
	fmt.Printf("Starting Node %s\n", nodeID)

	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address!")
		}
	}
	network.StartServer(nodeID, minerAddress)
}

func (cli *CommandLine) createWallet(nodeID string) {
	wallets, _ := wallet.CreateWallets(nodeID)
	address := wallets.AddWallet()
	wallets.SaveFile(nodeID)

	fmt.Printf("New address: %s\n", address)
}

func (cli *CommandLine) printChain(nodeID string) {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockChain(address, nodeID string) {

	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}

	chain := blockchain.InitBlockChain(address, nodeID)
	defer chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{chain}
	UTXOSet.Reindex()

	fmt.Println("Finished")
}

func (cli *CommandLine) getBalance(address, nodeID string) {

	if !wallet.ValidateAddress(address) {
		log.Panic("Address is not Valid")
	}

	chain := blockchain.ContinueBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	balance := 0
	pubKeyHash := wallet.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := UTXOSet.FindUnspentTransactions(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balanceof %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int, nodeID string, mineNow bool) {

	if !wallet.ValidateAddress(to) {
		log.Panic("Address is not Valid")
	}

	if !wallet.ValidateAddress(from) {
		log.Panic("Address is not Valid")
	}

	chain := blockchain.ContinueBlockChain(nodeID)
	UTXOSet := blockchain.UTXOSet{chain}
	defer chain.Database.Close()

	wallets, err := wallet.CreateWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}

	wallet := wallets.GetWallet(from)

	tx := blockchain.NewTransaction(&wallet, to, amount, &UTXOSet)

	if mineNow {
		cbTx := blockchain.CoinbaseTx(from, "")
		txs := []*blockchain.Transaction{cbTx, tx}
		block := chain.MineBlock(txs)
		UTXOSet.Update(block)
	} else {
		fmt.Println("aaa")
		network.SendTx(network.KnownNodes[0], tx)
		fmt.Println("send tx")
	}

	fmt.Println("Success")
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID env is not set:")
		runtime.Goexit()
	}

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)
	reindexUTXOCmd := flag.NewFlagSet("reindexutxo", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance from")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to create")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")
	sendMine := sendCmd.Bool("mine", false, "Mine immediately on the same node")
	startNodeMiner := startNodeCmd.String("miner", "", "Enable mining node and send reward")

	switch os.Args[1] {
	case "reindexutxo":
		err := reindexUTXOCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if reindexUTXOCmd.Parsed() {
		cli.reindexUTXO(nodeID)
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			cli.printUsage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress, nodeID)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			cli.printUsage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress, nodeID)
	}

	if printChainCmd.Parsed() {
		cli.printChain(nodeID)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			cli.printUsage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount, nodeID, *sendMine)
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses(nodeID)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet(nodeID)
	}

	if startNodeCmd.Parsed() {
		cli.StartNode(nodeID, *startNodeMiner)
	}
}
