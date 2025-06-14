package middilware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func RouTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UnixMilli()

		fmt.Println("Runtime运行前")

		c.Next()

		end := time.Now().UnixMilli()

		fmt.Sprintf("代码运行时间%d", end-start)
	}
}
