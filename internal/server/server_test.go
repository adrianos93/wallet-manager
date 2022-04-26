package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adrianos93/wallet-manager/internal/user"
	"github.com/adrianos93/wallet-manager/internal/wallet"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestServer_HandleCreateUser(t *testing.T) {
	for name, test := range map[string]struct {
		wantCode int
	}{
		"golden path": {
			wantCode: 201,
		},
	} {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/v1/user", nil)
			HandleCreateUser(w, r)
			require.Equal(t, test.wantCode, w.Code)
		})
	}
}

func TestServer_HandleCreateWallet(t *testing.T) {
	for name, test := range map[string]struct {
		wantCode int
		wantErr  bool
	}{
		"golden path": {
			wantCode: 201,
		},
		"user not found": {
			wantCode: 404,
			wantErr:  true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user.Users["user1"] = &user.User{
				Id:      "user1",
				Wallets: map[string]*wallet.Wallet{},
			}
			vars := map[string]string{
				"user": "user1",
			}
			if test.wantErr {
				delete(user.Users, "user1")
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/v1/user/user1/wallet", nil)
			r = mux.SetURLVars(r, vars)
			HandleCreateWallet(w, r)
			require.Equal(t, test.wantCode, w.Code)
		})
	}
}

func TestServer_HandleDeposit(t *testing.T) {
	input := wallet.Deposit{Amount: 100}
	for name, test := range map[string]struct {
		wantCode                                   int
		body                                       []byte
		wantUserErr, wantWalletErr, wantDepositErr bool
	}{
		"golden path": {
			wantCode: 200,
			body:     func() (b []byte) { b, _ = json.Marshal(input); return }(),
		},
		"user not found": {
			wantCode:    404,
			wantUserErr: true,
		},
		"wallet not found": {
			wantCode:      404,
			wantWalletErr: true,
		},
		"bad request": {
			wantCode: 400,
			body:     []byte(`i'm not json`),
		},
		"not your wallet": {
			wantCode:       401,
			body:           func() (b []byte) { b, _ = json.Marshal(input); return }(),
			wantDepositErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user.Users["user1"] = &user.User{
				Id: "user1",
				Wallets: map[string]*wallet.Wallet{
					"wallet1": {Id: "wallet1", Balance: 0},
				},
			}
			wallet.Wallets["wallet1"] = &wallet.Wallet{
				Id:      "wallet1",
				Balance: 0,
			}
			vars := map[string]string{
				"user":   "user1",
				"wallet": "wallet1",
			}
			if test.wantUserErr {
				delete(user.Users, "user1")
			}
			if test.wantWalletErr {
				delete(wallet.Wallets, "wallet1")
			}
			if test.wantDepositErr {
				walletToDelete := user.Users["user1"]
				delete(walletToDelete.Wallets, "wallet1")
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/v1/user/user1/wallet/wallet1/deposit", strings.NewReader(string(test.body)))
			r = mux.SetURLVars(r, vars)
			HandleDeposit(w, r)
			require.Equal(t, test.wantCode, w.Code)
		})
	}
}

func TestServer_HandleWithdrawal(t *testing.T) {
	input := wallet.Withdraw{Amount: 100}
	for name, test := range map[string]struct {
		wantCode                                    int
		body                                        []byte
		wantUserErr, wantWalletErr, wantWithdrawErr bool
	}{
		"golden path": {
			wantCode: 200,
			body:     func() (b []byte) { b, _ = json.Marshal(input); return }(),
		},
		"user not found": {
			wantCode:    404,
			wantUserErr: true,
		},
		"wallet not found": {
			wantCode:      404,
			wantWalletErr: true,
		},
		"bad request": {
			wantCode: 400,
			body:     []byte(`i'm not json`),
		},
		"not your wallet": {
			wantCode:        401,
			body:            func() (b []byte) { b, _ = json.Marshal(input); return }(),
			wantWithdrawErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user.Users["user1"] = &user.User{
				Id: "user1",
				Wallets: map[string]*wallet.Wallet{
					"wallet1": {Id: "wallet1", Balance: 100},
				},
			}
			wallet.Wallets["wallet1"] = &wallet.Wallet{
				Id:      "wallet1",
				Balance: 0,
			}
			vars := map[string]string{
				"user":   "user1",
				"wallet": "wallet1",
			}
			if test.wantUserErr {
				delete(user.Users, "user1")
			}
			if test.wantWalletErr {
				delete(wallet.Wallets, "wallet1")
			}
			if test.wantWithdrawErr {
				walletToDelete := user.Users["user1"]
				delete(walletToDelete.Wallets, "wallet1")
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/v1/user/user1/wallet/wallet1/withdraw", strings.NewReader(string(test.body)))
			r = mux.SetURLVars(r, vars)
			HandleWithdrawal(w, r)
			require.Equal(t, test.wantCode, w.Code)
		})
	}
}

