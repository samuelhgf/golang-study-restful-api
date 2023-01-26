package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samuelhgf/golang-study-restful-api/src/controllers"
	"github.com/samuelhgf/golang-study-restful-api/src/middleware"
	"github.com/samuelhgf/golang-study-restful-api/src/services"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup, userService services.UserService) {
	router := rg.Group("/auth")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(userService), rc.authController.LogoutUser)
}
