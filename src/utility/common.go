package utility

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuchangeu/kong-assignment/src/constant"
)

func GenerateUUIDForContext(ctx *gin.Context) *gin.Context {
	uuid.SetClockSequence(-1)
	uuid.SetNodeID([]byte{0xa1})

	//current time should be determined
	uuID, _ := uuid.NewUUID()
	_, exists := ctx.Get(constant.KeyContextId)
	if !exists {
		ctx.Set(constant.KeyContextId, uuID)
	}

	return ctx
}
