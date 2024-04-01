package repository

import (
	// "fmt"
	// "golang_socmed/internal/model"
	// "log"
	// "strconv"

	// "github.com/gofiber/fiber/v2"

	"log"

	"github.com/jmoiron/sqlx"
	// "github.com/lib/pq"
	"github.com/paimon_bank/internal/entity"
	"github.com/paimon_bank/internal/model"
)

type TransactionRepository struct {
	DB *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (b *TransactionRepository) Create(request *entity.Transaction) error {
	query := `INSERT INTO transactions (id, userid, bankaccountnumber, bankname, balance, currency) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := b.DB.Exec(query, request.ID, request.UserId, request.BankAccountNumber, request.BankName, request.Balance, request.Currency)
	return err
}

func (b *TransactionRepository) CreateTransaction(request *entity.Transaction) (err error) {
	tx := b.DB.MustBegin()

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = rollbackErr
			}
		}
	}()

	query := `INSERT INTO transactions (id, userid, bankaccountnumber, bankname, balance, currency, flow) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.Exec(query, request.ID, request.UserId, request.BankAccountNumber, request.BankName, request.Balance, request.Currency, "transaction")
	if err != nil {
		return err
	}

	queryBalance := `UPDATE balances SET balance = balance - $1 WHERE userid = $2 and currency = $3`

	_, err = tx.Exec(queryBalance, request.Balance, request.UserId, request.Currency)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (b *TransactionRepository) GetTransactions(filter *model.TransactionFilter, userId string) (*[]entity.Transaction, error) {
	var transaction []entity.Transaction
	query := `SELECT * FROM transactions WHERE userid = $1 LIMIT $2 OFFSET $3`

	err := b.DB.Select(&transaction, query, userId, filter.Limit, filter.Offset)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}

	return &transaction, nil
}
