package models

import (
	"github.com/rajatgpt1521/cachingSystem/service/pkg/database"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)
const SIZE = 5 // size of cache
type Cache struct {
	gorm.Model
	Data string `gorm:"size:255";not null;unique`
}


func AutoMigrateSQL()  {
	if !database.Instance.HasTable(&Cache{}) && !database.Instance.HasTable("cache"){
		database.Instance.CreateTable(&Cache{})
	}else{
		database.Instance.AutoMigrate(Cache{})
	}

}
// Insert data in case not found in DB
func InsertOrUpdate(c string)(error){
	cache := Cache{Data:c}
	var count int
	if database.Instance.Where(" data = ?", c).First(&Cache{}).Count(&count); count == 0{
		if err := database.Instance.Create(&cache).Where("data",c).Error; err != nil{
			log.Error().Err(err).Msg("Unable to insert")
			return err
		}
	}
	return nil
}

// One returns a single instance record from the query.
func One(c string)(error, string){
	cache := Cache{}
	if err := database.Instance.Where("data = ?",c).First(&cache).Error; err != nil {
		log.Error().Err(err).Msg("Data not found")
		return err, ""
	}
	return nil, cache.Data
}

func All()(error, []string)  {
	var cache []Cache
	if err := database.Instance.Select("data").Limit(SIZE).Find(&cache).Error;err !=nil {
		log.Error().Err(err).Msg("Data not found")
		return err, nil
	}
	var data []string
	for _, entry := range cache{
		data = append(data, entry.Data)
	}
	return nil, data
}

