package config

type JWT struct {
	SigningKey  string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`    // jwt签名
	Timeout int64  `mapstructure:"timeout" json:"timeout" yaml:"timeout"` // 过期时间
	MaxRefresh  int64  `mapstructure:"max-refresh" json:"maxRefresh" yaml:"max-refresh"`    // 缓冲时间
}
