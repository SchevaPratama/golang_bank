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

type BalanceRepository struct {
	DB *sqlx.DB
}

func NewBalanceRepository(db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{DB: db}
}

func (b *BalanceRepository) Create(request *entity.Balance, balanceRequest *model.BalanceRequest) (err error) {
	tx := b.DB.MustBegin()
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = rollbackErr
			}
		}
	}()

	queryBalance := `INSERT INTO balances VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(queryBalance, request.ID, request.UserId, request.Balance, request.Currency)
	if err != nil {
		return err
	}

	queryTransaction := `INSERT INTO transactions (id, userid, bankaccountnumber, bankname, transferproofimage, balance, currency, flow) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(queryTransaction, request.ID, request.UserId, balanceRequest.SenderBankAccountNumber, balanceRequest.SenderBankName, balanceRequest.TransferProofImage, request.Balance, request.Currency, "topup")
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (b *BalanceRepository) GetBalance(userId string) (*[]entity.Balance, error) {
	var balance []entity.Balance
	query := `SELECT * FROM balances WHERE userid = $1`

	err := b.DB.Select(&balance, query, userId)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}

	return &balance, nil
}

func (b *BalanceRepository) UpdateBalance(request *entity.Balance, balanceRequest *model.BalanceRequest) (err error) {

	tx := b.DB.MustBegin()
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = rollbackErr
			}
		}
	}()

	queryBalance := `UPDATE balances SET balance = balance + $1 WHERE userid = $2 and currency = $3`

	_, err = tx.Exec(queryBalance, request.Balance, request.UserId, request.Currency)
	if err != nil {
		return err
	}

	queryTransaction := `INSERT INTO transactions (id, userid, bankaccountnumber, bankname, transferproofimage, balance, currency, flow) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(queryTransaction, request.ID, request.UserId, balanceRequest.SenderBankAccountNumber, balanceRequest.SenderBankName, balanceRequest.TransferProofImage, request.Balance, request.Currency, "topup")
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
