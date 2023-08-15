package oss

type Oss struct {
	AliYun AliYun `json:"aliYun" yaml:"aliYun"`
}

type AliYun struct {
	AccessId      string `json:"accessId" yaml:"accessId"`
	AccessSecret  string `json:"accessSecret" yaml:"accessSecret"`
	DefaultBucket string `json:"defaultBucket" yaml:"defaultBucket"`
	EndPoint      string `json:"endPoint" yaml:"endPoint"`
	CustomDomain  string `json:"customDomain" yaml:"customDomain"`
}