func TestServer_HandleBalanceCheck(t *testing.T) {
	for name, test := range map[string]struct {
		wantCode                                   int
		wantUserErr, wantWalletErr, wantBalanceErr bool
	}{
		"golden path": {
			wantCode: 200,
		},
		"user not found": {
			wantCode:    404,
			wantUserErr: true,
		},
		"wallet not found": {
			wantCode:      404,
			wantWalletErr: true,
		},
		"not your wallet": {
			wantCode:       401,
			wantBalanceErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			user.Users["user1"] = &user.User{
				Id: "user1",
				Wallets: map[string]*wallet.Wallet{
					"wallet1": {Id: "wallet1", Balance: 0},
				},
			}
			wallet.Wallets["wallet1"] = &wallet.Wallet{
				Id:      "wallet1",
				Balance: 0,
			}
			vars := map[string]string{
				"user":   "user1",
				"wallet": "wallet1",
			}
			if test.wantUserErr {
				delete(user.Users, "user1")
			}
			if test.wantWalletErr {
				delete(wallet.Wallets, "wallet1")
			}
			if test.wantBalanceErr {
				walletToDelete := user.Users["user1"]
				delete(walletToDelete.Wallets, "wallet1")
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/v1/user/user1/wallet/wallet1/balance", nil)
			r = mux.SetURLVars(r, vars)
			HandleBalanceCheck(w, r)
			require.Equal(t, test.wantCode, w.Code)
		})
	}
}

func TestServer_HandlePayment(t *testing.T) {
	input := wallet.PaymentRequest{
		TargetWallet: "wallet2",
		Amount:       50,
	}
	inputInsufficient := wallet.PaymentRequest{
		TargetWallet: "wallet2",
		Amount:       150,
	}
	for name, test := range map[string]struct {
		wantCode                                   int
		body                                       []byte
		wantUserErr, wantWalletErr, wantPaymentErr bool
	}{
		"golden path": {
			wantCode: 200,
			body:     func() (b []byte) { b, _ = json.Marshal(input); return }(),
		},
		"user not found": {
			wantCode:    404,
			wantUserErr: true,
		},
		"wallet not found": {
			wantCode:      404,
			wantWalletErr: true,
		},
		"bad request": {
			wantCode: 400,
			body:     []byte(`i'm not json`),
		},
		"not your wallet": {
			wantCode:       401,
			body:           func() (b []byte) { b, _ = json.Marshal(input); return }(),
			wantPaymentErr: true,
		},
		"inssufficient funds": {
			wantCode: 403,
			body:     func() (b []byte) { b, _ = json.Marshal(inputInsufficient); return }(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			user.Users["user1"] = &user.User{
				Id: "user1",
				Wallets: map[string]*wallet.Wallet{
					"wallet1": {Id: "wallet1", Balance: 100},
				},
			}
			wallet.Wallets["wallet1"] = &wallet.Wallet{
				Id:      "wallet1",
				Balance: 100,
			}
			wallet.Wallets["wallet2"] = &wallet.Wallet{
				Id:      "wallet2",
				Balance: 0,
			}
			vars := map[string]string{
				"user":   "user1",
				"wallet": "wallet1",
			}
			if test.wantUserErr {
				delete(user.Users, "user1")
			}
			if test.wantWalletErr {
				delete(wallet.Wallets, "wallet1")
			}
			if test.wantPaymentErr {
				walletToDelete := user.Users["user1"]
				delete(walletToDelete.Wallets, "wallet1")
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/v1/user/user1/wallet/wallet1/payment", strings.NewReader(string(test.body)))
			r = mux.SetURLVars(r, vars)
			HandlePayment(w, r)
			require.Equal(t, test.wantCode, w.Code)
		})
	}
}
