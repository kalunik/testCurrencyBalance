package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kalunik/testCurrencyBalance/internal/entity"
)

type PgRepository interface {
	DepositFunds(accountId int, transaction entity.Transaction) (entity.Transaction, error)
	WithdrawFunds(accountId int, transaction entity.Transaction) (entity.Transaction, error)
	GetBalance(accountId int) (entity.Account, error)
	ActivateTransaction(data entity.RequestTransaction) error
}

type PgRepo struct {
	pg *pgxpool.Pool
}

func (p PgRepo) DepositFunds(userId int, transaction entity.Transaction) (entity.Transaction, error) {
	ctx := context.Background()
	t, _ := p.pg.Begin(ctx)
	defer t.Rollback(ctx)

	_, err := p.pg.Exec(ctx, createWallet, userId, transaction.Currency)
	if err != nil {
		return entity.Transaction{}, err
	}

	wallet := entity.Wallet{}
	err = p.pg.QueryRow(ctx, checkFrozenBalance, userId, transaction.Currency).Scan(&wallet.Id, &wallet.FrozenBalance)
	if err != nil {
		return entity.Transaction{}, err
	}

	err = p.pg.QueryRow(ctx, addTransaction, wallet.Id, transaction.Amount, false, transaction.CardNumber).Scan(&transaction.Id)
	if err != nil {
		return entity.Transaction{}, err
	}
	transaction.Status = "CREATED"

	if transaction.Amount <= 0 {
		transaction.Status = "ERROR"
		_, err = p.pg.Exec(ctx, updateStatusTransaction, transaction.Id, transaction.Status)
		if err != nil {
			return entity.Transaction{}, err
		}
		_ = t.Commit(ctx)
		return transaction, errors.New("amount should be more than 0")
	}

	_, err = p.pg.Exec(ctx, updateFrozenBalance, wallet.Id, wallet.FrozenBalance+transaction.Amount)
	if err != nil {
		return entity.Transaction{}, err
	}

	_ = t.Commit(ctx)
	return transaction, nil
}

func (p PgRepo) WithdrawFunds(accountId int, transaction entity.Transaction) (entity.Transaction, error) {
	ctx := context.Background()
	t, err := p.pg.Begin(ctx)
	if err != nil {
		return entity.Transaction{}, err
	}
	defer t.Rollback(ctx)

	w := entity.Wallet{}
	err = t.QueryRow(ctx, checkActiveBalance, accountId, transaction.Currency).Scan(&w.Id, &w.ActiveBalance)
	if err != nil {
		return entity.Transaction{}, err
	}

	err = p.pg.QueryRow(ctx, addTransaction, w.Id, transaction.Amount, true, transaction.CardNumber).Scan(&transaction.Id)
	if err != nil {
		return entity.Transaction{}, err
	}
	transaction.Status = "CREATED"

	if w.ActiveBalance < transaction.Amount {
		transaction.Status = "ERROR"
		_, err = p.pg.Exec(ctx, updateStatusTransaction, transaction.Id, transaction.Status)
		if err != nil {
			return entity.Transaction{}, err
		}
		_ = t.Commit(ctx)
		return transaction, errors.New("amount should be more than 0")
	}

	_, err = p.pg.Exec(ctx, updateActiveFrozen, w.Id,
		w.ActiveBalance-transaction.Amount,
		w.FrozenBalance+transaction.Amount,
	)
	if err != nil {
		return entity.Transaction{}, err
	}

	_ = t.Commit(ctx)
	return transaction, err
}

func (p PgRepo) GetBalance(accountId int) (entity.Account, error) {
	ctx := context.Background()
	t, err := p.pg.Begin(ctx)
	if err != nil {
		return entity.Account{}, err
	}
	defer t.Rollback(ctx)

	acc := entity.Account{accountId, nil}
	rows, err := t.Query(ctx, showBalance, accountId)
	if err != nil {
		return acc, err
	}

	wallets := []entity.Wallet{}
	for rows.Next() {
		w := entity.Wallet{}
		err := rows.Scan(&w.ActiveBalance, &w.FrozenBalance, &w.Currency)
		if err != nil {
			return acc, err
		}
		wallets = append(wallets, w)
	}
	acc.Wallet = wallets
	_ = t.Commit(ctx)
	return acc, nil
}

func (p PgRepo) ActivateTransaction(data entity.RequestTransaction) error {
	if data.Trans.Status == "SUCCESS" {
		return nil
	}

	ctx := context.Background()
	t, err := p.pg.Begin(ctx)
	if err != nil {
		return err
	}
	defer t.Rollback(ctx)

	w := entity.Wallet{}
	if err = t.QueryRow(ctx, checkBalance, data.UserId, data.Trans.Currency).Scan(&w.Id, &w.ActiveBalance, &w.FrozenBalance); err != nil {
		return err
	}

	if data.Trans.Withdraw {
		if _, err = t.Exec(ctx, updateFrozenBalance, w.Id,
			w.FrozenBalance-data.Trans.Amount); err != nil {
			return err
		}
	} else {
		if _, err = t.Exec(ctx, updateActiveFrozen,
			w.Id, w.ActiveBalance+data.Trans.Amount,
			w.FrozenBalance-data.Trans.Amount); err != nil {
			return err
		}
	}

	if _, err = t.Exec(ctx, updateStatusTransaction, data.Trans.Id, "SUCCESS"); err != nil {
		return err
	}

	_ = t.Commit(ctx)
	return nil
}

func NewPgRepository(pgpool *pgxpool.Pool) PgRepository {
	return &PgRepo{pg: pgpool}
}
