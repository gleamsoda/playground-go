package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morikuni/failure"

	"playground/internal/pkg/apperr"
	"playground/internal/pkg/token"
	"playground/internal/wallet"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (h handler) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if a, err := h.u.CreateAccount(c, &wallet.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, a)
	}
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (h handler) getAccount(c *gin.Context) {
	var req getAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if a, err := h.u.GetAccount(c, &wallet.GetAccountsParams{
		ID:    req.ID,
		Owner: authPayload.Username,
	}); err != nil {
		if code, ok := failure.CodeOf(err); ok {
			switch code {
			case apperr.NotFound:
				c.JSON(http.StatusNotFound, errorResponse(err))
			case apperr.Unauthorized:
				c.JSON(http.StatusUnauthorized, errorResponse(err))
			default:
				c.JSON(http.StatusInternalServerError, errorResponse(err))
			}
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, a)
	}
}

type listAccountRequest struct {
	Limit  int32 `form:"limit" binding:"required,min=1,max=100"`
	Offset int32 `form:"offset"`
}

func (h handler) listAccounts(c *gin.Context) {
	var req listAccountRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if as, err := h.u.ListAccounts(c, &wallet.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.Limit,
		Offset: req.Offset,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, as)
	}
}
