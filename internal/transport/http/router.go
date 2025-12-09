package http

import (
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(userSvc *service.UserService) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
			"service": "user",
		})
	})

	fbAuth := NewFirebaseAuth()
	v1 := r.Group("/api/v1")
	{
		usrs := v1.Group("/users")
		usrs.Use(fbAuth.Middleware())
		// usrs.Use(firebaseAuth.Middleware())
		{
			userHandler := NewUserHandler(userSvc)
			usrs.POST("", userHandler.Create)
			usrs.GET("", userHandler.List)
			usrs.GET("/:dni", userHandler.GetByDni)
			usrs.PUT("/:user", userHandler.Update)
			// usrs.DELETE("/:dni", userHandler.Delete)
		}
	}
	return r
}
