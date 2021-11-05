package main

import (
	"context"
	"inshorts/models"
	"inshorts/utils"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/joho/godotenv"

	"inshorts/config"

	"github.com/labstack/echo/v4"

	_ "inshorts/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// Ping
// @Summary Ping for server health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "server is healthy",
	})
}

// GetCovidData
// @Summary Get covid data from coordinates
// @Param lat query string true "latitude"
// @Param long query string true "longitude"
// @Success 200 {object} models.CovidData
// @Router /coviddata [get]
func get_covid_data(c echo.Context) error {
	lat := c.QueryParam("lat")
	long := c.QueryParam("long")

	reverse, err := utils.FetchReverseGeo(lat, long)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "internal server error, reason: " + err.Error(),
		})
	}
	if reverse.Country != "India" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "location out of india not supported",
		})
	}
	coll := config.DB.Database("inshorts").Collection("covid_data")
	var data models.CovidData
	res := coll.FindOne(context.TODO(), bson.D{{"state", reverse.Region_code}})
	res.Decode(&data)
	return c.JSON(http.StatusOK, data)
}

// @title Covid Data API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5050
// @BasePath /
// @schemes http
func main() {
	godotenv.Load()
	var err error
	config.DB, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = config.DB.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := config.DB.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	e := echo.New()

	// data := utils.FetchData()
	// var data_bson = make([]interface{}, len(data))

	// for i := 0; i < len(data); i++ {
	// 	data_bson[i] = data[i]
	// }

	// coll := config.DB.Database("inshorts").Collection("covid_data")

	// coll.InsertMany(context.TODO(), data_bson, nil)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/ping", ping)
	e.GET("/coviddata", get_covid_data)

	e.Logger.Fatal(e.Start(":5050"))
}
