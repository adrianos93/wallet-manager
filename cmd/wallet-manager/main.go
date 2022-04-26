package main

import (
	"fmt"
	"net/http"

	manager "github.com/adrianos93/wallet-manager"
	"github.com/adrianos93/wallet-manager/internal/server"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc(fmt.Sprintf("/v1/health/%s", manager.ServiceName), func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)
	r.HandleFunc("/v1/user", server.HandleCreateUser).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/{user:[A-Za-z0-9]{1,64}}/wallet", server.HandleCreateWallet).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/{user:[A-Za-z0-9]{1,64}}/wallet/{wallet:[A-Za-z0-9]{1,64}}/balance", server.HandleBalanceCheck).Methods(http.MethodGet)
	r.HandleFunc("/v1/user/{user:[A-Za-z0-9]{1,64}}/wallet/{wallet:[A-Za-z0-9]{1,64}}/deposit", server.HandleDeposit).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/{user:[A-Za-z0-9]{1,64}}/wallet/{wallet:[A-Za-z0-9]{1,64}}/withdraw", server.HandleWithdrawal).Methods(http.MethodPost)
	r.HandleFunc("/v1/user/{user:[A-Za-z0-9]{1,64}}/wallet/{wallet:[A-Za-z0-9]{1,64}}/payment", server.HandlePayment).Methods(http.MethodPost)

	fmt.Println("Listening on port 8080")
	_ = http.ListenAndServe(":8080", r)
}
