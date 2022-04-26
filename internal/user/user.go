package user

import (
	"errors"

	manager "github.com/adrianos93/wallet-manager"
	"github.com/adrianos93/wallet-manager/internal/wallet"
)

type User struct {
	Id      string                    `json:"Id"`
	Wallets map[string]*wallet.Wallet `json:"Wallets,omitempty"`
}

const userIdSize = 16

var Users = map[string]*User{}

func New() *User {
	user := &User{
		Id:      manager.GenerateId(userIdSize),
		Wallets: map[string]*wallet.Wallet{},
	}
	Users[user.Id] = user
	return user
}

func (u *User) CreateWallet() *wallet.Wallet {
	wallet := wallet.New()
	u.Wallets[wallet.Id] = wallet
	return wallet
}

func (u *User) Deposit(walletId string, amount float64) (wallet.Balance, error) {
	userWallet, found := u.Wallets[walletId]
	if !found {
		return wallet.Balance{}, errors.New("unauthorized transaction")
	}
	return userWallet.Deposit(amount), nil
}

func (u *User) Withdraw(walletId string, amount float64) (wallet.Balance, error) {
	userWallet, found := u.Wallets[walletId]
	if !found {
		return wallet.Balance{}, errors.New("unauthorized transaction")
	}
	return userWallet.Withdraw(amount)
}

func (u *User) CheckBalance(walletId string) (wallet.Balance, error) {
	userWallet, found := u.Wallets[walletId]
	if !found {
		return wallet.Balance{}, errors.New("unauthorized transaction")
	}
	return userWallet.CheckBalance(), nil
}

func (u *User) InitiatePayment(sourceWalletId, targetWalletId string, amount float64) (wallet.Payment, error) {
	intiatorWallet, found := u.Wallets[sourceWalletId]
	if !found {
		return wallet.Payment{}, errors.New("unauthorized transaction")
	}
	return intiatorWallet.InitiatePayment(targetWalletId, amount)
}
