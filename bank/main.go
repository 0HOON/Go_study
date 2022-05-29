package main

import (
	"fmt"

	accounts "example.com/bank/banking"
)

func main() {
	account := accounts.NewAccount("YH")
	account.Deposit(10)
	err := account.Withdraw(100)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(account)
}
