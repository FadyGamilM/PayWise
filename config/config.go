package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type PG_Config struct {
	Postgresdb struct {
		Host     string
		User     string
		Password string
		Dbname   string
		Sslmode  string
	}
}

type PG_Test_Config struct {
	Testpostgresdb struct {
		Host     string
		User     string
		Password string
		Dbname   string
		Sslmode  string
	}
}

func LoadPostgresConfig(path string) (*PG_Config, error) {
	config := new(PG_Config)

	// tell viper from where to read
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// configure the feature to override the vars from the yaml file via the environment variables
	viper.AutomaticEnv()
	/*
		 viper reads the vars from the yaml file as following :
			SERVER.PORT , but we can't define an env variable with the dot notation, so we will define it with _ and replace the default behaviour of viper
	*/
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	// read the configs
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("error reading data from config file : %v \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("error unmarshling the data from config file : %v \n", err)
	}
	return config, nil
}

func LoadPostgresTestConfig() (*PG_Test_Config, error) {
	config := new(PG_Test_Config)

	// tell viper from where to read
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")

	// configure the feature to override the vars from the yaml file via the environment variables
	viper.AutomaticEnv()
	/*
		 viper reads the vars from the yaml file as following :
			SERVER.PORT , but we can't define an env variable with the dot notation, so we will define it with _ and replace the default behaviour of viper
	*/
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	// read the configs
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("error reading data from config file : %v \n", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("error unmarshling the data from config file : %v \n", err)
	}
	return config, nil
}
