package controllers

import (
	"net/http"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShowsController struct {
	DB *gorm.DB
}

func NewShowsController(DB *gorm.DB) ShowsController {
	return ShowsController{DB}
}

func (c *ShowsController) GetProductions(ctx *gin.Context) {
	var productions []models.Production
	c.DB.Find(&productions).Where("has_tkts_data = true").Order("last_scanned_at DESC").Order("last_shown_at DESC")

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"productions": productions}})
}
