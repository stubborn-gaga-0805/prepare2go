package jwt

import "time"

type Jwt struct {
	SecretKey   string        `json:"secretKey" yaml:"secretKey"`
	EffectAfter time.Duration `json:"effectAfter" yaml:"effectAfter"`
	MaxAge      time.Duration `json:"maxAge" yaml:"maxAge"`
}
