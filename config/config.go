package config

type Config struct {
	Shared string `json:"shared"`
	Token  string `json:"-"`
	Client string `json:"client"`
}

type Ping struct {
	Port string `json:"port"`
}

var (
	config *Config
)

func Init() {
	config = &Config{}
}

func GetConfig() *Config {
	return config
}
