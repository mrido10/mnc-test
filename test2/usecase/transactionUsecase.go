package usecase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"test2/model"
	"test2/model/entity"
	"test2/repository/cacheRepository"
	"test2/repository/sqlRepository"
	"test2/util"
)

type Transactions interface {
	DoTransactions(tu model.TransactionRequest, calcBalance func(b, nb float64) float64, additionalFunction func(tx *gorm.DB, tr model.TransactionRequest) *model.Error) (model.TransactionResponse, *model.Error)
	CalcBalanceTopUp(b, nb float64) float64
	CalcBalancePaymentAndTransfer(b, nb float64) float64
	TransferToAnotherUser(tx *gorm.DB, tr model.TransactionRequest) *model.Error
}
type Transaction struct {
	userRepository        sqlRepository.Users
	transactionRepository sqlRepository.Transactions
	sqlRepo               sqlRepository.SqlDB
	cache                 cacheRepository.Cache
}

func NewTransaction(
	userRepository sqlRepository.Users,
	transactionRepository sqlRepository.Transactions,
	sqlRepo sqlRepository.SqlDB,
	cache cacheRepository.Cache,
) Transactions {
	return Transaction{
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		sqlRepo:               sqlRepo,
		cache:                 cache,
	}
}

func (t Transaction) DoTransactions(
	tr model.TransactionRequest,
	calcBalance func(b, nb float64) float64,
	additionalFunction func(tx *gorm.DB, tr model.TransactionRequest) *model.Error,
) (model.TransactionResponse, *model.Error) {

	// validate request
	if err := t.validateRequest(tr); err != nil {
		return model.TransactionResponse{}, err
	}

	// get data user before
	balance, err := t.getBalanceUserData(tr.UserID)
	if err != nil {
		return model.TransactionResponse{}, err
	}
	newBalance := calcBalance(balance, tr.Amount)
	if newBalance < 0 {
		return model.TransactionResponse{}, model.NewError(400, "balance is not enough", nil)
	}

	tx := t.sqlRepo.DBBeginTX()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	userID, errs := uuid.Parse(tr.UserID)
	if errs != nil {
		return model.TransactionResponse{}, model.NewError(500, "Internal server error", errs)
	}
	// insert data transaction with type TOPUP
	transaction := entity.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		TransactionType: tr.Type,
		Amount:          tr.Amount,
		Remarks:         tr.Remarks,
		BalanceBefore:   balance,
		BalanceAfter:    newBalance,
	}
	if err = t.transactionRepository.Insert(tx, &transaction); err != nil {
		return model.TransactionResponse{}, err
	}

	// update user balance
	if err = t.userRepository.UpdateBalance(tx, tr.UserID, newBalance); err != nil {
		return model.TransactionResponse{}, err
	}

	var topUpID, paymentID string
	if tr.Type == "TOPUP" {
		topUpID = transaction.ID.String()
	} else {
		paymentID = transaction.ID.String()
	}
	return model.TransactionResponse{
		TopUpID:       topUpID,
		PaymentID:     paymentID,
		Amount:        transaction.Amount,
		BalanceBefore: balance,
		BalanceAfter:  newBalance,
		CreatedDate:   transaction.CreatedAt.String(),
	}, nil
}

func (Transaction) validateRequest(tr model.TransactionRequest) *model.Error {
	if tr.Amount < 10000 {
		return model.NewError(400, "minimum amount is 10000", nil)
	}
	return nil
}

func (t Transaction) getBalanceUserData(userID string) (float64, *model.Error) {
	user, err := t.userRepository.GetByID(userID)
	if err != nil {
		return 0, err
	}
	if util.IsEmptyString(user.ID.String()) {
		return 0, model.NewError(400, "unknown user", nil)
	}
	return user.Balance, nil
}

func (Transaction) CalcBalanceTopUp(b, nb float64) float64 {
	return b + nb
}

func (Transaction) CalcBalancePaymentAndTransfer(b, nb float64) float64 {
	return b - nb
}

func (t Transaction) TransferToAnotherUser(tx *gorm.DB, tr model.TransactionRequest) *model.Error {
	balance, err := t.getBalanceUserData(tr.TargetUser)
	if err != nil {
		return err
	}
	newBalance := balance + tr.Amount
	return t.userRepository.UpdateBalance(tx, tr.UserID, newBalance)
}
