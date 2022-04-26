package user

import (
	"testing"

	"github.com/adrianos93/wallet-manager/internal/wallet"
	"github.com/stretchr/testify/require"
)

func TestUser_New(t *testing.T) {
	for name, test := range map[string]struct {
		wantUsers int
	}{
		"returns 1 new User": {
			wantUsers: 1,
		},
		"returns 2 new Users": {
			wantUsers: 2,
		},
	} {
		t.Run(name, func(t *testing.T) {
			loops := 0
			for loops < test.wantUsers {
				New()
				loops++
			}
			require.Equal(t, test.wantUsers, len(Users))
			Users = map[string]*User{}
		})
	}
}

func TestUser_CreateWallet(t *testing.T) {
	for name, test := range map[string]struct {
		wallets int
	}{
		"creates a wallet": {
			wallets: 1,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user := &User{
				Wallets: map[string]*wallet.Wallet{},
			}
			got := user.CreateWallet()
			require.Equal(t, test.wallets, len(wallet.Wallets))
			require.Equal(t, wallet.Wallets[got.Id], got)
		})
	}
}

func TestUser_Deposit(t *testing.T) {
	for name, test := range map[string]struct {
		walletId string
		amount   float64

		wantResult wallet.Balance
		wantErr    bool
	}{
		"process a deposit": {
			walletId:   "somerandomID",
			amount:     100,
			wantResult: wallet.Balance{},
		},
		"fail to process a deposit": {
			walletId:   "somerandomID",
			amount:     100,
			wantResult: wallet.Balance{},
			wantErr:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user := &User{
				Wallets: map[string]*wallet.Wallet{
					test.walletId: {
						Id: test.walletId,
					},
				},
			}
			if test.wantErr {
				delete(user.Wallets, test.walletId)
			}

			got, err := user.Deposit(test.walletId, test.amount)
			if test.wantErr {
				require.Error(t, err)
			}
			require.IsType(t, test.wantResult, got)
		})
	}
}

func TestUser_Withdraw(t *testing.T) {
	for name, test := range map[string]struct {
		walletId string
		amount   float64

		wantResult wallet.Balance
		wantErr    bool
	}{
		"process a withdrawal": {
			walletId:   "somerandomID",
			amount:     100,
			wantResult: wallet.Balance{},
		},
		"fail to process a withdrawal": {
			walletId:   "somerandomID",
			amount:     100,
			wantResult: wallet.Balance{},
			wantErr:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user := &User{
				Wallets: map[string]*wallet.Wallet{
					test.walletId: {
						Id:      test.walletId,
						Balance: 100,
					},
				},
			}
			if test.wantErr {
				delete(user.Wallets, test.walletId)
			}

			got, err := user.Withdraw(test.walletId, test.amount)
			if test.wantErr {
				require.Error(t, err)
			}
			require.IsType(t, test.wantResult, got)
		})
	}
}

func TestUser_CheckBalance(t *testing.T) {
	for name, test := range map[string]struct {
		walletId string

		wantResult wallet.Balance
		wantErr    bool
	}{
		"get balance": {
			walletId:   "somerandomID",
			wantResult: wallet.Balance{},
		},
		"fail to get balance": {
			walletId:   "somerandomID",
			wantResult: wallet.Balance{},
			wantErr:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user := &User{
				Wallets: map[string]*wallet.Wallet{
					test.walletId: {
						Id:      test.walletId,
						Balance: 100,
					},
				},
			}
			if test.wantErr {
				delete(user.Wallets, test.walletId)
			}

			got, err := user.CheckBalance(test.walletId)
			if test.wantErr {
				require.Error(t, err)
			}
			require.IsType(t, test.wantResult, got)
		})
	}
}

func TestUser_InitiatePayment(t *testing.T) {
	for name, test := range map[string]struct {
		sourceWalletId, targetWalletId string
		amount                         float64

		wantResult wallet.Payment
		wantErr    bool
	}{
		"successfully initiate payment": {
			sourceWalletId: "somerandomID",
			targetWalletId: "someotherID",
			amount:         100,
			wantResult:     wallet.Payment{},
		},
		"fail to initiate payment": {
			sourceWalletId: "somerandomID",
			targetWalletId: "someotherID",
			amount:         100,
			wantResult:     wallet.Payment{},
			wantErr:        true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user := &User{
				Wallets: map[string]*wallet.Wallet{
					test.sourceWalletId: {
						Id:      test.sourceWalletId,
						Balance: 100,
					},
				},
			}
			if test.wantErr {
				delete(user.Wallets, test.sourceWalletId)
			}

			got, err := user.InitiatePayment(test.sourceWalletId, test.targetWalletId, test.amount)
			if test.wantErr {
				require.Error(t, err)
			}
			require.IsType(t, test.wantResult, got)
		})
	}
}
