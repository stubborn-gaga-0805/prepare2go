package pkg

import (
	"github.com/stubborn-gaga-0805/prepare2go/pkg/jwt"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/oss"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/wechat"
)

var pkg *Pkg

type Pkg struct {
	config Config

	Jwt     *jwt.Instance
	Oss     *oss.AliYunOss
	MiniApp *wechat.MiniApp
}

type Option func(*Pkg)

// NewPkg Pkg实例
func NewPkg(opt ...Option) *Pkg {
	var cfg = GetConfig()
	if pkg == nil {
		pkg = &Pkg{
			config: cfg,
		}
		for _, o := range opt {
			o(pkg)
		}
	}

	return pkg
}

// WithJwt 返回Jwt的实例
func WithJwt() Option {
	return func(p *Pkg) {
		if p.Jwt == nil {
			p.Jwt = jwt.NewJwtInstance(p.config.Jwt)
		}
	}
}

// WithOss 返回OSS的实例
func WithOss() Option {
	return func(p *Pkg) {
		if p.Oss == nil {
			p.Oss = oss.NewAliYunOss(p.config.Oss.AliYun)
		}
	}
}

// WithMiniApp 返回小程序的实例
func WithMiniApp() Option {
	return func(p *Pkg) {
		if p.MiniApp == nil {
			p.MiniApp = wechat.NewMiniApp(p.config.Wechat.MiniApp)
		}
	}
}
