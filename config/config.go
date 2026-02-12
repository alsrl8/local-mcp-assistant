package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Databases []DatabaseConfig `mapstructure:"databases"`
}

type DatabaseConfig struct {
	Name           string   `mapstructure:"name"`
	Host           string   `mapstructure:"host"`
	Port           int      `mapstructure:"port"`
	User           string   `mapstructure:"user"`
	Password       string   `mapstructure:"password"`
	DBName         string   `mapstructure:"dbname"`
	WritableTables []string `mapstructure:"writable_tables"`
}

var (
	instance *Config
	once     sync.Once
)

func Load(path string) (*Config, error) {
	var err error
	once.Do(func() {
		viper.SetConfigFile(path)
		if err = viper.ReadInConfig(); err != nil {
			return
		}
		instance = &Config{}
		err = viper.Unmarshal(instance)
	})
	return instance, err
}

func Get() *Config {
	return instance
}
