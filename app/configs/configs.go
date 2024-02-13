package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configurations struct {
	Server ServerConfigs
	Bot    BotConfigs
	Redis  RedisConfigs
	Email  EmailConfigs
}

type ServerConfigs struct {
	Port string
}

type BotConfigs struct {
	Name  string
	Token string
}

type RedisConfigs struct {
	Hostname string
	Port     string
	Username string
	Password string
}

type EmailConfigs struct {
	Sender   string
	Receiver string
	Password string
	SmtpHost string
	SmtpPort string
}

var instantiated *Configurations = nil

func GetInstance() *Configurations {
	if instantiated == nil {
		//read the configs and put it in the struct
		viper.SetConfigName("configs")
		viper.AddConfigPath("./app/configs")
		viper.AutomaticEnv()
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file, %s", err)
		}
		if err := viper.Unmarshal(&instantiated); err != nil {
			fmt.Printf("Unable to decode into struct, %v", err)
		}
	}
	return instantiated
}
