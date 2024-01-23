package usecase

import (
	"encoding/json"
	"github.com/kalunik/testCurrencyBalance/internal/entity"
	repo "github.com/kalunik/testCurrencyBalance/internal/repository"
)

type Usecase interface {
	DepositFunds(accountId int, transaction entity.Transaction) (entity.Transaction, error)
	WithdrawFunds(accountId int, transaction entity.Transaction) (entity.Transaction, error)
	GetBalance(accountId int) (entity.Account, error)

	ActivateTransaction()
}

type walletUc struct {
	pgRepo   repo.PgRepository
	natsRepo repo.NatsRepository
}

func (w walletUc) DepositFunds(accountId int, transaction entity.Transaction) (entity.Transaction, error) {
	transaction, err := w.pgRepo.DepositFunds(accountId, transaction)
	if err != nil {
		return transaction, err
	}
	tr := entity.RequestTransaction{
		UserId: accountId,
		Trans:  transaction,
	}
	bytes, err := json.Marshal(tr)
	if err != nil {
		return entity.Transaction{}, err
	}

	if err = w.natsRepo.Publish(bytes); err != nil {
		return entity.Transaction{}, err
	}
	return transaction, nil
}

func (w walletUc) WithdrawFunds(accountId int, transaction entity.Transaction) (entity.Transaction, error) {
	transaction.Withdraw = true
	transaction, err := w.pgRepo.WithdrawFunds(accountId, transaction)
	if err != nil {
		return transaction, err
	}

	tr := entity.RequestTransaction{
		UserId: accountId,
		Trans:  transaction,
	}
	bytes, err := json.Marshal(tr)
	if err != nil {
		return entity.Transaction{}, err
	}

	if err = w.natsRepo.Publish(bytes); err != nil {
		return entity.Transaction{}, err
	}

	return transaction, nil
}

func (w walletUc) GetBalance(accountId int) (entity.Account, error) {
	wallets, err := w.pgRepo.GetBalance(accountId)
	if err != nil {
		return entity.Account{}, err
	}
	return wallets, nil
}

func (w walletUc) ActivateTransaction() {
	dataChan := make(chan []byte)
	defer close(dataChan)
	sub, err := w.natsRepo.Subscribe(dataChan)
	if err != nil {
		return
	}
	defer sub.Unsubscribe()

	data := entity.RequestTransaction{}
	if err = json.Unmarshal(<-dataChan, &data); err != nil {
		return
	}

	if err = w.pgRepo.ActivateTransaction(data); err != nil {
		return
	}
}

func NewWalletUsecase(pgRepo repo.PgRepository, natsRepo repo.NatsRepository) Usecase {
	return &walletUc{pgRepo: pgRepo, natsRepo: natsRepo}
}
