package gin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"playground/app"
	"playground/pkg/token"
)

type entryHandler struct {
	u app.EntryUsecase
}

func NewEntryHandler(u app.EntryUsecase) entryHandler {
	return entryHandler{
		u: u,
	}
}

func (h entryHandler) List(c *gin.Context) {
	var args app.ListEntriesInputParams
	if err := c.ShouldBindQuery(&args); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var err error
	args.WalletID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	args.RequestUserID = authPayload.UserID

	fmt.Println("ListEntriesParams", args)
	if es, err := h.u.List(c, args); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		c.JSON(http.StatusOK, es)
	}
}
