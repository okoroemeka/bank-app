package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	customerror "github.com/okoroemeka/simple_bank/custom-error"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/token"
	"net/http"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if errCode := customerror.ErrorCode(err); errCode == customerror.UniqueViolation || errCode == customerror.ForeignKeyViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if errors.Is(err, customerror.ErrorNoRecordFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != account.Owner {
		err = errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

type listAccountRequest struct {
	Offset int32 `form:"offset" binding:"required,min=0"`
	Limit  int32 `form:"limit" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Owner:  authPayload.Username,
		Offset: int64(req.Offset),
		Limit:  int64(req.Limit),
	})

	if err != nil {
		if errors.Is(err, customerror.ErrorNoRecordFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
	return
}

// TODO: add delete and update account
type accountIDQueryRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type updateAccountBodyRequest struct {
	Amount int64 `form:"amount" binding:"required,min=1"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var reqQuery accountIDQueryRequest
	var reqBody updateAccountBodyRequest

	if err := ctx.ShouldBindUri(&reqQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.UpdateAccount(ctx, db.UpdateAccountParams{
		Balance: reqBody.Amount,
		ID:      reqQuery.ID,
	})

	if err != nil {
		if errors.Is(err, customerror.ErrorNoRecordFound) {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, account)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var reqQuery accountIDQueryRequest

	if err := ctx.ShouldBindUri(&reqQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := server.store.DeleteAccount(ctx, reqQuery.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
