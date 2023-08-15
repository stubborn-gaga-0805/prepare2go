package rabbitmq

import "time"

type Config struct {
	Host           string        `json:"host"`
	Port           string        `json:"port"`
	User           string        `json:"user"`
	Password       string        `json:"password"`
	Vhost          string        `json:"vhost"`
	MaxDialTimeout time.Duration `json:"maxDialTimeout"`
}
