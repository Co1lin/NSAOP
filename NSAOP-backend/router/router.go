package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"nsaop/config"
	_ "nsaop/docs"
	Location "nsaop/router/api/v2/location"
	Service "nsaop/router/api/v2/service"
	User "nsaop/router/api/v2/user"
	"nsaop/router/jwt"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	return r
}

func SetupRouter(router *gin.Engine) {
	r := router.Group("/" + config.Router.GetString("version"))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user := r.Group("/user")
	{
		user.POST("/check/username", User.CheckUsername)
		user.POST("/signup", User.Signup)
		user.POST("/login", User.Login)
		user.POST("/refresh", User.Refresh)
		user.POST("/logout", User.Logout)
		user.GET("/resetpasswd", User.ResetPasswdRequest)
		user.POST("/resetpasswd", User.ResetPassword)
		user.Use(jwt.Auth())
		{
			user.GET("/email", User.SendEmail)
			user.POST("/check/password", User.CheckPassword)
			user.Any("/detail", User.Detail)
		}
	}

	service := r.Group("/service").Use(jwt.Auth())
	{
		service.Any("", Service.Service)
		service.Use(Service.IdMiddleware()).Any("/:id", Service.ServiceId)
		service.Use(Service.IdMiddleware()).Any("/:id/device", Service.Devices)
		service.Use(Service.IdMiddleware()).Any("/:id/ssid", Service.SSID)
		service.Use(Service.IdMiddleware()).GET("/:id/traffic", Service.GetTraffic)
	}

	location := r.Group("/location").Use(jwt.Auth())
	{
		location.Any("", Location.Location)
		location.Any("/:id", Location.LocationId)
	}
}

func Run() {
	router := InitRouter()
	SetupRouter(router)
	router.Run(config.Router.GetString("port"))
}
