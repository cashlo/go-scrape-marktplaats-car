package main

import (
	"fmt"
	"os"
	"regexp"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atotto/encoding/csv"
)

type Car struct {
	Brand	string
	Model	string
	Mileage	string
	Year	string
	Price	string
	Fuel	string
	Gearbox	string
	Plate	string

	Category string
	Engine  string
	Consumption string
	Feature string
}

func FillCar(carData map[string]string) *Car {
	c := new(Car)
	brandnmodel := strings.Split(carData["Merk & Model:"], " ")
	c.Brand = brandnmodel[0]
	if (len(brandnmodel) > 1) {
		c.Model = brandnmodel[1]
	}
	c.Mileage = carData["Kilometerstand:"]
	c.Year = carData["Bouwjaar:"]
	c.Price = carData["Prijs:"]
	c.Fuel = carData["Brandstof:"]
	c.Gearbox = carData["Transmissie:"]
	c.Plate = carData["Kenteken:"]
	c.Category = carData["Carrosserie:"]
	c.Engine = carData["Motorinhoud:"]
	c.Consumption = carData["Verbruik:"]
	c.Feature = carData["Opties:"]
	return c
}

func ScrapeMarktplaats(url string, result *[]Car) {
	fmt.Printf("Scraping: %s\n", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		return
	}

  	doc.Find(".listing-title-description a").Each(func(i int, s *goquery.Selection) {
  		link, _ := s.Attr("href")

  		marktplaatsInternalLinkRegex, _ := regexp.Compile("www\\.marktplaats\\.nl")
  		if ( !marktplaatsInternalLinkRegex.MatchString(link) ) {
  			return
  		}

    	cardata := ScrapeMarktplaatsCarDetails(link)
		*result = append(*result, *FillCar(cardata))
 	})

}


func ScrapeMarktplaatsCarDetails(carUrl string) map[string]string {
	fmt.Printf("Scraping details: %s\n", carUrl)

	doc, err := goquery.NewDocument(carUrl)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	carData := make(map[string]string)


  	doc.Find(".spec").Each(func(i int, s *goquery.Selection) {

  		key := s.Find(".key").Text()
  		value := s.Find(".value").Text()

  		carData[key] = value
 	})

 	return carData
}

func main() {
	var results []Car

	baseUrl := "http://www.marktplaats.nl/z/auto-s/volkswagen.html?categoryId=157&attributes=S,10882&currentPage=%d"

	for i := 1; i <= 167; i++ {
		ScrapeMarktplaats(fmt.Sprintf(baseUrl, i), &results)
	}

	f, _ := os.Create("car.txt")
	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteStructHeader(results[0])
	w.WriteStructAll(results)
}
