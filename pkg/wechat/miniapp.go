package wechat

import (
	"github.com/ArtisanCloud/PowerWeChat/v2/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v2/src/miniProgram"
	"go.uber.org/zap"
)

// MiniApp Doc: https://powerwechat.artisan-cloud.com/zh/mini-program/index.html
type MiniApp struct {
	appId string

	secret string
	token  string
	aesKey string

	Client *miniProgram.MiniProgram
}

func NewMiniApp(c MiniAppConfig) *MiniApp {
	var (
		err     error
		miniApp = &MiniApp{
			appId:  c.AppId,
			secret: c.AppSecret,
			token:  c.Token,
			aesKey: c.AESKey,
		}
		config = &miniProgram.UserConfig{
			AppID:     miniApp.appId,
			Secret:    miniApp.secret,
			AESKey:    miniApp.aesKey,
			HttpDebug: c.HttpDebug,
			Log: miniProgram.Log{
				Level: "info",
				File:  "./logs/wechat.log",
			},
			Cache: kernel.NewRedisClient(&kernel.RedisOptions{
				Addr:     c.CacheRedis.Addr,
				Password: c.CacheRedis.Password,
				DB:       c.CacheRedis.Db,
			}),
		}
	)
	if miniApp.Client, err = miniProgram.NewMiniProgram(config); err != nil {
		zap.S().Errorf("Init MiniApp Failed: conf.app.mini: %+v, config: %#v, err: %v", c, config, err)
		return nil
	}

	return miniApp
}

// GetAppId 返回AppId
func (miniApp *MiniApp) GetAppId() string {
	return miniApp.appId
}
