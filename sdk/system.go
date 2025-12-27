package sdk

import (
	"github.com/gin-gonic/gin"

	"gucooing/lolo/pkg/alg"
	"gucooing/lolo/protocol/quick"
)

func systemInit(c *gin.Context) {
	req := new(quick.SystemInitRequest)
	rsp := quick.NewResponse()
	defer c.JSON(200, rsp)
	if err := alg.DecryptedData(c, &req); err != nil {
		rsp.SetError("解密失败")
		return
	}
	rsp.SetData(&quick.SystemInitResult{
		ClientIp: c.ClientIP(),
		PtConfig: &quick.PtConfig{
			UseSms: "1",
			FcmTips: &quick.FcmTips{
				NoAdultLogoutTip: "根据法规管控，当前为防沉迷管控时间，您将被强制下线。",
			},
			JoinQQGroup: new(quick.JoinQQGroup),
		},
		PtVer: &quick.PtVer{
			VersionName: "empty",
			VersionNo:   0,
			VersionUrl:  "empty",
			UpdateTime:  "empty",
			IsMust:      "empty",
			UpdateTips:  "empty",
		},
		RealnameNode: "2",
	})
}
