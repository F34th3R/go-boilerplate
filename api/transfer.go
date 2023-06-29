package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/F34th3R/go_simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type transferRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required,uuid"`
	ToAccountID   string `json:"to_account_id" binding:"required,uuid"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	FromAccountID := uuid.MustParse(req.FromAccountID)
	ToAccountID := uuid.MustParse(req.ToAccountID)

	if !server.validAccount(ctx, FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: FromAccountID,
		ToAccountID:   ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID uuid.UUID, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%s] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
