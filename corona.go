package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	setupRoutes(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupRoutes(r *gin.Engine) {
	r.GET("/cases/new/country/:country", Dummy)
	r.GET("/cases/total/country/:fromDate", Dummy1)

}

func convertToTimeFormat(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		fmt.Println("error: ", err)
	}
	return t
}

//Dummy function
func Dummy(c *gin.Context) {

	records := readCsvFile("./full_data.csv")
	country, ok := c.Params.Get("country")
	date := getDate(records, country)
	if ok == false {
		res := gin.H{
			"error": "name is missing",
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res := gin.H{
		"Name":    country,
		"Date":    date,
		"Country": country,
	}
	c.JSON(http.StatusOK, res)
}

//Dummy function
func Dummy1(c *gin.Context) {
	records := readCsvFile("./full_data.csv")
	from_date_string, ok := c.Params.Get("fromDate")
	fromDate := convertToTimeFormat(from_date_string)
	totalCases := getfromDate(records, fromDate)
	if ok == false {
		res := gin.H{
			"error": "name is missing",
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res := gin.H{
		"total_case_count": totalCases,
		"FromDate":         fromDate,
	}
	c.JSON(http.StatusOK, res)
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func getDate(records [][]string, date string) []string {
	newCase := []string{}
	for i := 1; i < len(records); i++ {
		if records[i][0] == date {
			newCase = append(newCase, records[i][2])
		}
	}
	return newCase
}
func getfromDate(records [][]string, fromDate time.Time) float64 {
	totalCase := float64(0.0)
	for i := 1; i < len(records); i++ {
		date := convertToTimeFormat(records[i][0])
		if date.After(fromDate) {
			total, err := strconv.ParseFloat(records[i][4], 64)
			if err == nil {
				totalCase = totalCase + total
			}
		}
	}
	return totalCase
}
