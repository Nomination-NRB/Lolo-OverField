package quick

type SystemInitRequest struct {
	Imsi         string `json:"imsi"`
	OsLang       string `json:"os_lang"`
	ScreenWidth  string `json:"screen_width"`
	OsName       string `json:"os_name"`
	ScreenHeight string `json:"screen_height"`
	AuthToken    string `json:"auth_token"`
	GameVer      string `json:"game_ver"`
	DevImei      string `json:"dev_imei"`
	ChanelCkey   string `json:"chanel_ckey"`
	DevName      string `json:"dev_name"`
	ProductCkey  string `json:"product_ckey"`
	Platform     int    `json:"platform"`
	TimeStamp    string `json:"time_stamp"`
	Oaid         string `json:"oaid"`
	DeviceId     string `json:"device_id"`
	PushToken    string `json:"push_token"`
	CountryCode  string `json:"country_code"`
	OsVer        string `json:"os_ver"`
	SdkVer       string `json:"sdk_ver"`
}

type SystemInitResult struct {
	OrigPwd      int       `json:"origPwd"`
	ClientIp     string    `json:"clientIp"`
	PtConfig     *PtConfig `json:"pt_config"`
	PtVer        *PtVer    `json:"pt_ver"`
	RealnameNode string    `json:"realname_node"`
}
