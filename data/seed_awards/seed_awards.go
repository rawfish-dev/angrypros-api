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

	jsonFile, err := os.Open("rarity.json")
	if err != nil {
		log.Fatalf("could not open rarity json file due to %s", err)
	}
	defer jsonFile.Close()

	rawData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("could not read countries json file due to %s", err)
	}

	var rarityData []models.Rarity

	err = json.Unmarshal(rawData, &rarityData)
	if err != nil {
		log.Fatalf("could not unmarshal rarity json file due to %s", err)
	}

	log.Println("processing seed rarity...")

	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&rarityData)
	if result.Error != nil {
		log.Fatalf("could not upsert rarity records due to %s", result.Error)
	}

	log.Printf("processed %d rarity", len(rarityData))

	groupedRarityData := make(map[string]models.Rarity)
	for idx, rarity := range rarityData {
		groupedRarityData[rarity.Label] = rarityData[idx]
	}

	groupedAwardsData := make(map[string][]models.Award)

	jsonFile, err = os.Open("awards.json")
	if err != nil {
		log.Fatalf("could not open awards json file due to %s", err)
	}
	defer jsonFile.Close()

	rawData, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("could not read awards json file due to %s", err)
	}

	err = json.Unmarshal(rawData, &groupedAwardsData)
	if err != nil {
		log.Fatalf("could not unmarshal awards json file due to %s", err)
	}

	log.Println("processing seed awards...")

	var awards []models.Award

	for rarityGroup, groupedAwards := range groupedAwardsData {
		if rarity, ok := groupedRarityData[rarityGroup]; ok {
			for _, groupedAward := range groupedAwards {
				groupedAward.RarityId = rarity.Id
				awards = append(awards, groupedAward)
			}
		} else {
			log.Fatalf("could not find rarity data for %s", rarityGroup)
		}
	}

	result = db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&awards)
	if result.Error != nil {
		log.Fatalf("could not upsert awards records due to %s", result.Error)
	}

	log.Printf("processed %d awards", len(awards))
}
