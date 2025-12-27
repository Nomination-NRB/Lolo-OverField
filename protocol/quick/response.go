package quick

type Response struct {
	Result bool        `json:"result"`
	Data   interface{} `json:"data"`
	Error  *Error      `json:"error"`
}

func NewResponse() *Response {
	return &Response{
		Result: true,
		Data:   nil,
		Error:  &Error{},
	}
}

type Error struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

func (r *Response) SetError(e string) {
	r.Result = false
	r.Error.Id = 1
	r.Error.Message = e
}

func (r *Response) SetData(data interface{}) {
	r.Data = data
}

type ExtInfo struct {
	OauthType   int    `json:"oauthType"`
	OauthId     string `json:"oauthId"`
	AccessToken string `json:"access_token"`
}

type UserData struct {
	Uid       string `json:"uid"`
	Username  string `json:"username"`
	Mobile    string `json:"mobile"`
	IsGuest   string `json:"isGuest"`
	RegDevice string `json:"regDevice"`
	SexType   string `json:"sexType"`
	IsMbUser  int    `json:"isMbUser"`
	IsSnsUser int    `json:"isSnsUser"`
	Token     string `json:"token"`
}

type PtConfig struct {
	UseSms      string       `json:"useSms"`
	FcmTips     *FcmTips     `json:"fcmTips"`
	JoinQQGroup *JoinQQGroup `json:"joinQQGroup"`
	ServiceInfo string       `json:"serviceInfo"`
	IsFloat     string       `json:"isFloat"`
	MainLogin   string       `json:"mainLogin"`
	UseService  string       `json:"useService"`
}

type FcmTips struct {
	NoAdultLogoutTip string `json:"noAdultLogoutTip"`
	MinorTimeTip     string `json:"minorTimeTip"`
	AgeLimitTip      string `json:"ageLimitTip"`
	AgeMaxLimitTip   string `json:"ageMaxLimitTip"`
	NoAdultCommonTip string `json:"noAdultCommonTip"`
	ShiMingTip8      string `json:"shiMingTip8"`
	ShiMingTip816    string `json:"shiMingTip8_16"`
	ShiMingTip1618   string `json:"shiMingTip16_18"`
}

type JoinQQGroup struct {
	GroupNum string `json:"groupNum"`
	GroupKey string `json:"groupKey"`
}

type PtVer struct {
	VersionName string `json:"versionName"`
	VersionNo   int    `json:"versionNo"`
	VersionUrl  string `json:"versionUrl"`
	UpdateTime  string `json:"updateTime"`
	IsMust      string `json:"isMust"`
	UpdateTips  string `json:"updateTips"`
}
