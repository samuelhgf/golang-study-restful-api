package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samuelhgf/golang-study-restful-api/src/config"
	"github.com/samuelhgf/golang-study-restful-api/src/services"
	"github.com/samuelhgf/golang-study-restful-api/src/utils"
	"net/http"
	"strings"
)

func DeserializeUser(userService services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		configEnvs, _ := config.LoadConfig(".")
		sub, err := utils.ValidateToken(accessToken, configEnvs.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		user, err := userService.FindUserById(fmt.Sprint(sub))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "The user belonging to this token no longer exists"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
