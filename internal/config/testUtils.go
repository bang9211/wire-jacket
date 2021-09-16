package config

func initTest() (Config, error) {
	cfg := NewViperConfig()
	return cfg, nil
}

func closeTest(cfg Config) error {
	return cfg.Close()
}
