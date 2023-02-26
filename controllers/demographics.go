package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DemographicController struct {
	DB *gorm.DB
}

func NewDemographicController(DB *gorm.DB) DemographicController {
	return DemographicController{DB}
}

func (c *DemographicController) CreateAnalytic(ctx *gin.Context) {
	ip := ctx.ClientIP()
	// log.Println(ip, ctx.RemoteIP(), ctx.ClientIP())
	userAgent := ctx.Request.UserAgent()
	domain := ctx.Request.Referer()

	config := utils.GetConfig()

	url := fmt.Sprintf("https://ipinfo.io/%s?token=%s", ip, config.IpInfoKey)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var ipData map[string]interface{}
	err = json.Unmarshal(body, &ipData)
	if err != nil {
		log.Fatalln(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var payload map[string]string
	if err := ctx.BindJSON(&payload); err != nil {
		log.Fatalln(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var country string
	icountry := ipData["country"]
	if icountry != nil {
		country = fmt.Sprint(icountry)
	}
	var city string
	icity := ipData["city"]
	if icity != nil {
		city = fmt.Sprint(icity)
	}
	var state string
	istate := ipData["region"]
	if istate != nil {
		state = fmt.Sprint(istate)
	}
	analytic := models.Analytic{
		Page:      payload["page"],
		Query:     payload["query"],
		UserAgent: userAgent,
		Country:   country,
		City:      city,
		State:     state,
		Ip:        ip,
		Domain:    domain,
	}

	c.DB.Create(&analytic)

	// ctx.Status(http.StatusAccepted)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"info": analytic}})
}

type GetTalliesResponse struct {
	Domestic      uint `json:"domestic"`
	International uint `json:"international"`
}

func (c *DemographicController) GetTallies(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var analytics []models.Analytic
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Select("country, created_at").
		Where("created_at > ?", earliest).
		Where("created_at < ?", latest).
		Find(&analytics)

	response := make(map[string]GetTalliesResponse)

	for _, analytic := range analytics {
		date := utils.FormatDate(analytic.CreatedAt)
		_, ok := response[date]
		if !ok {
			response[date] = GetTalliesResponse{}
		}
		val := response[date]

		switch analytic.Country {
		case "US":
			val.Domestic++
		default:
			val.International++
		}

		response[date] = val
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"demographics": response}})
}

func (c *DemographicController) GetDomesticTallies(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var analytics []models.Analytic
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Select("state, created_at").
		Where("country = ?", "US").
		Where("created_at > ?", earliest).
		Where("created_at < ?", latest).
		Find(&analytics)

	response := make(map[string]map[string]uint)

	for _, analytic := range analytics {
		date := utils.FormatDate(analytic.CreatedAt)
		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]uint)
		}
		val := response[date]

		val[analytic.State]++

		response[date] = val
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"demographics": response}})
}

func (c *DemographicController) GetInternationalTallies(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var analytics []models.Analytic
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Select("country, created_at").
		// Where("country != ?", "US").
		Where("created_at > ?", earliest).
		Where("created_at < ?", latest).
		Find(&analytics)

	response := make(map[string]map[string]uint)

	for _, analytic := range analytics {
		date := utils.FormatDate(analytic.CreatedAt)
		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]uint)
		}
		val := response[date]

		val[utils.GetCountryFromCode(analytic.Country)]++

		response[date] = val
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"demographics": response}})
}
