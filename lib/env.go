package lib

import "github.com/spf13/viper"

type Env struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	Environment string `mapstructure:"ENVIRONMENT"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`

	DBUsername string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASS"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBType     string `mapstructure:"DB_TYPE"`

	MaxMultipartMemory int64  `mapstructure:"MAX_MULTIPART_MEMORY"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`

	TimeZone      string `mapstructure:"TIMEZONE"`
	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`
}

var globalEnv = Env{
	MaxMultipartMemory: 10 << 20, // 10 MB
}

func GetEnv() Env {
	return globalEnv
}

func NewEnv(logger Logger) *Env {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("cannot read cofiguration", err)
	}

	viper.SetDefault("TIMEZONE", "UTC")

	err = viper.Unmarshal(&globalEnv)
	if err != nil {
		logger.Fatal("environment cant be loaded: ", err)
	}

	return &globalEnv
}
