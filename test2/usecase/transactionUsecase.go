package usecase

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"test2/model"
	"test2/model/entity"
	"test2/repository/cacheRepository"
	"test2/repository/sqlRepository"
)

type Transactions interface {
	DoTransactions(tu model.TransactionRequest, calcBalance func(b, nb float64) float64, additionalFunction func(tx *gorm.DB, tr model.TransactionRequest) *model.Error) (model.TransactionResponse, *model.Error)
	CalcBalanceTopUp(b, nb float64) float64
	CalcBalancePaymentAndTransfer(b, nb float64) float64
	TransferToAnotherUser(tx *gorm.DB, tr model.TransactionRequest) *model.Error
	GetList(tr model.TransactionListRequest) ([]model.TransactionListResponse, *model.Error)
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

	var targetUserID uuid.UUID
	if tr.Type == "TRANSFER" {
		targetUserID, errs = uuid.Parse(tr.UserID)
		if errs != nil {
			return model.TransactionResponse{}, model.NewError(400, "invalid target user", errs)
		}
	}

	// insert data transaction with
	transaction := entity.Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		TargetUserID:    targetUserID,
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

	// update user if transfer
	if tr.Type == "TRANSFER" {
		if err = additionalFunction(tx, tr); err != nil {
			return model.TransactionResponse{}, err
		}
	}

	var topUpID, paymentID, transferID string
	if tr.Type == "TOPUP" {
		topUpID = transaction.ID.String()
	} else if tr.Type == "PAYMENT" {
		paymentID = transaction.ID.String()
	} else {
		transferID = transaction.ID.String()
	}
	return model.TransactionResponse{
		TopUpID:       topUpID,
		PaymentID:     paymentID,
		TransferID:    transferID,
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
	if user.ID == uuid.Nil {
		return 0, model.NewError(400, "unknown target user", nil)
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
	if tr.TargetUser == tr.UserID {
		return model.NewError(400, "target user is my self, not allowed to transfer", nil)
	}

	balance, err := t.getBalanceUserData(tr.TargetUser)
	if err != nil {
		return err
	}
	newBalance := balance + tr.Amount
	return t.userRepository.UpdateBalance(tx, tr.TargetUser, newBalance)
}

func (t Transaction) GetList(tr model.TransactionListRequest) ([]model.TransactionListResponse, *model.Error) {
	if tr.Page < 1 || tr.Limit < 1 {
		return nil, model.NewError(400, "invalid parameter page or limit", nil)
	}

	list, err := t.transactionRepository.GetList(tr.UserID, sqlRepository.Clause{
		Limit:  tr.Limit,
		Offset: (tr.Page - 1) * tr.Limit,
		Order:  "updated_at asc",
	})
	if err != nil {
		return nil, err
	}

	var results []model.TransactionListResponse
	for _, l := range list {
		var topup, pay, transfer string
		if l.TransactionType == "TOPUP" {
			topup = l.ID.String()
		} else if l.TransactionType == "PAYMENT" {
			pay = l.ID.String()
		} else {
			transfer = l.ID.String()
		}

		results = append(results, model.TransactionListResponse{
			TopUpID:       topup,
			PaymentID:     pay,
			TransferID:    transfer,
			UserID:        l.UserID.String(),
			Remarks:       l.Remarks,
			BalanceBefore: l.BalanceBefore,
			BalanceAfter:  l.BalanceAfter,
			CreatedDate:   l.CreatedAt.String(),
		})
	}
	return results, nil
}
