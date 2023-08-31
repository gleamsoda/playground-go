package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"playground/app"
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

	// TODO: Auth
	if a, err := h.u.CreateAccount(c, &app.CreateAccountParams{
		Owner:    "example",
		Balance:  1000000,
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

	if a, err := h.u.GetAccount(c, req.ID); err != nil {
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

	// TODO: Get Owner from AuthPayload
	if as, err := h.u.ListAccounts(c, &app.ListAccountsParams{
		Owner:  "example",
		Limit:  req.Limit,
		Offset: req.Offset,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, as)
	}
}
