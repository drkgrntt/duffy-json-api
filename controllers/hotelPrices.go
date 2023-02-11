package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelPricesController struct {
	DB *mongo.Database
}

func NewHotelPricesController(DB *mongo.Database) HotelPricesController {
	return HotelPricesController{DB}
}

func (c *HotelPricesController) GetHotelPrices(ctx *gin.Context) {
	var prices []*models.HotelPrice

	// pastWeek := time.Now().AddDate(0, 0, -7)
	date := time.Now()
	dates := []string{}

	for i := 0; i < 7; i++ {
		dates = append(dates, formatHpfDate(date.AddDate(0, 0, (-1*i))))
	}

	cursor, err := c.DB.Collection("prices").Find(context.TODO(), bson.D{{"date", bson.D{{"$in", dates}}}})
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
	cursor.Close(context.TODO())

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"prices": prices}})
}

func formatHpfDate(date time.Time) string {
	return date.Format("Mon Jan 02 2006")
	// return fmt.Sprintf("%s %s %s %s", date.Format(time.), date.Month(), date.Day(), date.Year())
}
