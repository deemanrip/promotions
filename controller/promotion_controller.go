package controller

import (
	"github.com/deemanrip/promotions/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GinInit() {
	router := gin.Default()
	router.Use(ErrorHandler)
	router.GET("/promotions/:id", GetPromotion)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func GetPromotion(context *gin.Context) {
	id := context.Param("id")
	promotion, err := service.GetPromotionById(&id)

	if err != nil {
		log.Error(err)
		_ = context.AbortWithError(http.StatusInternalServerError, err)
	} else if promotion != nil {
		context.JSON(http.StatusOK, *promotion)
	} else {
		context.JSON(http.StatusNotFound, gin.H{"message": "promotion not found"})
	}
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
}
