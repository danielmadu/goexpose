package config

import (
	"encoding/json"
)

type MessageInterface interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
	GetPath() string
	GetMethod() string
	GetStatus() int
}

type Config struct {
	Shared   string `json:"shared"`
	Token    string `json:"-"`
	Client   string `json:"client"`
	CertFile string
	KeyFile  string
}

type Ping struct {
	Port string `json:"port"`
}

type LocalConfig struct {
	SharedHostname   string `json:"shared_hostname"`
	BasicAuth        string
	EnabledBasicAuth bool
}

type Message struct {
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Method  string            `json:"method"`
	Status  int               `json:"status"`
}

func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Decode(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Message) GetPath() string {
	return m.Path
}

func (m *Message) GetMethod() string {
	return m.Method
}

func (m *Message) GetStatus() int {
	return m.Status
}

var (
	config      *Config
	localConfig *LocalConfig
	message     *Message
	channel     chan Message
)

func Init() {
	config = &Config{}
	localConfig = &LocalConfig{}
	message = &Message{}
	channel = make(chan Message)
}

func GetConfig() *Config {
	return config
}

func GetLocalConfig() *LocalConfig {
	return localConfig
}

func GetLocalMessage() *Message {
	return message
}

func GetChannel() chan Message {
	return channel
}
