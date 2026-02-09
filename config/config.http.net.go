package config

type HttpNet struct {
	InnerIp   string `json:"InnerIp"`
	InnerPort string `json:"InnerPort"`
	Tls       bool   `json:"Tls"`
	HttpsPort string `json:"HttpsPort"`
	CertPath  string `json:"CertPath"`
	KeyPath   string `json:"KeyPath"`
}

var defaultHttpNet = &HttpNet{
	InnerIp:   "0.0.0.0",
	InnerPort: "8080",
	Tls:       true,
	HttpsPort: "4430",
	CertPath:  "./data/cert.pem",
	KeyPath:   "./data/key.pem",
}

func GetHttpNet() *HttpNet {
	conf := GetConfig()
	if conf.HttpNet == nil {
		return defaultHttpNet
	}
	return conf.HttpNet
}

func (x *HttpNet) GetInnerIp() string {
	return x.InnerIp
}

func (x *HttpNet) GetInnerPort() string {
	return x.InnerPort
}

func (x *HttpNet) GetTls() bool {
	return x.Tls
}

func (x *HttpNet) GetHttpsPort() string {
	return x.HttpsPort
}

func (x *HttpNet) GetCertFile() string {
	return x.CertPath
}

func (x *HttpNet) GetKeyFile() string {
	return x.KeyPath
}
