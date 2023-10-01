package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type PGConfig struct {
	Postgresdb struct {
		Host     string
		User     string
		Password string
		Dbname   string
		Sslmode  string
	}
}

type PGTestConfig struct {
	Testpostgresdb struct {
		Host     string
		User     string
		Password string
		Dbname   string
		Sslmode  string
	}
}

type PasetoTokenConfig struct {
	Paseto struct {
		SymmetricKey             string
		Access_token_expiration  time.Duration
		Refresh_token_expiration time.Duration
	}
}

func LoadPasetoTokenConfig(path string) (*PasetoTokenConfig, error) {
	viper.Reset()
	config := new(PasetoTokenConfig)

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

	// Add this after unmarshaling the configuration
	return config, nil
}

func LoadPostgresConfig(path string) (*PGConfig, error) {
	viper.Reset()
	config := new(PGConfig)

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

	log.Println("Environment Variable:", viper.Get("PGCONFIG_POSTGRESDB_FADY"))

	// If you've set the environment variable correctly, it will override the value from config.yaml
	log.Printf("Config Value: %s", config.Postgresdb.Host)

	return config, nil
}

func LoadPostgresConfig_v2() (*PGConfig, error) {
	viper.Reset()
	config := new(PGConfig)

	// tell viper from where to read
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file : %v", err)
	}

	if hostOverrided := os.Getenv("POSTGRESDB_HOST"); hostOverrided != "" {
		viper.Set("postgresdb.host", hostOverrided)
	}
	config.Postgresdb.Host = viper.GetString("postgresdb.host")

	if userOverrided := os.Getenv("POSTGRESDB_USER"); userOverrided != "" {
		viper.Set("postgresdb.user", userOverrided)
	}
	config.Postgresdb.User = viper.GetString("postgresdb.user")

	if passwordOverrided := os.Getenv("POSTGRESDB_PASSWORD"); passwordOverrided != "" {
		viper.Set("postgresdb.password", passwordOverrided)
	}
	config.Postgresdb.Host = viper.GetString("postgresdb.password")

	if dbNameOverrided := os.Getenv("POSTGRESDB_DBNAME"); dbNameOverrided != "" {
		viper.Set("postgresdb.dbname", dbNameOverrided)
	}
	config.Postgresdb.Host = viper.GetString("postgresdb.dbname")

	if sslemodeOverrided := os.Getenv("POSTGRESDB_SSLMODE"); sslemodeOverrided != "" {
		viper.Set("postgresdb.sslmode", sslemodeOverrided)
	}
	config.Postgresdb.Host = viper.GetString("postgresdb.host")
	config.Postgresdb.User = viper.GetString("postgresdb.user")
	config.Postgresdb.Password = viper.GetString("postgresdb.password")
	config.Postgresdb.Dbname = viper.GetString("postgresdb.dbname")
	config.Postgresdb.Sslmode = viper.GetString("postgresdb.sslmode")

	// If you've set the environment variable correctly, it will override the value from config.yaml
	log.Printf("Config Value: %s", viper.GetString("postgresdb.host"))

	return config, nil
}

func LoadPostgresTestConfig() (*PGTestConfig, error) {
	config := new(PGTestConfig)

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

type ServerConfig struct {
	Server struct {
		Port string
	}
}

func LoadServerConfigs(path string) (*ServerConfig, error) {
	viper.Reset()
	config := new(ServerConfig)

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

	// Add this after unmarshaling the configuration
	return config, nil
}
