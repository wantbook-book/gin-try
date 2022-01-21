package main

import (
	"github.com/gin-gonic/gin"
	"hannibal/gin-try/controller"
	"hannibal/gin-try/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("api/auth/info", middleware.AuthMiddleware(), controller.Info)

	categoryRoutes := r.Group("categories")
	categoryController := controller.NewCategoryController()
	categoryRoutes.POST("", categoryController.Create)
	categoryRoutes.PUT("/:id", categoryController.Update)
	categoryRoutes.DELETE("/:id", categoryController.Delete)
	categoryRoutes.GET("/:id", categoryController.Show)
	return r
}
