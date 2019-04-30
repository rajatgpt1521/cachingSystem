package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type dbConnParam struct {
	User     string `yaml:"database_user"`
	Password string `yaml:"database_password"`
	Host     string `yaml:"database_host"`
	Port     int    `yaml:"database_port"`
	DBName   string `yaml:"database_name"`
}

var (
	//Instance DB connection object
	Instance *gorm.DB
)

func Initialize() {
	var (
		connParam dbConnParam
		err       error
	)
	configDefaultDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	if _, err := os.Stat(filepath.Join(configDefaultDir, "service", "configs", "db_config.yaml")); os.IsNotExist(err) {
		log.Fatal().Err(err).Msg("Unable to load db_config.yaml.")
	}
	yamlFile, err := ioutil.ReadFile(filepath.Join(configDefaultDir, "service", "configs", "db_config.yaml"))

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	if err := yaml.Unmarshal(yamlFile, &connParam); err != nil {
		log.Fatal().Err(err).Msg("Unable to create db connection string.")
	}
	Instance, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		connParam.Host, connParam.Port, connParam.User, connParam.DBName, connParam.Password))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create db connection")
	}
	Instance.LogMode(true)

}


//Close closes the database connection
func Close() {
	Instance.Close()
}
