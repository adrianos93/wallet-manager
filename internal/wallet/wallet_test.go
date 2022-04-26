package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWallet_New(t *testing.T) {
	for name, test := range map[string]struct {
		wantWallets int
	}{
		"returns 1 new Wallet": {
			wantWallets: 1,
		},
		"returns 2 new Wallets": {
			wantWallets: 2,
		},
	} {
		t.Run(name, func(t *testing.T) {
			loops := 0
			for loops < test.wantWallets {
				New()
				loops++
			}
			require.Equal(t, test.wantWallets, len(Wallets))
			Wallets = map[string]*Wallet{}
		})
	}
}

func TestWallet_Deposit(t *testing.T) {
	for name, test := range map[string]struct {
		amount      float64
		wantBalance Balance
	}{
		"adds money successfully": {
			amount:      100,
			wantBalance: Balance{100},
		},
	} {
		t.Run(name, func(t *testing.T) {
			wallet := &Wallet{
				Balance: 0,
			}

			got := wallet.Deposit(test.amount)
			require.Equal(t, test.wantBalance, got)
		})
	}
}

func TestWallet_Withdraw(t *testing.T) {
	for name, test := range map[string]struct {
		amount      float64
		wantBalance Balance
		wantErr     bool
	}{
		"withdraws money successfully": {
			amount:      100,
			wantBalance: Balance{0},
		},
		"fails to withdraw money": {
			amount:  100,
			wantErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			wallet := &Wallet{
				Balance: 0,
			}
			wallet.Balance = 100
			if test.wantErr {
				wallet.Balance = 0
			}

			got, err := wallet.Withdraw(test.amount)
			if test.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, test.wantBalance, got)
		})
	}
}

func TestWallet_CheckBalance(t *testing.T) {
	for name, test := range map[string]struct {
		wantBalance Balance
	}{
		"returns balance": {
			wantBalance: Balance{100},
		},
	} {
		t.Run(name, func(t *testing.T) {
			wallet := &Wallet{
				Balance: 100,
			}
			got := wallet.CheckBalance()
			require.Equal(t, test.wantBalance, got)
		})
	}
}

func TestWallet_InitiatePayment(t *testing.T) {
	for name, test := range map[string]struct {
		sourceWalletId, targetWalletId string
		initialAmount, amountToPay     float64

		wantPayment            Payment
		wantErr, wantTargetErr bool
	}{
		"successfully initiates a payment": {
			sourceWalletId: "sourceId",
			targetWalletId: "targetId",
			initialAmount:  100,
			amountToPay:    50,
			wantPayment: Payment{
				Balance: 50,
			},
		},
		"fails to pay because due to inssuficient funds": {
			sourceWalletId: "sourceId",
			targetWalletId: "targetId",
			initialAmount:  40,
			amountToPay:    50,
			wantErr:        true,
		},
		"fails to pay because target wallet does not exist": {
			sourceWalletId: "sourceId",
			targetWalletId: "targetId",
			initialAmount:  100,
			amountToPay:    50,
			wantErr:        true,
			wantTargetErr:  true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			Wallets[test.targetWalletId] = &Wallet{
				Id:      test.targetWalletId,
				Balance: 0,
			}

			if test.wantTargetErr {
				delete(Wallets, test.targetWalletId)
			}

			sourceWallet := &Wallet{
				Id:      test.sourceWalletId,
				Balance: test.initialAmount,
			}

			got, err := sourceWallet.InitiatePayment(test.targetWalletId, test.amountToPay)
			if test.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, test.wantPayment.Balance, got.Balance)
		})
	}
}
