package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type account struct {
	owner   string
	balance int
}

var errNoMoney = errors.New("can't withdraw")

// NewAccount creates Account
func NewAccount(Owner string) *account {
	account := account{owner: Owner, balance: 0}
	return &account
}

// Deposit x amount on your account
func (a *account) Deposit(amount int) {
	a.balance += amount
}

// Balance of your account
func (a account) Balance() int {
	return a.balance
}

// Withdraw x amount from your account
func (a *account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

// ChangeOwner of the account
func (a *account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner of the account
func (a account) Owner() string {
	return a.owner
}

func (a account) String() string {
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Balance())
}
