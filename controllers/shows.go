package controllers

import (
	"fmt"
	"log"
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

	minPrice, maxPrice := utils.GetMinAndMaxPrice(ctx)
	minDiscount, maxDiscount := utils.GetMinAndMaxDiscount(ctx)

	var productions []models.Production

	c.DB.Select("id, name").
		Preload("Shows", "showtime BETWEEN ? AND ?", earliest, latest).
		Preload("Shows.Listings", "broadway = ?", true).
		Find(&productions)

	response := make(map[int]string)

	utils.ForEach(productions, func(production models.Production, i int, slice []models.Production) {
		utils.ForEach(production.Shows, func(show *models.Show, j int, shows []*models.Show) {
			show.Listings = utils.Filter(show.Listings, func(listing *models.Listing, x int, listings []*models.Listing) bool {
				return listing.IsWithinParams(minPrice, maxPrice, minDiscount, maxDiscount)
			})
		})
	})

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

func (c *ShowController) GetPerformanceTallies(ctx *gin.Context) {
	showType := ctx.Query("type")
	if showType != "musicals" && showType != "plays" {
		showType = ""
	}

	earliest, latest := utils.GetEarliestAndLatest(ctx)
	minPrice, maxPrice := utils.GetMinAndMaxPrice(ctx)
	minDiscount, maxDiscount := utils.GetMinAndMaxDiscount(ctx)

	var productions []models.Production

	c.DB.
		Preload("Shows", "showtime >= ? AND showtime < ?", earliest, latest).
		Preload("Shows.Listings", "broadway = ?", true).
		Find(&productions)

	utils.ForEach(productions, func(production models.Production, i int, slice []models.Production) {
		utils.ForEach(production.Shows, func(show *models.Show, j int, shows []*models.Show) {
			show.Listings = utils.Filter(show.Listings, func(listing *models.Listing, x int, listings []*models.Listing) bool {
				return listing.IsWithinParams(minPrice, maxPrice, minDiscount, maxDiscount)
			})
		})
	})

	broadwayProductions := utils.Filter(productions, func(production models.Production, i int, slice []models.Production) bool {
		return utils.Some(production.Shows, func(show *models.Show, j int, shows []*models.Show) bool {
			return len(show.Listings) > 0
		})
	})

	if showType != "" {
		broadwayProductions = utils.Filter(broadwayProductions, func(production models.Production, i int, slice []models.Production) bool {
			isPlay := production.Shows[0].Listings[0].IsPlayOnly
			if showType == "plays" {
				return isPlay
			} else {
				return !isPlay
			}
		})
	}

	response := utils.Reduce(broadwayProductions, func(acc map[string]int, production models.Production, i int, slice []models.Production) map[string]int {
		acc[production.Name] = len(production.Shows)
		return acc
	}, make(map[string]int))

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"performances": response}})
}

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

	minPrice, maxPrice := utils.GetMinAndMaxPrice(ctx)
	minDiscount, maxDiscount := utils.GetMinAndMaxDiscount(ctx)

	var shows []models.Show

	c.DB.Where("showtime > ?", earliest).
		Where("showtime < ?", latest).
		Joins("Production").
		Preload("Listings", "broadway = ?", true).
		Find(&shows)

	response := make(map[string]map[string]map[string]float64)

	refForAverages := make(map[string]map[string]map[string]float64)

	for _, show := range shows {
		show.Listings = utils.Filter(show.Listings, func(listing *models.Listing, i int, listings []*models.Listing) bool {
			return listing.IsWithinParams(minPrice, maxPrice, minDiscount, maxDiscount)
		})

		date := utils.FormatDate(show.Showtime)

		val, ok := response[date]
		if !ok {
			response[date] = make(map[string]map[string]float64)

			val = response[date]
			val["all"] = make(map[string]float64)
			val["all"]["low"] = 0
			val["all"]["high"] = 0

			refForAverages[date] = make(map[string]map[string]float64)
			refForAverages[date]["all"] = make(map[string]float64)
			refForAverages[date]["all"]["price"] = 0
			refForAverages[date]["all"]["showCount"] = 0

			for _, productionId := range productionIds {
				prodId := fmt.Sprint(productionId)
				response[date][prodId] = make(map[string]float64)
				response[date][prodId]["low"] = 0
				response[date][prodId]["high"] = 0

				refForAverages[date][prodId] = make(map[string]float64)
				refForAverages[date][prodId]["price"] = 0
				refForAverages[date][prodId]["showCount"] = 0
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
		showRange := make(map[string]float64)
		isShow := false

		// Handle individual productions
		for _, productionId := range productionIds {
			if show.ProductionId != productionId {
				continue
			}
			prodId := fmt.Sprint(show.ProductionId)
			showRange = val[prodId]
			isShow = true

			averageShowPrice := show.AveragePrice()
			refForAverages[date][prodId]["price"] += averageShowPrice
			refForAverages[date][prodId]["showCount"]++
			refForAverages[date]["all"]["price"] += averageShowPrice
			refForAverages[date]["all"]["showCount"]++
		}

		for _, listing := range show.Listings {
			priceRange := listing.ParsedPriceRange()
			if all["low"] == 0 || all["low"] > priceRange.Low {
				all["low"] = priceRange.Low
			}
			if all["high"] < priceRange.High {
				all["high"] = priceRange.High
			}

			if isShow {
				if showRange["low"] == 0 || showRange["low"] > priceRange.Low {
					showRange["low"] = priceRange.Low
				}
				if showRange["high"] < priceRange.High {
					showRange["high"] = priceRange.High
				}
				val[fmt.Sprint(show.ProductionId)] = showRange
			}

			val["all"] = all
		}

		if !isShow {
			log.Println(date, refForAverages[date])
			refForAverages[date]["all"]["price"] += show.AveragePrice()
			refForAverages[date]["all"]["showCount"]++
		}
	}

	for date, averages := range refForAverages {
		for id, info := range averages {
			var average float64
			if info["showCount"] > 0 {
				average = info["price"] / info["showCount"]
			}
			response[date][id]["average"] = average
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

	minPrice, maxPrice := utils.GetMinAndMaxPrice(ctx)
	minDiscount, maxDiscount := utils.GetMinAndMaxDiscount(ctx)

	var shows []models.Show

	c.DB.Where("showtime > ?", earliest).
		Where("showtime < ?", latest).
		Joins("Production").
		Preload("Listings", "broadway = ?", true).
		Find(&shows)

	response := make(map[string]map[string]float64)

	// Handle individual productions
	totalsMap := make(map[string]map[string]map[string]int)

	for _, show := range shows {
		show.Listings = utils.Filter(show.Listings, func(listing *models.Listing, i int, listings []*models.Listing) bool {
			return listing.IsWithinParams(minPrice, maxPrice, minDiscount, maxDiscount)
		})

		date := utils.FormatDate(show.Showtime)

		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]float64)
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
			response[date]["all"] = float64(totals["all"]["totalDiscount"]) / float64(totals["all"]["totalListings"])
		}

		for _, productionId := range productionIds {
			if totals[productionId]["totalListings"] == 0 {
				continue
			}
			response[date][productionId] = float64(totals[productionId]["totalDiscount"]) / float64(totals[productionId]["totalListings"])
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

	response := make(map[string]float64)
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

		percentage := float64(len(showsInRange)) / float64(total)
		response[label] = percentage * 100
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"percentages": response}})
}
