package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	customerror "github.com/okoroemeka/simple_bank/custom-error"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/token"
	"net/http"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, isFromAccountCurrencyValid := server.validateCurrency(ctx, req.FromAccountID, req.Currency)

	if fromAccount.Owner != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("you are not authorised to perform this transaction")))
		return
	}
	_, istoAccountCurrencyValid := server.validateCurrency(ctx, req.ToAccountID, req.Currency)

	if !isFromAccountCurrencyValid || !istoAccountCurrencyValid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		if errors.Is(err, customerror.ErrorNoRecordFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
