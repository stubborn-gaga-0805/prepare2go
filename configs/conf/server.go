package conf

import "time"

type Server struct {
	Http Http `json:"http" yaml:"http"`
	Grpc Grpc `json:"grpc" yaml:"grpc"`
	WS   WS   `json:"ws" yaml:"ws"`
}

type Http struct {
	Network string        `json:"network" yaml:"network"`
	Addr    string        `json:"addr" yaml:"addr"`
	Port    int           `json:"port" yaml:"port"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

type Grpc struct {
	Network string        `json:"network" yaml:"network"`
	Addr    string        `json:"addr" yaml:"addr"`
	Port    string        `json:"port" yaml:"port"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

type WS struct {
	Network      string        `json:"network" yaml:"network"`
	Addr         string        `json:"addr" yaml:"addr"`
	Port         int           `json:"port" yaml:"port"`
	PingTimeout  time.Duration `json:"pingTimeout" yaml:"pingTimeout"`
	PingInterval time.Duration `json:"pingInterval" yaml:"pingInterval"`
}
