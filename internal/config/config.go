package config

import (
	"homework_platform/internal/utils"
	// "fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var CONFIG_FILE_PATH = "./config.yml"

// Config ...
type Config struct {
	JWTSigningString string `yaml:"jwt_signing_string"`
	SQLDSN           string `yaml:"sql_dsn"`
}

func DefaultACHConfig() *Config {
	return &Config{
		JWTSigningString: utils.RandStr(6),
		SQLDSN:           "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
	}
}

func ReadConfig() (*Config, error) {
	// log.Println("[config/ReadConfig]: Reading " + CONFIG_FILE_PATH + "...")
	configStr, err := os.ReadFile(CONFIG_FILE_PATH)
	if err != nil { // 读取文件发生错误
		return DefaultACHConfig(), err
	}

	config := &Config{}

	// 可以读取config.yml，清空 config
	// log.Println("[config/ReadConfig]: Parsing...")
	err = yaml.Unmarshal(configStr, config)
	if err != nil {
		// log.Println(err)
	}
	// log.Print("[config/ReadConfig]: config:", config, '\n')
	return config, nil
}

func Save(config *Config) {
	// log.Println("[config/SaveConfig]: Saving config to " + CONFIG_FILE_PATH + "...")
	configStr, _ := yaml.Marshal(config)
	os.WriteFile(CONFIG_FILE_PATH, configStr, 0666)
}
