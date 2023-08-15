package secretariat

import (
	"bytes"
	"compress/gzip"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/logger"
	"io/ioutil"
)

// 压缩数据
func compressData(data []byte) ([]byte, error) {
	var (
		err error
		buf bytes.Buffer
		gz  = gzip.NewWriter(&buf)
	)

	if _, err = gz.Write(data); err != nil {
		logger.Helper().Errorf("CompressData[gz.Write] Error! err:[%v]", err)
		return nil, err
	}

	if err = gz.Close(); err != nil {
		logger.Helper().Errorf("CompressData[gz.Close] Error! err:[%v]", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

// 解压缩数据
func decompressData(data []byte) ([]byte, error) {
	var (
		buf = bytes.NewReader(data)
		gz  *gzip.Reader
		err error
	)

	if gz, err = gzip.NewReader(buf); err != nil {
		logger.Helper().Errorf("DecompressData[gzip.NewReader] Error! err:[%v]", err)
		return nil, err
	}
	defer gz.Close()

	return ioutil.ReadAll(gz)
}
