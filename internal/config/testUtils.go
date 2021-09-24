package config

func initTest(serviceName string) (Config, error) {
	cfg := NewViperConfig(serviceName)
	return cfg, nil
}

func closeTest(cfg Config) error {
	return cfg.Close()
}
