package config

type Config struct {
	Database struct {
		Address  string
		Name     string
		User     string
		Password string
	}
	RunMode string `mapstructure:"RUN_MODE"`
	Server  struct {
		Http struct {
			Port int
		}
		Grpc struct {
			Port int
		}
	}
}
