package lib

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	Port        string `mapstructure:"PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	LogOutput   string `mapstructure:"LOG_OUTPUT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`

	DatabaseType string `mapstructure:"DATABASE_TYPE"`
	DatabaseUrl  string `mapstructure:"DATABASE_URL"`

	MaxMultipartMemory   int64         `mapstructure:"MAX_MULTIPART_MEMORY"`
	JWTAccessSecret      string        `mapstructure:"JWT_ACCESS_SECRET"`
	JWTRefreshSecret     string        `mapstructure:"JWT_REFRESH_SECRET"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`

	AWSAccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion          string `mapstructure:"AWS_REGION"`
	AWSS3Bucket        string `mapstructure:"AWS_S3_BUCKET"`

	TimeZone      string `mapstructure:"TIMEZONE"`
	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
}

var globalEnv *Env

func setEnv() {
	globalEnv = &Env{
		MaxMultipartMemory: 10 << 20, // 10 MB
	}

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot read cofiguration", err)
	}

	viper.SetDefault("TIMEZONE", "UTC")

	err = viper.Unmarshal(&globalEnv)
	if err != nil {
		log.Fatal("environment cant be loaded: ", err)
	}
}

func GetEnv() *Env {
	if globalEnv == nil {
		setEnv()
	}
	return globalEnv
}
