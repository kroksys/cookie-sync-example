package service

import (
	"strings"

	"github.com/spf13/viper"
)

// Server struct located in server.go file is considered the config struct for this project.

// Reads the configuration file and returns pointer to Config
func ReadConfig(name string) (*Server, error) {
	name = prepareConfigFileName(name)
	viper.SetConfigName(name)
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// move config data into struct
	s := Server{}
	err = viper.Unmarshal(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Removes .json, .yaml or .toml file extensions from config file name
func prepareConfigFileName(conf string) string {
	conf = strings.TrimSuffix(conf, ".yaml")
	conf = strings.TrimSuffix(conf, ".toml")
	conf = strings.TrimSuffix(conf, ".json")
	return conf
}
