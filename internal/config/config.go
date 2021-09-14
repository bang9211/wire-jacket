package config

import "time"

type Config interface {
	// Load reads and sets config.
	Load() error
	// GetBool gets boolean value of the key if present, Otherwise, it gets defaultVal.
	GetBool(key string, defaultVal bool) bool
	// GetBool gets string value of the key if present, Otherwise, it gets defaultVal.
	GetString(key string, defaultVal string) string
	// GetBool gets int value of the key if present, Otherwise, it gets defaultVal.
	GetInt(key string, defaultVal int) int
	// GetBool gets int32 value of the key if present, Otherwise, it gets defaultVal.
	GetInt32(key string, defaultVal int32) int32
	// GetBool gets int64 value of the key if present, Otherwise, it gets defaultVal.
	GetInt64(key string, defaultVal int64) int64
	// GetBool gets uint value of the key if present, Otherwise, it gets defaultVal.
	GetUint(key string, defaultVal uint) uint
	// GetBool gets uint32 value of the key if present, Otherwise, it gets defaultVal.
	GetUint32(key string, defaultVal uint32) uint32
	// GetBool gets uint64 value of the key if present, Otherwise, it gets defaultVal.
	GetUint64(key string, defaultVal uint64) uint64
	// GetBool gets float64 value of the key if present, Otherwise, it gets defaultVal.
	GetFloat64(key string, defaultVal float64) float64
	// GetBool gets time.Time value of the key if present, Otherwise, it gets defaultVal.
	GetTime(key string, defaultVal time.Time) time.Time
	// GetBool gets time.Duration value of the key if present, Otherwise, it gets defaultVal.
	GetDuration(key string, defaultVal time.Duration) time.Duration
	// GetBool gets []int value of the key if present, Otherwise, it gets defaultVal.
	GetIntSlice(key string, defaultVal []int) []int
	// GetBool gets []string value of the key if present, Otherwise, it gets defaultVal.
	GetStringSlice(key string, defaultVal []string) []string
	// GetBool gets map[string]interface{} value of the key if present, Otherwise, it gets defaultVal.
	GetStringMap(key string, defaultVal map[string]interface{}) map[string]interface{}
	// GetBool gets map[string]string value of the key if present, Otherwise, it gets defaultVal.
	GetStringMapString(key string, defaultVal map[string]string) map[string]string
	// GetBool gets map[string][]string value of the key if present, Otherwise, it gets defaultVal.
	GetStringMapSlice(key string, defaultVal map[string][]string) map[string][]string
	// Close closes config.
	Close() error
}
