package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelPriceController struct {
	DB *mongo.Database
}

func NewHotelPriceController(DB *mongo.Database) HotelPriceController {
	return HotelPriceController{DB}
}

func (c *HotelPriceController) GetThisWeeksAverage(ctx *gin.Context) {
	var prices []*models.HotelPrice

	date := time.Now()
	dates := []string{}

	for i := 0; i < 7; i++ {
		dates = append(dates, formatHpfDate(date.AddDate(0, 0, i)))
	}

	cursor, err := c.DB.Collection("prices").Find(context.TODO(), bson.D{{Key: "date", Value: bson.D{{Key: "$in", Value: dates}}}})

	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.HotelPrice
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		prices = append(prices, &elem)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	defer cursor.Close(context.TODO())

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"prices": prices}})
}

func (c *HotelPriceController) GetTodaysAverage(ctx *gin.Context) {
	var price models.HotelPrice

	date := formatHpfDate(time.Now())

	cursor, err := c.DB.Collection("prices").
		Find(context.TODO(), bson.D{{Key: "date", Value: date}})

	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		err := cursor.Decode(&price)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	defer cursor.Close(context.TODO())

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"price": price}})
}

func (c *HotelPriceController) GetHotelPrices(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var prices []*models.HotelPrice

	date := time.Now()
	dates := []string{}

	for i := 0; i < days; i++ {
		if i < skip {
			continue
		}
		dates = append(dates, formatHpfDate(date.AddDate(0, 0, (-1*i))))
	}

	cursor, err := c.DB.Collection("prices").Find(context.TODO(), bson.D{{Key: "date", Value: bson.D{{Key: "$in", Value: dates}}}})
	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.HotelPrice
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		prices = append(prices, &elem)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	defer cursor.Close(context.TODO())

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"prices": prices}})
}

func formatHpfDate(date time.Time) string {
	return date.Format("Mon Jan 02 2006")
}
