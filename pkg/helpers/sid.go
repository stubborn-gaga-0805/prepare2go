package helpers

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sony/sonyflake"
)

const (
	base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Sid struct {
	sf *sonyflake.Sonyflake
}

func NewSid() *Sid {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}
	return &Sid{sf}
}

func (s Sid) GenString() (string, error) {
	// 生成分布式ID
	id, err := s.sf.NextID()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate sonyflake ID")
	}
	// 将ID转换为字符串
	return IntToBase62(int(id)), nil
}

func (s Sid) GenUint64() (uint64, error) {
	// 生成分布式ID
	return s.sf.NextID()
}

func IntToBase62(n int) string {
	if n == 0 {
		return string(base62[0])
	}

	var result []byte
	for n > 0 {
		result = append(result, base62[n%62])
		n /= 62
		fmt.Println(n)
	}

	// 反转字符串
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
		fmt.Println(result)
	}

	return string(result)
}
