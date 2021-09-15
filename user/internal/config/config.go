package config

type Config struct {
	Server struct {
		Http struct {
			Port int
		}
		Grpc struct {
			Port int
		}
	}
}
