package auth

import (
	"github.com/gin-gonic/gin"
	"backend/handlers"
)

func Logout(c *gin.Context) {
	// Clear the token cookie by setting maxAge to -1
	c.SetCookie(
		"token",  // name
		"",       // value (empty)
		-1,       // maxAge (negative value deletes the cookie)
		"/",      // path
		"",       // domain
		false,    // secure
		true,     // httpOnly
	)

	handlers.SuccessResponse(c, gin.H{
		"message": "Logged out successfully",
	})
}
