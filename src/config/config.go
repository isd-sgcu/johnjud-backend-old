// package config

// import (
// 	"github.com/pkg/errors"
// 	"github.com/spf13/viper"
// )

// type Database struct {
// 	Host     string `mapstructure:"host"`
// 	Port     int    `mapstructure:"port"`
// 	Name     string `mapstructure:"name"`
// 	Username string `mapstructure:"username"`
// 	Password string `mapstructure:"password"`
// 	SSL      string `mapstructure:"ssl"`
// }

// type App struct {
// 	Port  int  `mapstructure:"port"`
// 	Debug bool `mapstructure:"debug"`
// }

// type Service struct {
// 	File string `mapstructure:"file"`
// }

// type Config struct {
// 	App      App      `mapstructure:"app"`
// 	Database Database `mapstructure:"database"`
// 	Service  Service  `mapstructure:"service"`
// }

// func LoadConfig() (config *Config, err error) {
// 	viper.AddConfigPath("./config")
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")

// 	viper.AutomaticEnv()

// 	err = viper.ReadInConfig()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "error occurs while reading the config")
// 	}

// 	err = viper.Unmarshal(&config)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "error occurs while unmarshal the config")
// 	}

// 	return
// }

package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Database struct {
	Url string `mapstructure:"db_url"`
}

type App struct {
	Port int    `mapstructure:"app_port"`
	Env  string `mapstructure:"app_env"`
}

type Service struct {
	File string `mapstructure:"service_file"`
}

type Config struct {
	App      App
	Database Database
	Service  Service
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).
			Str("service", "file").
			Msg("Failed to load .env file")
	}

	var dbConfig Database
	if err := viper.Unmarshal(&dbConfig); err != nil {
		return nil, err
	}

	var appConfig App
	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	var serviceConfig Service
	if err := viper.Unmarshal(&serviceConfig); err != nil {
		return nil, err
	}

	config := &Config{
		Database: dbConfig,
		App:      appConfig,
		Service:  serviceConfig,
	}

	return config, nil
}

func (ac *App) IsDevelopment() bool {
	return ac.Env == "development"
}
