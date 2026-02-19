package pkg

import (
	"time"

	"github.com/gin-gonic/gin"
)

func SetTokenToCookie(c *gin.Context, name, token, origin string) {
	c.SetCookie(name, token, int(time.Duration(15*time.Minute)), "/", "localhost", false, false)
}
