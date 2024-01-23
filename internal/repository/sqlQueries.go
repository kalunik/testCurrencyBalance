package repository

const (
	createWallet = `INSERT INTO wallet (user_id, currency, active_balance, frozen_balance)
						SELECT $1, $2, 0.0, 0.0
						WHERE NOT EXISTS (SELECT user_id,currency FROM wallet WHERE user_id = $1 AND currency = $2);`

	checkFrozenBalance = `SELECT id, frozen_balance FROM wallet WHERE user_id = $1 AND currency = $2;`

	checkActiveBalance = `SELECT id, active_balance FROM wallet WHERE user_id = $1 AND currency = $2;`

	checkBalance = `SELECT id, active_balance, frozen_balance FROM wallet WHERE user_id = $1 AND currency = $2;`

	updateFrozenBalance = `UPDATE wallet SET frozen_balance = $2 WHERE id = $1;`

	updateActiveFrozen = `UPDATE wallet SET active_balance = $2, frozen_balance = $3 WHERE id = $1;`

	showBalance = `SELECT active_balance, frozen_balance, currency FROM wallet WHERE user_id = $1;`

	addTransaction = `INSERT INTO transactions (wallet_id, status, amount, withdraw, card_number, created_at) 
						VALUES ($1, 'CREATED', $2, $3, $4, (SELECT LOCALTIMESTAMP))
						RETURNING id;`

	updateStatusTransaction = `UPDATE transactions SET status = $2 WHERE id = $1`
)
