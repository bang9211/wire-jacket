package config

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bang9211/wire-jacket/internal/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const flagName = "viperjacket"
const defaultConfigFile = "app.conf"

var once sync.Once
var viperJacket *ViperJacket

type ViperJacket struct {
	viper *viper.Viper
	flag  *flag.FlagSet
}

// GetOrCreate returns new ViperJacket.
//
// Implicitly, Wire-Jacket use app.conf as default config file if exists.
// Wire-Jacket uses envfile format if there is no file extension.
// You can use the --config flag to explicitly specify a config file.
//
// The Twelve-Factors recommends use environment variables for
// configuration. Because it's good in container, cloud environments.
// However, it is not efficient to express all configs as environment
// variables.
//
// So, we adopted viper, which integrates most config formats with
// environment variables in go.
// By using viper, you can use various config file formats without
// worrying about conversion to environment variables even if it is
// not in env format.
//
// Viper v1.8.1 supports "json", "toml", "yaml", "yml", "properties",
// "props", "prop", "hcl", "dotenv", "env", "ini".
//
// But viper can't integrates all complex structure with environment
// variables. you can check it in TestLoadDefault of viperjacket_test.go.
func GetOrCreate() Config {
	if viperJacket == nil {
		once.Do(func() {
			viperJacket = &ViperJacket{viper: viper.New(), flag: flag.NewFlagSet(flagName, flag.ExitOnError)}
			viperJacket.init()
		})
	}
	return viperJacket
}

func (vc *ViperJacket) init() {
	vc.setFlags(defaultConfigFile)
	// only use 'config' flag for reading config file path
	configFilePath := vc.GetString("config", defaultConfigFile)
	vc.preconfigForRead(configFilePath)
	err := vc.Load()
	if err != nil {
		log.Printf("Failed to load config : %s", err)
	}
}

func (vc *ViperJacket) setFlags(defaultConfigFile string) {
	vc.flag.String("config", defaultConfigFile,
		"Config file(envfile)[default : "+defaultConfigFile+"]")
	pf := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	pf.AddGoFlagSet(vc.flag)
	pf.Parse(os.Args[1:])
	vc.viper.BindPFlags(pf)
}

func (vc *ViperJacket) preconfigForRead(configFilePath string) {
	vc.viper.AddConfigPath(utils.GetFileDir(configFilePath))
	ext := utils.GetFileExtension(configFilePath)
	if ext == "conf" || ext == "" {
		vc.viper.SetConfigName(utils.GetFileName(configFilePath))
		vc.viper.SetConfigType("env")
	} else {
		vc.viper.SetConfigFile(configFilePath)
	}
	vc.viper.AutomaticEnv()
}

// Load loads config file from path, if the same key exists in
// environment variables, Viper overwrites value of same key to
// environment variables. All keys are converted to lowercase and
// stored.
func (vc *ViperJacket) Load() error {
	if err := vc.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Failed to find config file default values will be used : %s", err)
			return err
		}
		log.Printf("Failed to read config file default values will be used : %s", err)
	}

	return nil
}

