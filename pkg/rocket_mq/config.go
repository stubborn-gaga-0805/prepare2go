package rocket_mq

import "time"

type Config struct {
	EndPoint     string        `json:"endPoint" yaml:"endPoint"`
	Group        string        `json:"group" yaml:"group"`
	Region       string        `json:"region" yaml:"region"`
	AccessKey    string        `json:"accessKey" yaml:"accessKey"`
	AccessSecret string        `json:"accessSecret" yaml:"accessSecret"`
	DialTimeout  time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
}
