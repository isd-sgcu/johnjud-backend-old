package config

import (
	"github.com/spf13/viper"
)

type Database struct {
	Url string `mapstructure:"URL"`
}

type App struct {
	Port int    `mapstructure:"PORT"`
	Env  string `mapstructure:"ENV"`
}

type Service struct {
	File string `mapstructure:"FILE"`
}

type Config struct {
	App      App
	Database Database
	Service  Service
}

func LoadConfig() (*Config, error) {
	dbCfgLdr := viper.New()
	dbCfgLdr.SetEnvPrefix("DB")
	dbCfgLdr.AutomaticEnv()
	dbCfgLdr.AllowEmptyEnv(false)
	dbConfig := Database{}
	if err := dbCfgLdr.Unmarshal(&dbConfig); err != nil {
		return nil, err
	}

	appCfgLdr := viper.New()
	appCfgLdr.SetEnvPrefix("APP")
	appCfgLdr.AutomaticEnv()
	dbCfgLdr.AllowEmptyEnv(false)
	appConfig := App{}
	if err := appCfgLdr.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	serviceCfgLdr := viper.New()
	serviceCfgLdr.SetEnvPrefix("SERVICE")
	serviceCfgLdr.AutomaticEnv()
	dbCfgLdr.AllowEmptyEnv(false)
	serviceConfig := Service{}
	if err := serviceCfgLdr.Unmarshal(&serviceConfig); err != nil {
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
