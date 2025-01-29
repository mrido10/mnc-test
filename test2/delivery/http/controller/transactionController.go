package controller

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"test2/model"
	"test2/usecase"
)

type TransactionController struct {
	transaction usecase.Transactions
}

func NewTransactionController(trans usecase.Transactions) TransactionController {
	return TransactionController{transaction: trans}
}

func (a TransactionController) TopUp(c *fiber.Ctx) (interface{}, *model.Error) {
	var req model.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return model.Response{}, model.NewError(400, "Body format error", nil)
	}

	req.UserAccess, _ = c.Locals("auth").(model.UserAccess)
	req.Type = "TOPUP"
	return a.transaction.DoTransactions(req, a.transaction.CalcBalanceTopUp, nil)
}

func (a TransactionController) Payment(c *fiber.Ctx) (interface{}, *model.Error) {
	var req model.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return model.Response{}, model.NewError(400, "Body format error", nil)
	}

	req.UserAccess, _ = c.Locals("auth").(model.UserAccess)
	req.Type = "PAYMENT"
	return a.transaction.DoTransactions(req, a.transaction.CalcBalancePaymentAndTransfer, nil)
}

func (a TransactionController) Transfer(c *fiber.Ctx) (interface{}, *model.Error) {
	var req model.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return model.Response{}, model.NewError(400, "Body format error", nil)
	}

	req.UserAccess, _ = c.Locals("auth").(model.UserAccess)
	req.Type = "TRANSFER"
	return a.transaction.DoTransactions(
		req,
		a.transaction.CalcBalancePaymentAndTransfer,
		a.transaction.TransferToAnotherUser)
}

func (a TransactionController) List(c *fiber.Ctx) (interface{}, *model.Error) {
	var req model.TransactionListRequest
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)

	req.Page = int(page)
	req.Limit = int(limit)
	req.UserAccess, _ = c.Locals("auth").(model.UserAccess)

	return a.transaction.GetList(req)
}
