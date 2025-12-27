package quick

type LoginByNameRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ProductCkey string `json:"product_ckey"`
	DeviceId    string `json:"device_id"`
	Platform    int    `json:"platform"`
	TimeStamp   string `json:"time_stamp"`
	ChanelCkey  string `json:"chanel_ckey"`
	SdkVer      string `json:"sdk_ver"`
	GameVer     string `json:"game_ver"`
	AuthToken   string `json:"auth_token"`
	Oaid        string `json:"oaid"`
}

type AutoLoginRequest struct {
	ProductCkey string `json:"product_ckey"`
	DeviceId    string `json:"device_id"`
	Platform    int    `json:"platform"`
	TimeStamp   string `json:"time_stamp"`
	ChanelCkey  string `json:"chanel_ckey"`
	SdkVer      string `json:"sdk_ver"`
	GameVer     string `json:"game_ver"`
	AuthToken   string `json:"auth_token"`
	Oaid        string `json:"oaid"`
}

type LoginResult struct {
	ExtInfo       *ExtInfo  `json:"extInfo"`
	IsAdult       bool      `json:"isAdult"`
	UAge          int       `json:"uAge"`
	CkPlayTime    int       `json:"ckPlayTime"`
	GuestRealName int       `json:"guestRealName"`
	Id            int       `json:"id"`
	Message       string    `json:"message"`
	AuthToken     string    `json:"auth_token"`
	UserData      *UserData `json:"user_data"`
	CheckRealname int       `json:"check_realname"`
}
