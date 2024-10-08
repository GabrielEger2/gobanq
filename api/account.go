package api

import (
	"database/sql"
	"net/http"

	db "github.com/GabrielEger2/gobanq/db/sqlc"
	"github.com/GabrielEger2/gobanq/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// createAccount creates a new account
// @Summary Create a new account
// @Tags accounts
// @Description create a new account with the input payload
// @Accept json
// @Produce json
// @Param input body createAccountRequest true "Account info"
// @Success 200 {object} db.Account
// @Security ApiKeyAuth
// @Router /accounts [post]
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet("authorization_payload").(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusConflict, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount returns an account by its ID
// @Summary Get an account
// @Tags accounts
// @Description get account by ID
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} db.Account
// @Security ApiKeyAuth
// @Router /accounts/{id} [get]
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccountById(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet("authorization_payload").(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errorResponse(err)
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=25"`
}

// listAccounts returns a list of accounts
// @Summary List accounts
// @Tags accounts
// @Description list accounts
// @Accept json
// @Produce json
// @Param page_id query int true "Page ID"
// @Param page_size query int true "Page Size"
// @Success 200 {object} []db.Account
// @Security ApiKeyAuth
// @Router /accounts [get]
func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet("authorization_payload").(*token.Payload)
	arg := db.GetAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
