package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	consumerIdCtx = "id"
)

func getId(ctx *gin.Context) (int, error) {
	id := ctx.Query(consumerIdCtx)

	intId, ok := strconv.Atoi(id)
	if ok != nil {
		return 0, errors.New("id is of invalid type")
	}

	return intId, nil
}
