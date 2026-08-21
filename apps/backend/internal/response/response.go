package response

import "github.com/gin-gonic/gin"

func OKWithData[T any](c *gin.Context, code int, data T) {
	c.JSON(
		code,
		OKResponse[T]{
			Status: "ok",
			Data:   data,
		},
	)
}

func OK(c *gin.Context, code int) {
	c.JSON(
		code,
		OKResponse[any]{
			Status: "ok",
		},
	)
}

func FailedWithReason(c *gin.Context, code int, error string, reason string) {
	c.AbortWithStatusJSON(code,
		FailedResponse{
			Status: "failed",
			Error:  error,
			Reason: reason,
		},
	)
}

func Failed(c *gin.Context, code int, error string) {
	c.AbortWithStatusJSON(code,
		FailedResponse{
			Status: "failed",
			Error:  error,
		},
	)
}
