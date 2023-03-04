package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/rawfish-dev/angrypros-api/config"
	"github.com/rawfish-dev/angrypros-api/models"
)

func main() {
	appConfig := config.NewAppConfig(os.Getenv("APP_ENVIRONMENT"), "../..")

	connectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		appConfig.PostgresConfig.Host, appConfig.PostgresConfig.Port,
		appConfig.PostgresConfig.Username, appConfig.PostgresConfig.Password,
		appConfig.PostgresConfig.Database, appConfig.PostgresConfig.SSLMode)

	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{
		// Logger: logger.New(
		// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		// 	logger.Config{
		// 		SlowThreshold: time.Second, // Slow SQL threshold
		// 		LogLevel:      logger.Info, // Log level
		// 	},
		// ),
	})
	if err != nil {
		log.Fatalf("could not open db connection due to %s", err)
	}

	jsonFile, err := os.Open("countries.json")
	if err != nil {
		log.Fatalf("could not open countries json file due to %s", err)
	}
	defer jsonFile.Close()

	rawData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("could not read countries json file due to %s", err)
	}

	data := make(map[string]string)

	err = json.Unmarshal(rawData, &data)
	if err != nil {
		log.Fatalf("could not unmarshal countries json file due to %s", err)
	}

	log.Println("processing seed countries...")

	var countries []models.Country
	for isoAlpha2Code, countryName := range data {
		countries = append(countries, models.Country{
			IsoAlpha2Code: isoAlpha2Code,
			Name:          countryName,
		})
	}

	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&countries)
	if result.Error != nil {
		log.Fatalf("could not upsert country records due to %s", result.Error)
	}

	log.Printf("processed %d countries", len(countries))
}