func (vc *ViperJacket) GetBool(key string, defaultVal bool) bool {
	if vc.viper.IsSet(key) {
		return vc.viper.GetBool(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetString(key string, defaultVal string) string {
	if vc.viper.IsSet(key) {
		v, success := expand(vc.viper.GetString(key))
		if success {
			return v
		}
	}
	v, _ := expand(defaultVal)
	return v
}

// expand replaces ${var} or $var in the string,
// if s includes no exists $var, it returns false.
func expand(s string) (string, bool) {
	var buf []byte
	success := true
	// ${} is all ASCII, so bytes are fine for this operation.
	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			if buf == nil {
				buf = make([]byte, 0, 2*len(s))
			}
			buf = append(buf, s[i:j]...)
			name, w := getShellName(s[j+1:])
			if name == "" && w > 0 {
				// Encountered invalid syntax; eat the
				// characters.
			} else if name == "" {
				// Valid syntax, but $ was not followed by a
				// name. Leave the dollar character untouched.
				buf = append(buf, s[j])
			} else {
				envValue, isExist := os.LookupEnv(name)
				if !isExist {
					success = false
				}
				buf = append(buf, envValue...)
			}
			j += w
			i = j + 1
		}
	}
	if buf == nil {
		return s, success
	}
	return string(buf) + s[i:], success
}

// getShellName returns the name that begins the string and the number of bytes
// consumed to extract it. If the name is enclosed in {}, it's part of a ${}
// expansion and two more bytes are needed than the length of the name.
func getShellName(s string) (string, int) {
	switch {
	case s[0] == '{':
		if len(s) > 2 && isShellSpecialVar(s[1]) && s[2] == '}' {
			return s[1:2], 3
		}
		// Scan to closing brace
		for i := 1; i < len(s); i++ {
			if s[i] == '}' {
				if i == 1 {
					return "", 2 // Bad syntax; eat "${}"
				}
				return s[1:i], i + 1
			}
		}
		return "", 1 // Bad syntax; eat "${"
	case isShellSpecialVar(s[0]):
		return s[0:1], 1
	}
	// Scan alphanumerics.
	var i int
	for i = 0; i < len(s) && isAlphaNum(s[i]); i++ {
	}
	return s[:i], i
}

// isShellSpecialVar reports whether the character identifies a special
// shell variable such as $*.
func isShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func (vc *ViperJacket) GetInt(key string, defaultVal int) int {
	if vc.viper.IsSet(key) {
		return vc.viper.GetInt(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetInt32(key string, defaultVal int32) int32 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetInt32(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetInt64(key string, defaultVal int64) int64 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetInt64(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetUint(key string, defaultVal uint) uint {
	if vc.viper.IsSet(key) {
		return vc.viper.GetUint(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetUint32(key string, defaultVal uint32) uint32 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetUint32(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetUint64(key string, defaultVal uint64) uint64 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetUint64(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetFloat64(key string, defaultVal float64) float64 {
	if vc.viper.IsSet(key) {
		return vc.viper.GetFloat64(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetTime(key string, defaultVal time.Time) time.Time {
	if vc.viper.IsSet(key) {
		return vc.viper.GetTime(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetDuration(key string, defaultVal time.Duration) time.Duration {
	if vc.viper.IsSet(key) {
		return vc.viper.GetDuration(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetIntSlice(key string, defaultVal []int) []int {
	if vc.viper.IsSet(key) {
		return vc.viper.GetIntSlice(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetStringSlice(key string, defaultVal []string) []string {
	if vc.viper.IsSet(key) {
		slice := vc.viper.GetStringSlice(key)
		if slice == nil {
			return []string{}
		}
		for k, v := range slice {
			slice[k] = os.ExpandEnv(v)
		}
		return slice
	}
	return defaultVal
}

func (vc *ViperJacket) GetStringMap(key string, defaultVal map[string]interface{}) map[string]interface{} {
	if vc.viper.IsSet(key) {
		return vc.viper.GetStringMap(key)
	}
	return defaultVal
}

func (vc *ViperJacket) GetStringMapString(key string, defaultVal map[string]string) map[string]string {
	if vc.viper.IsSet(key) {
		m := vc.viper.GetStringMapString(key)
		for k, v := range m {
			m[k] = os.ExpandEnv(v)
		}
		return m
	}
	for k, v := range defaultVal {
		defaultVal[k] = os.ExpandEnv(v)
	}
	return defaultVal
}

func (vc *ViperJacket) GetStringMapSlice(key string, defaultVal map[string][]string) map[string][]string {
	if vc.viper.IsSet(key) {
		return vc.viper.GetStringMapStringSlice(key)
	}
	return defaultVal
}

func (vc *ViperJacket) Close() error {
	return nil
}
