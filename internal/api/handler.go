package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/kalunik/testCurrencyBalance/internal/entity"
	"github.com/kalunik/testCurrencyBalance/internal/usecase"
	"github.com/kalunik/testCurrencyBalance/pkg/log"
	"net/http"
	"strconv"
)

type Handlers interface {
	depositFunds(w http.ResponseWriter, r *http.Request)
	withdrawFunds(w http.ResponseWriter, r *http.Request)
	getBalance(w http.ResponseWriter, r *http.Request)
}

type walletHandlers struct {
	walletUsecase usecase.Usecase
	log           log.Logger
}

func (h walletHandlers) depositFunds(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	trans := entity.Transaction{}
	if err := json.NewDecoder(r.Body).Decode(&trans); err != nil {
		responseError(http.StatusInternalServerError, "deposit: json decode fail", err, w, h.log)
		return
	}

	transaction, err := h.walletUsecase.DepositFunds(accountID, trans)
	if err != nil {
		responseError(http.StatusBadRequest, "deposit: can't create transaction", err, w, h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		responseError(http.StatusInternalServerError, "deposit: json encode fail", err, w, h.log)
		return
	}
}

func (h walletHandlers) withdrawFunds(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	trans := entity.Transaction{}
	if err := json.NewDecoder(r.Body).Decode(&trans); err != nil {
		responseError(http.StatusInternalServerError, "withdraw: json decode fail", err, w, h.log)
		return
	}

	transaction, err := h.walletUsecase.WithdrawFunds(accountID, trans)
	if err != nil {
		responseError(http.StatusBadRequest, "deposit: can't create transaction", err, w, h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		responseError(http.StatusInternalServerError, "withdraw: json encode fail", err, w, h.log)
		return
	}
}

func (h walletHandlers) getBalance(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.Atoi(chi.URLParam(r, "id"))
	wallets, err := h.walletUsecase.GetBalance(accountID)
	if err != nil {
		responseError(http.StatusInternalServerError, "balance: can't perform", err, w, h.log)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(wallets); err != nil {
		return
	}

}

func NewWalletHandlers(walletUC usecase.Usecase, log log.Logger) Handlers {
	return &walletHandlers{
		walletUsecase: walletUC,
		log:           log,
	}
}
