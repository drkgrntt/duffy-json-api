package controllers

import (
	"net/http"
	"time"

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

	pastWeek := time.Now().AddDate(0, 0, -7)

	c.DB.Where("has_tkts_data = ?", true).
		Where("last_scanned_at > ?", pastWeek).
		Joins("CompetitionGroup").
		Preload("Shows", "showtime > ? ORDER BY showtime DESC", pastWeek).
		Preload("Shows.Listings", "scanned_at > ? ORDER BY scanned_at DESC", pastWeek).
		Preload("CompetitionGroup.Productions").
		Order("last_scanned_at DESC").
		Order("last_shown_at DESC").
		Find(&productions)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"productions": productions}})
}
