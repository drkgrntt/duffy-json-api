package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
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
		Preload("Shows", "showtime > ?", pastWeek).
		Preload("Shows.Listings", "scanned_at > ? AND broadway = ?", pastWeek, true).
		Preload("CompetitionGroup.Productions").
		Order("last_scanned_at DESC").
		Order("last_shown_at DESC").
		Find(&productions)

	var response []models.Production
	for _, production := range productions {
		isBroadway := false
		for _, show := range production.Shows {
			if len(show.Listings) > 0 {
				isBroadway = true
			}
		}
		if isBroadway {
			response = append(response, production)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"shows": response}})
}

func (c *ShowController) GetNames(ctx *gin.Context) {
	earliest, latest := utils.GetEarliestAndLatest(ctx)

	showType := ctx.Query("type")
	if showType != "musicals" && showType != "plays" {
		showType = ""
	}

	var productions []models.Production

	c.DB.Select("id, name").
		Preload("Shows", "showtime BETWEEN ? AND ?", earliest, latest).
		Preload("Shows.Listings", "broadway = ?", true).
		Find(&productions)

	response := make(map[int]string)

	for _, production := range productions {
		include := false
		for _, show := range production.Shows {
			if len(show.Listings) > 0 {
				if showType == "plays" {
					include = show.Listings[0].IsPlayOnly
				} else if showType == "musicals" {
					include = !show.Listings[0].IsPlayOnly
				} else {
					include = true
				}
				break
			}
		}
		if include {
			response[production.Id] = production.Name
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"shows": response}})
}

// =============================================== //
//                     Tallies                     //
// =============================================== //

func (c *ShowController) GetPriceRangeTallies(ctx *gin.Context) {
	showType := ctx.Query("type")
	if showType != "musicals" && showType != "plays" {
		showType = ""
	}

	earliest, latest := utils.GetEarliestAndLatest(ctx)
	productionIdsQuery := ctx.QueryArray("productionIds")
	productionIds := make([]int, 0)
	for _, productionId := range productionIdsQuery {
		id, err := strconv.Atoi(productionId)
		if err != nil {
			continue
		}
		productionIds = append(productionIds, id)
	}

	var shows []models.Show

	c.DB.Where("showtime > ?", earliest).
		Where("showtime < ?", latest).
		Joins("Production").
		Preload("Listings", "broadway = ?", true).
		Find(&shows)

	response := make(map[string]map[string]models.PriceRange)

	for _, show := range shows {
		date := utils.FormatDate(show.Showtime)

		val, ok := response[date]
		if !ok {
			response[date] = make(map[string]models.PriceRange)
			val = response[date]
			for _, productionId := range productionIds {
				prodId := fmt.Sprint(productionId)
				response[date][prodId] = models.PriceRange{}
			}
		}

		include := false
		if len(show.Listings) > 0 {
			if showType == "plays" {
				include = show.Listings[0].IsPlayOnly
			} else if showType == "musicals" {
				include = !show.Listings[0].IsPlayOnly
			} else {
				include = true
			}
		}
		if !include {
			continue
		}

		all := val["all"]
		var showRange models.PriceRange
		isShow := false

		// Handle individual productions
		for _, productionId := range productionIds {
			if show.ProductionId != productionId {
				continue
			}
			showRange = val[fmt.Sprint(show.ProductionId)]
			isShow = true
		}

		for _, listing := range show.Listings {
			priceRange := listing.ParsedPriceRange()
			if all.Low == 0 || all.Low > priceRange.Low {
				all.Low = priceRange.Low
			}
			if all.High < priceRange.High {
				all.High = priceRange.High
			}

			if isShow {
				if showRange.Low == 0 || showRange.Low > priceRange.Low {
					showRange.Low = priceRange.Low
				}
				if showRange.High < priceRange.High {
					showRange.High = priceRange.High
				}
				val[fmt.Sprint(show.ProductionId)] = showRange
			}

			val["all"] = all
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"priceRanges": response}})
}

