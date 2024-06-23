package env

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func init() {
	// optionally look for config in the working directory
	viper.AddConfigPath("./")
	viper.AddConfigPath("./env/")
	viper.AddConfigPath("../env/")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func GetEnv(envName string) string {
	if viper.Get(envName) == nil {
		return ""
	}
	return viper.Get(envName).(string)
}

func GetEnvBool(envName string) bool {
	if viper.Get(envName) == nil {
		return false
	}
	return viper.GetBool(envName)
}

func GetEnvInt(envName string) int {
	if viper.Get(envName) == nil {
		return 0
	}
	return viper.GetInt(envName)
}

func GetEnvFloat(envName string) float64 {
	return viper.GetFloat64(envName)
}

func GetEnvInterface(envName string) interface{} {
	if viper.Get(envName) == nil {
		return nil
	}
	return viper.Get(envName)
}

func GetTimeout(envName string) time.Duration {
	secondStr := GetEnv(envName)
	secondInt, err := strconv.Atoi(secondStr)
	if err != nil {
		return 5 * time.Second
	}
	return time.Duration(secondInt) * time.Second
}

func GetViper() *viper.Viper {
	return viper.GetViper()
}
