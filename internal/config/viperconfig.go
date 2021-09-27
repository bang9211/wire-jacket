package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/bang9211/wire-jacket/internal/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ViperConfig struct {
	viper       *viper.Viper
	flag        *flag.FlagSet
	serviceName string
}

// NewViperConfig returns new ViperConfig.
// serviceName uses as default config file {serviceName}.conf.
func NewViperConfig(serviceName string) Config {
	vc := ViperConfig{viper: viper.New(), flag: flag.NewFlagSet(os.Args[0], flag.ExitOnError), serviceName: serviceName}
	vc.init()
	return &vc
}

func (vc *ViperConfig) init() {
	defaultConfigFile := vc.serviceName + ".conf"
	vc.setFlags(defaultConfigFile)
	// only use 'config' flag for reading config file path
	configFilePath := vc.GetString("config", defaultConfigFile)
	vc.preconfigForRead(configFilePath)
	vc.Load()
}

func (vc *ViperConfig) setFlags(defaultConfigFile string) {
	vc.flag.String("config", defaultConfigFile,
		"Config file(envfile)[default : "+vc.serviceName+".conf]")

	pflag.CommandLine.AddGoFlagSet(vc.flag)
	pflag.Parse()
	vc.viper.BindPFlags(pflag.CommandLine)
}

func (vc *ViperConfig) preconfigForRead(configFilePath string) {
	vc.viper.AddConfigPath(utils.GetFileDir(configFilePath))
	ext := utils.GetFileExtension(configFilePath)
	if ext == "conf" {
		vc.viper.SetConfigName(utils.GetFileName(configFilePath))
		vc.viper.SetConfigType("env")
	} else {
		supportedType := utils.IsContain(viper.SupportedExts, ext)
		if supportedType {
			vc.viper.SetConfigFile(configFilePath)
		} else {
			vc.viper.SetConfigName(utils.GetFileName(configFilePath))
		}
	}
	vc.viper.AutomaticEnv()
}

// Load loads config file from path, if the same key exists in environment variables
// Viper overwrites value of same key to environment variables.
// all the keys store to lowercase.
func (vc *ViperConfig) Load() error {
	if err := vc.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Failed to find config file default values will be used : %s", err)
		} else {
			log.Fatal(err)
		}
	}

	err := vc.viper.Unmarshal(vc)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (vc *ViperConfig) GetBool(key string, defaultVal bool) bool {
	if vc.viper.IsSet(key) {
		return vc.viper.GetBool(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetString(key string, defaultVal string) string {
	if vc.viper.IsSet(key) {
		return vc.viper.GetString(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetInt(key string, defaultVal int) int {
	if vc.viper.IsSet(key) {
		return vc.viper.GetInt(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetInt32(key string, defaultVal int32) int32 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetInt32(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetInt64(key string, defaultVal int64) int64 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetInt64(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetUint(key string, defaultVal uint) uint {
	if vc.viper.IsSet(key) {
		return vc.viper.GetUint(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetUint32(key string, defaultVal uint32) uint32 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetUint32(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetUint64(key string, defaultVal uint64) uint64 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetUint64(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetFloat64(key string, defaultVal float64) float64 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetFloat64(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetTime(key string, defaultVal time.Time) time.Time {
	if vc.viper.IsSet(key) {
		return vc.viper.GetTime(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetDuration(key string, defaultVal time.Duration) time.Duration {
	if vc.viper.IsSet(key) {
		return vc.viper.GetDuration(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetIntSlice(key string, defaultVal []int) []int {
	if vc.viper.IsSet(key) {
		return vc.viper.GetIntSlice(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetStringSlice(key string, defaultVal []string) []string {
	if vc.viper.IsSet(key) {
		ret := vc.viper.GetStringSlice(key)
		if ret == nil {
			return []string{}
		}
		return ret
	}
	return defaultVal
}

func (vc *ViperConfig) GetStringMap(key string, defaultVal map[string]interface{}) map[string]interface{} {
	if vc.viper.IsSet(key) {
		return vc.viper.GetStringMap(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetStringMapString(key string, defaultVal map[string]string) map[string]string {
	if vc.viper.IsSet(key) {
		return vc.viper.GetStringMapString(key)
	}
	return defaultVal
}

func (vc *ViperConfig) GetStringMapSlice(key string, defaultVal map[string][]string) map[string][]string {
	if vc.viper.IsSet(key) {
		return vc.viper.GetStringMapStringSlice(key)
	}
	return defaultVal
}

func (vc *ViperConfig) Close() error {
	return nil
}
