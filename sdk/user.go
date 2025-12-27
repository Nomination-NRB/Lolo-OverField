package sdk

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"gucooing/lolo/db"
	"gucooing/lolo/pkg/alg"
	"gucooing/lolo/pkg/log"
	"gucooing/lolo/protocol/quick"
)

func getLoginResult(user *db.OFQuick) *quick.LoginResult {
	result := &quick.LoginResult{
		ExtInfo:       nil,
		IsAdult:       true,
		UAge:          9999,
		CkPlayTime:    0,
		GuestRealName: 1,
		Id:            0,
		Message:       "",
		AuthToken:     user.AuthToken,
		UserData:      getUserData(user),
		CheckRealname: 0,
	}
	return result
}

func getUserData(user *db.OFQuick) *quick.UserData {
	data := &quick.UserData{
		Uid:       strconv.FormatUint(uint64(user.ID), 10),
		Username:  user.Username,
		Mobile:    "188****8888",
		IsGuest:   "",
		RegDevice: user.RegDevice,
		SexType:   "",
		IsMbUser:  1,
		IsSnsUser: 0,
		Token:     user.UserToken,
	}
	return data
}

func loginByName(c *gin.Context) {
	req := new(quick.LoginByNameRequest)
	rsp := quick.NewResponse()
	defer c.JSON(200, rsp)
	if err := alg.DecryptedData(c, &req); err != nil {
		rsp.SetError("解密失败")
		log.App.Debugf("gin req autoLogin error: %v", err)
		return
	}
	user, err := db.GetOFQuick(req.Username, req.Password)
	if err != nil {
		rsp.SetError("解密失败")
		return
	}
	if user.Password != req.Password {
		rsp.SetError("解密失败")
		return
	}

	// 更新token
	if user.AuthToken == "" {
		user.AuthToken = alg.RandStr(40)
	}
	user.UserToken = alg.RandStr(40)
	if err := db.UpOFQuick(user); err != nil {
		rsp.SetError("更新失败")
		return
	}

	rsp.SetData(getLoginResult(user))
}

func autoLogin(c *gin.Context) {
	req := new(quick.AutoLoginRequest)
	rsp := quick.NewResponse()
	defer c.JSON(200, rsp)
	if err := alg.DecryptedData(c, &req); err != nil {
		rsp.SetError("解密失败")
		log.App.Debugf("gin req autoLogin error: %v", err)
		return
	}
	user, err := db.GetOFQuickByAuthToken(req.AuthToken)
	if err != nil {
		rsp.SetError("没有该账号")
		return
	}
	// 更新token
	// user.AuthToken = alg.RandStr(40)
	user.UserToken = alg.RandStr(40)
	if err := db.UpOFQuick(user); err != nil {
		rsp.SetError("更新失败")
		return
	}

	rsp.SetData(getLoginResult(user))
}
