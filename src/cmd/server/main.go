package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samuelhgf/golang-study-restful-api/src/config"
	"github.com/samuelhgf/golang-study-restful-api/src/controllers"
	"github.com/samuelhgf/golang-study-restful-api/src/routes"
	"github.com/samuelhgf/golang-study-restful-api/src/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	postService         services.PostService
	PostController      controllers.PostController
	postCollection      *mongo.Collection
	PostRouteController routes.PostRouteController
)

func init() {
	configEnvs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(configEnvs.DBUri)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		panic(err)
	}

	if err = mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	authCollection = mongoclient.Database("golang_mongodb").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	postCollection = mongoclient.Database("golang_mongodb").Collection("posts")
	postService = services.NewPostService(postCollection, ctx)
	PostController = controllers.NewPostController(postService)
	PostRouteController = routes.NewPostControllerRoute(PostController)

	server = gin.Default()
}

func main() {
	configEnvs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load configEnvs", err)
	}

	defer mongoclient.Disconnect(ctx)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	startGinServer(configEnvs)
}

func startGinServer(config config.Config) {
	router := server.Group("/api")
	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "All systems running"})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	PostRouteController.PostRoute(router)
	log.Fatal(server.Run(":" + config.Port))
}
