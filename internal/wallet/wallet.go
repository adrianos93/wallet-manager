package wallet

import (
	"errors"
	"fmt"
	"sync"
	"time"

	manager "github.com/adrianos93/wallet-manager"
)

type Wallet struct {
	Id           string  `json:"Id"`
	Balance      float64 `json:"Balance"`
	Transactions map[string]*Transaction
	sync.Mutex
}

type Transaction struct {
	Id             string
	AmountChanged  float64
	Balance        float64
	Timestamp      time.Time
	SourceWalletID string
	Reference      string
}

type Payment struct {
	TransactionId string  `json:"TransactionId"`
	Balance       float64 `json:"Balance"`
}

type Balance struct {
	Balance float64 `json:"Balance"`
}

type Deposit struct {
	Amount float64 `json:"Amount"`
}

type Withdraw struct {
	Amount float64 `json:"Amount"`
}

type PaymentRequest struct {
	TargetWallet string  `json:"Creditor"`
	Amount       float64 `json:"Amount"`
}

const (
	walletIdSize      = 16
	transactionIdSize = 32
)

var Wallets = map[string]*Wallet{}

func New() *Wallet {
	wallet := &Wallet{
		Id:      manager.GenerateId(walletIdSize),
		Balance: 0,
	}
	Wallets[wallet.Id] = wallet
	return wallet
}

func (w *Wallet) Deposit(amount float64) Balance {
	w.Balance += amount
	return Balance{
		w.Balance,
	}
}

func (w *Wallet) Withdraw(amount float64) (Balance, error) {
	if amount > w.Balance {
		return Balance{}, errors.New("inssuficient funds in wallet")
	}
	w.Balance = w.Balance - amount
	return Balance{
		w.Balance,
	}, nil
}

func (w *Wallet) CheckBalance() Balance {
	return Balance{w.Balance}
}

func (w *Wallet) InitiatePayment(walletId string, amount float64) (Payment, error) {
	w.Lock()
	defer w.Unlock()
	targetWallet, found := Wallets[walletId]
	if !found {
		return Payment{}, fmt.Errorf("wallet with ID: %s does not exist", walletId)
	}
	if amount > w.Balance {
		return Payment{}, errors.New("insufficient funds")
	}
	w.Balance = w.Balance - amount
	targetWallet.Balance = targetWallet.Balance + amount
	transactionId := manager.GenerateId(transactionIdSize)
	w.Transactions[transactionId] = &Transaction{
		Id:             transactionId,
		AmountChanged:  amount,
		Balance:        w.Balance,
		Timestamp:      time.Now(),
		SourceWalletID: targetWallet.Id,
		Reference:      "",
	}
	Wallets[w.Id] = w
	Wallets[targetWallet.Id] = targetWallet

	return Payment{
		TransactionId: transactionId,
		Balance:       w.Balance,
	}, nil
}
