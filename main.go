package main

import (
	"belgianrailway/apilogic"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)
	
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("**/*.html")
	r.Static("/static", "./static")

	r.GET("/", indexHandler)
	r.GET("/stations", CallGetStationInformation)
	r.GET("/liveboards", LiveboardIndexHandler)
	r.POST("/liveboards/search", LiveboardSearchHandler)
	r.GET("/liveboards/:station", CallGetLiveboardInformation)
	r.GET("/disturbances", CallGetDisturbanceInformation)

	err := r.Run(":4000")
	if err != nil {
		log.Fatal(err)
	}
}
	
func indexHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func LiveboardIndexHandler(c *gin.Context) {
	c.HTML(200, "liveboards.html", nil)
}

func CallGetStationInformation(c *gin.Context) {
	stations, err := apilogic.GetStationInformation()
	if err != nil {
		fmt.Printf("Could not get station information: %v", err)
		// Log the error and return an HTTP 500 response
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Could not get station information: %v", err),
		})
		return
	}

	c.HTML(200, "stations.html", gin.H{
		"Stations": stations,
	})
}

func CallGetLiveboardInformation(c *gin.Context) {
	station := c.Param("station")
	if station == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "station parameter is required"})
		return
	}

	liveboard, err := apilogic.GetLiveboardInformation(station)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "liveboardsinfo.html", gin.H{
		"Station":    liveboard.Station,
		"Departures": liveboard.Departures.DepartureInfo,
	})
}

func LiveboardSearchHandler(c *gin.Context) () {
	station := c.PostForm("station") // Get station from form data
	if station == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"Error": "Station name is required",
		})
		return
	}

	// Redirect to /liveboards/:station
	c.Redirect(http.StatusFound, fmt.Sprintf("/liveboards/%s", station))
}

func CallGetDisturbanceInformation(c *gin.Context) {
	info, err := apilogic.GetDisturbanceInformation()
	if err != nil {
		fmt.Printf("Could not get disturbance information: %v", err)
		// Log the error and return an HTTP 500 response
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Could not get disturbance information: %v", err),
		})
		return
	}

	c.HTML(200, "disturbances.html", gin.H{
		"Disturbances": info,
		"Date":         time.Now().In(time.FixedZone("CEST", 2*3600)).Format("2006-01-02"),
	})
}