package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/adrianos93/wallet-manager/internal/user"
	"github.com/adrianos93/wallet-manager/internal/wallet"
	"github.com/gorilla/mux"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	createdUser := user.New()
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdUser)
}

func HandleCreateWallet(w http.ResponseWriter, r *http.Request) {
	userRequested := mux.Vars(r)["user"]
	userData, found := user.Users[userRequested]
	if !found {
		http.Error(w, fmt.Errorf("user %s not found", userRequested).Error(), http.StatusNotFound)
		return
	}
	walletToReturn := userData.CreateWallet()
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(walletToReturn)
}

func HandleDeposit(w http.ResponseWriter, r *http.Request) {
	userRequested, walletRequested := mux.Vars(r)["user"], mux.Vars(r)["wallet"]
	userData, found := user.Users[userRequested]
	if !found {
		http.Error(w, fmt.Errorf("user %s not found", userRequested).Error(), http.StatusNotFound)
		return
	}
	_, found = wallet.Wallets[walletRequested]
	if !found {
		http.Error(w, fmt.Errorf("wallet %s not found", walletRequested).Error(), http.StatusNotFound)
		return
	}
	var input wallet.Deposit
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, errors.New("invalid json").Error(), http.StatusBadRequest)
		return
	}
	balanceToReturn, err := userData.Deposit(walletRequested, input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(balanceToReturn)
}

func HandleWithdrawal(w http.ResponseWriter, r *http.Request) {
	userRequested, walletRequested := mux.Vars(r)["user"], mux.Vars(r)["wallet"]
	userData, found := user.Users[userRequested]
	if !found {
		http.Error(w, fmt.Errorf("user %s not found", userRequested).Error(), http.StatusNotFound)
		return
	}
	_, found = wallet.Wallets[walletRequested]
	if !found {
		http.Error(w, fmt.Errorf("wallet %s not found", walletRequested).Error(), http.StatusNotFound)
		return
	}
	var input wallet.Withdraw
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, errors.New("invalid json").Error(), http.StatusBadRequest)
		return
	}
	balanceToReturn, err := userData.Withdraw(walletRequested, input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(balanceToReturn)
}

func HandleBalanceCheck(w http.ResponseWriter, r *http.Request) {
	userRequested, walletRequested := mux.Vars(r)["user"], mux.Vars(r)["wallet"]
	userData, found := user.Users[userRequested]
	if !found {
		http.Error(w, fmt.Errorf("user %s not found", userRequested).Error(), http.StatusNotFound)
		return
	}
	_, found = wallet.Wallets[walletRequested]
	if !found {
		http.Error(w, fmt.Errorf("wallet %s not found", walletRequested).Error(), http.StatusNotFound)
		return
	}
	balanceToReturn, err := userData.CheckBalance(walletRequested)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	_ = json.NewEncoder(w).Encode(balanceToReturn)
}

func HandlePayment(w http.ResponseWriter, r *http.Request) {
	userRequested, walletRequested := mux.Vars(r)["user"], mux.Vars(r)["wallet"]
	userData, found := user.Users[userRequested]
	if !found {
		http.Error(w, fmt.Errorf("user %s not found", userRequested).Error(), http.StatusNotFound)
		return
	}
	_, found = wallet.Wallets[walletRequested]
	if !found {
		http.Error(w, fmt.Errorf("wallet %s not found", walletRequested).Error(), http.StatusNotFound)
		return
	}
	var paymentRequest wallet.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
		http.Error(w, errors.New("invalid json").Error(), http.StatusBadRequest)
		return
	}
	payment, err := userData.InitiatePayment(walletRequested, paymentRequest.TargetWallet, paymentRequest.Amount)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "unauthorized"):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		case strings.Contains(err.Error(), "insufficient funds"):
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
	}
	_ = json.NewEncoder(w).Encode(payment)
}
