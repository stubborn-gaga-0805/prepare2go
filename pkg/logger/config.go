package logger

type Config struct {
	Zap Zap `json:"zap" yaml:"zap"`
}

type Zap struct {
	Mode       string `json:"mode" yaml:"mode"`
	FilePath   string `json:"filePath" yaml:"filePath"`
	FileName   string `json:"fileName" yaml:"fileName"`
	MaxSize    int    `json:"maxSize" yaml:"maxSize"`
	MaxAge     int    `json:"maxAge" yaml:"maxAge"`
	MaxBackups int    `json:"maxBackups" yaml:"maxBackups"`
}
