package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShowController struct {
	DB *gorm.DB
}

func NewShowController(DB *gorm.DB) ShowController {
	return ShowController{DB}
}

func (c *ShowController) GetProductions(ctx *gin.Context) {
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

func (c *ShowController) GetPriceRanges(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	var shows []models.Show

	c.DB.Where("showtime > ?", earliest).
		Where("showtime < ?", latest).
		Joins("Production").
		Preload("Listings", "broadway = ?", true).
		Find(&shows)

	response := make(map[string]map[string]models.PriceRange)

	for _, show := range shows {
		date := utils.FormatDate(show.Showtime)
		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]models.PriceRange)
		}
		val := response[date]

		all := val["all"]
		for _, listing := range show.Listings {
			priceRange := listing.ParsedPriceRange()
			if all.Low == 0 || all.Low > priceRange.Low {
				all.Low = priceRange.Low
			}
			if all.High < priceRange.High {
				all.High = priceRange.High
			}
			val["all"] = all
		}

		// TODO: Handle itemized shows
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"priceRanges": response}})
}

func (c *ShowController) GetAverageDiscounts(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	var shows []models.Show

	c.DB.Where("showtime > ?", earliest).
		Where("showtime < ?", latest).
		Joins("Production").
		Preload("Listings", "broadway = ?", true).
		Find(&shows)

	response := make(map[string]map[string]float32)

	var totalDiscount int
	var totalListings int

	for _, show := range shows {
		date := utils.FormatDate(show.Showtime)
		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]float32)
		}
		val := response[date]

		for _, listing := range show.Listings {
			discount, _ := strconv.Atoi(listing.PercentDiscount)
			totalDiscount += discount
			totalListings++
		}
		val["all"] = float32(totalDiscount) / float32(totalListings)

		// TODO: Handle itemized shows
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"discounts": response}})
}
