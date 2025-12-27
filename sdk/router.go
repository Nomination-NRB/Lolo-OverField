package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gucooing/lolo/config"
	"gucooing/lolo/gdconf"
	"gucooing/lolo/pkg/alg"
)

func (s *Server) Router() {
	s.router.Any("/", HandleDefault)

	dispatch := s.router.Group("/dispatch")
	{
		dispatch.POST("/region_info", regionInfo)
		dispatch.POST("/client_hot_update", clientHotUpdate)
	}
	s.router.POST("/v3/bind", bindTest)
	v1 := s.router.Group("/v1", alg.AutoCryptoMiddleware())
	{
		system := v1.Group("/system")
		{
			system.POST("/init", systemInit)
		}
		user := v1.Group("/user")
		{
			user.POST("/loginByName", s.loginByName)
			user.POST("/autoLogin", s.autoLogin)
		}
	}

	v2 := s.router.Group("/v2", alg.AutoCryptoMiddleware())
	{
		system := v2.Group("/system")
		{
			system.POST("/init", systemInit)
		}
		user := v2.Group("/user")
		{
			user.POST("/loginByName", s.loginByName)
			user.POST("/autoLogin", s.autoLogin)
		}
	}
}

func HandleDefault(c *gin.Context) {
	c.String(200, "Lolo!")
}

type RegionInfoRequest struct {
	Version         string `form:"version" binding:"required"`
	Version2        string `form:"version2" binding:"required"`
	AccountType     string `form:"accountType" binding:"required"`
	OS              string `form:"os" binding:"required"`
	LastLoginSDKUID string `form:"lastloginsdkuid" binding:"required"`
}

type RegionInfo struct {
	Status           bool   `json:"status"`
	Message          string `json:"message"`
	GateTcpIp        string `json:"gate_tcp_ip"`
	GateTcpPort      int    `json:"gate_tcp_port"`
	IsServerOpen     bool   `json:"is_server_open"`
	Text             string `json:"text"`
	ClientLogTcpIp   string `json:"client_log_tcp_ip"`
	ClientLogTcpPort int    `json:"client_log_tcp_port"`
	CurrentVersion   string `json:"currentVersion"`
	PhotoShareCdnUrl string `json:"photo_share_cdn_url"`
}

func regionInfo(c *gin.Context) {
	var req RegionInfoRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	conf := config.GetGateWay()
	info := &RegionInfo{
		Status:           true,
		Message:          "success",
		GateTcpIp:        conf.GetOuterIp(),
		GateTcpPort:      conf.GetOuterPort(),
		IsServerOpen:     true,
		Text:             "",
		ClientLogTcpIp:   config.GetLogServer().GetOuterIp(),
		ClientLogTcpPort: config.GetLogServer().GetOuterPort(),
		CurrentVersion:   gdconf.GetClientVersion(req.Version),
		PhotoShareCdnUrl: "https://cdn-photo-of.inutan.com/cn_prod_main",
	}

	c.JSONP(http.StatusOK, info)
}

type GMClientConfig struct {
	Status               bool   `json:"status"`
	Message              string `json:"message"`
	HotOssUrl            string `json:"hotOssUrl"`
	CurrentVersion       string `json:"currentVersion"`
	Server               string `json:"server"`
	SsAppId              string `json:"ssAppId"`
	SsServerUrl          string `json:"ssServerUrl"`
	OpenGm               bool   `json:"open_gm"`
	OpenErrorLog         bool   `json:"open_error_log"`
	OpenNetConnectingLog bool   `json:"open_netConnecting_log"`
	IpAddress            string `json:"ipAddress"`
	PayUrl               string `json:"payUrl"`
	IsTestServer         bool   `json:"isTestServer"`
	ErrorLogLevel        int    `json:"error_log_level"`
	ServerId             string `json:"server_id"`
	OpenCs               bool   `json:"open_cs"`
}

func clientHotUpdate(c *gin.Context) {
	var req RegionInfoRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	info := &GMClientConfig{
		Status:               true,
		Message:              "success",
		HotOssUrl:            "http://cdn-of.inutan.com/Resources;https://cdn-of.inutan.com/Resources",
		CurrentVersion:       gdconf.GetClientVersion(req.Version),
		Server:               "cn_prod_main",
		SsAppId:              "c969ebf346794cc797ed6eb6c3eac089",
		SsServerUrl:          "https://te-of.inutan.com",
		OpenGm:               true,
		OpenErrorLog:         true,
		OpenNetConnectingLog: true,
		IpAddress:            c.ClientIP(),
		PayUrl:               "http://api-callback-of.inutan.com:19701",
		IsTestServer:         true,
		ErrorLogLevel:        0,
		ServerId:             "10001",
		OpenCs:               true,
	}

	c.JSONP(http.StatusOK, info)
}

func bindTest(c *gin.Context) {
	// reqBase64, _ := io.ReadAll(c.Request.Body)
	// reqBytes, _ := base64.StdEncoding.DecodeString(string(reqBase64))
	// rsaLen := binary.BigEndian.Uint32(reqBytes[:4])
	// rsaBytes := reqBytes[4 : 4+rsaLen]
	//
	// aesLen := binary.BigEndian.Uint32(reqBytes[4+rsaLen : 4+rsaLen+4])
	// aesBytes := reqBytes[4+rsaLen+4 : 4+rsaLen+4+aesLen]
	//
	// log.Printf("rsaLen:%d rsaBytes:%s \n aesLen:%d aesBytes:%s\n",
	// 	rsaLen, hex.EncodeToString(rsaBytes), aesLen, hex.EncodeToString(aesBytes))
	//
	// rsa := flyrsa.NewFlyRSA(1024)
	// e := new(big.Int)
	// n := new(big.Int)
	// d := new(big.Int)
	// e.SetString("9cbd92ccef123be840deec0c6ed0547194c1e471d11b6f375e56038458fb18833e5bab2e1206b261495d7e2d1d9e5aa859e6d4b671a8ca5d78efede48e291a3f", 16)
	// n.SetString("a387f05b88acf4898fb76054412d552b80160e6947b00153046fab67d49d97274839358a9182c30f6df4e2cbc461ed1e3721922c4034ba2ac38fe2258ae0a9d14f032fe5068d35d097bafbb9d4c020fbf921ab0b723bcfbcb804e51a23305da0cb9112855f0f4658a69ea78106692107793e1537dc2636a014a83cbc442a709b", 16)
	// d.SetString("3da02c026e5f807f15782f5fc09ddf131a050a40c09b1e0897cde76a1183dc6950c84008cb825a426c49b0ab297b6ff3e1117a695d39d91ea49a15d93fd125d281fa64a940ef85cf2fdf9641f0fb9bb39d83b7090cb5f5c757524bd20902b4408ae86eb974a5dafaddaf8e08f552db63d3396b088608643bf88cd34adec004ff", 16)
	//
	// aesKey, err := rsa.Decode(rsaBytes, d, n)
	// if err != nil {
	// 	fmt.Printf("解密失败: %v\n", err)
	// } else {
	// 	fmt.Printf("解密结果AesKey: %s\n", string(aesKey))
	// }
	//
	// data := &flyrsa.Data{}
	// req, err := data.AES128Decode(aesKey, aesBytes)
	// if err != nil {
	// 	fmt.Printf("AES解密失败: %v\n", err)
	// } else {
	// 	fmt.Printf("AES解密结果req: %s\n", string(req))
	// }
}
