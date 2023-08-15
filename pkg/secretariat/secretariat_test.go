package secretariat

import (
	"context"
	"encoding/xml"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"path/filepath"
	"testing"
	"time"
)

// SuiteCommandData 代开发应用的指令消息结构
type SuiteCommandData struct {
	XmlName     xml.Name `xml:"xml"`
	SuiteId     string   `xml:"SuiteId,omitempty"`     // 第三方应用的SuiteId或者代开发应用模板id
	InfoType    string   `xml:"InfoType,omitempty"`    // 事件类型
	TimeStamp   int      `xml:"TimeStamp,omitempty"`   // 时间戳
	SuiteTicket string   `xml:"SuiteTicket,omitempty"` // Ticket内容，最长为512字节 (推送suite_ticket)
	AuthCode    string   `xml:"AuthCode,omitempty"`    // 临时授权码,最长为512字节。用于获取企业永久授权码。10分钟内有效 (授权成功通知, 重置永久授权码通知)
	State       string   `xml:"State,omitempty"`       // 构造授权链接指定的state参数 (授权成功通知, 变更授权通知)
	AuthCorpId  string   `xml:"AuthCorpId,omitempty"`  // 授权方的corpId (授权成功通知, 变更授权通知, 取消授权通知, 获客助手权限确认事件, 获客助手权限取消事件)
	AuthType    string   `xml:"AuthType,omitempty"`    // 此时固定为customer_acquisition (授权成功通知, 变更授权通知, 取消授权通知, 获客助手权限确认事件, 获客助手权限取消事件)
}

var (
	xmlStr  = "<xml><SuiteId><![CDATA[dk8ef1d079e0dc723e]]></SuiteId><InfoType><![CDATA[suite_ticket]]></InfoType><TimeStamp>1691464017</TimeStamp><SuiteTicket><![CDATA[OF-OjdC7ffO5yCzksTtWs7bG2MBKOV--4rndGWBocwO57WJcSi_mOu7OOuGFOPXY]]></SuiteTicket></xml>"
	xmlData SuiteCommandData
)

func TestNewSecretary(t *testing.T) {
	startAt := time.Now()
	confPath, err := filepath.Abs("../../configs/config.local.yaml")
	if err != nil {
		t.Errorf("filepath.Abs[err: %v]", err)
		t.Fail()
	}
	logger.InitLog(conf.ReadConfig(confPath).Logger, true)
	secretary := NewSecretary(context.Background(), SecretaryConfig{
		Name:        "TestNewSecretary",
		BufferSize:  64 * 1000,
		ScannerRate: time.Millisecond * 200,
		WithCompass: true,
	})
	go secretary.Start()

	for i := 1; i <= 1e5; i++ {
		if err := xml.Unmarshal([]byte(xmlStr), &xmlData); err != nil {
			t.Errorf("xml.Unmarshal[err: %v]", err)
			t.Fail()
		}
		xmlData.AuthCode = helpers.GenUUID()
		secretary.Put(&xmlData)
	}
	t.Logf("duration: %s\n", time.Since(startAt))
}
