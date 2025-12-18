package config

// Config struct to hold configuration variables
type Config struct {
	DBSource string
}

func LoadConfig(path string) (config Config, err error) {
	// TODO: Implement config loading (e.g. from env or file)
	return
}