func (c *ShowController) GetAverageDiscountTallies(ctx *gin.Context) {
	showType := ctx.Query("type")
	if showType != "musicals" && showType != "plays" {
		showType = ""
	}

	earliest, latest := utils.GetEarliestAndLatest(ctx)
	productionIds := ctx.QueryArray("productionIds")

	var shows []models.Show

	c.DB.Where("showtime > ?", earliest).
		Where("showtime < ?", latest).
		Joins("Production").
		Preload("Listings", "broadway = ?", true).
		Find(&shows)

	response := make(map[string]map[string]float32)

	// Handle individual productions
	totalsMap := make(map[string]map[string]map[string]int)

	for _, show := range shows {
		date := utils.FormatDate(show.Showtime)

		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]float32)
			val := response[date]
			val["all"] = 0

			totalsMap[date] = make(map[string]map[string]int)

			totalsMap[date]["all"] = make(map[string]int)
			totalsMap[date]["all"]["totalDiscount"] = 0
			totalsMap[date]["all"]["totalListings"] = 0

			for _, productionId := range productionIds {
				val[productionId] = 0
				totalsMap[date][productionId] = make(map[string]int)
				totalsMap[date][productionId]["totalDiscount"] = 0
				totalsMap[date][productionId]["totalListings"] = 0
			}
		}

		include := false
		if len(show.Listings) > 0 {
			if showType == "plays" {
				include = show.Listings[0].IsPlayOnly
			} else if showType == "musicals" {
				include = !show.Listings[0].IsPlayOnly
			} else {
				include = true
			}
		}
		if !include {
			continue
		}

		allTotals := totalsMap[date]["all"]

		isShow := false
		var productionTotals map[string]int

		// Handle individual productions
		for _, productionId := range productionIds {
			if fmt.Sprint(show.ProductionId) != productionId {
				continue
			}
			isShow = true
			productionTotals = totalsMap[date][productionId]
		}

		for _, listing := range show.Listings {
			discount, _ := strconv.Atoi(listing.PercentDiscount)
			allTotals["totalDiscount"] += discount
			allTotals["totalListings"]++

			if isShow {
				productionTotals["totalDiscount"] += discount
				productionTotals["totalListings"]++
			}
		}
	}

	for date := range response {
		totals := totalsMap[date]

		if totals["all"]["totalListings"] != 0 {
			response[date]["all"] = float32(totals["all"]["totalDiscount"]) / float32(totals["all"]["totalListings"])
		}

		for _, productionId := range productionIds {
			if totals[productionId]["totalListings"] == 0 {
				continue
			}
			response[date][productionId] = float32(totals[productionId]["totalDiscount"]) / float32(totals[productionId]["totalListings"])
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"discounts": response}})
}

func (c *ShowController) GetPercentageAtTktsTallies(ctx *gin.Context) {
	earliest, latest := utils.GetEarliestAndLatest(ctx)

	var shows []models.Show
	var grosses []models.Gross

	c.DB.Where("created_at > ?", earliest).
		Where("created_at < ?", latest).
		Preload("Listings", "broadway = ?", true).
		Find(&shows)

	c.DB.Where("week_end_date > ?", earliest).
		Where("week_end_date < ?", latest).
		Find(&grosses)

	tmp := make(map[time.Time]int)
	for _, gross := range grosses {
		tmp[gross.WeekEndDate] += gross.Performances
	}

	response := make(map[string]float32)
	for weekEndDate, total := range tmp {
		start := now.With(weekEndDate).BeginningOfWeek()
		end := now.With(weekEndDate).EndOfWeek()
		label := fmt.Sprintf("%s to %s", utils.FormatDate(start), utils.FormatDate(end))

		var showsInRange []models.Show
		for _, show := range shows {
			if (show.CreatedAt.Equal(start) || show.CreatedAt.After(start)) &&
				show.CreatedAt.Before(end) &&
				len(show.Listings) > 0 {

				showsInRange = append(showsInRange, show)
			}
		}

		percentage := float32(len(showsInRange)) / float32(total)
		response[label] = percentage * 100
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"percentages": response}})
}
