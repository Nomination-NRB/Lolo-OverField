package game

import (
	"gucooing/lolo/game/model"
	"gucooing/lolo/pkg/alg"
	"gucooing/lolo/protocol/proto"
)

func (g *Game) SwitchGardenStatus(s *model.Player, msg *alg.GameMsg) {
	req := msg.Body.(*proto.SwitchGardenStatusReq)
	rsp := &proto.SwitchGardenStatusRsp{
		Status:             proto.StatusCode_StatusCode_Ok,
		IsOpen:             false,
		Password:           "",
		PasswordExpireTime: 0,
	}
	defer g.send(s, msg.PacketId, rsp)
	garden := model.GetSceneGardenData(s.UserId, 9999)
	garden.IsOpen = req.IsOpen
	garden.Password = req.Password
	garden.Save()
	rsp.IsOpen = garden.IsOpen
	rsp.Password = garden.Password
	rsp.PasswordExpireTime = garden.PasswordExpireTime
}

func (g *Game) GardenLikeRecord(s *model.Player, msg *alg.GameMsg) {
	// req := msg.Body.(*proto.GardenLikeRecordReq)
	rsp := &proto.GardenLikeRecordRsp{
		Status: proto.StatusCode_StatusCode_Ok,
		Record: make([]*proto.GardenLikeRecordInfo, 0),
	}
	defer g.send(s, msg.PacketId, rsp)
}

func (g *Game) GardenFurnitureScheme(s *model.Player, msg *alg.GameMsg) {
	// req := msg.Body.(*proto.GardenFurnitureSchemeReq)
	rsp := &proto.GardenFurnitureSchemeRsp{
		Status:  proto.StatusCode_StatusCode_Ok,
		Schemes: s.GetGardenModel().Schemes(),
	}
	defer g.send(s, msg.PacketId, rsp)
}

func (g *Game) GardenSchemeFurnitureList(s *model.Player, msg *alg.GameMsg) {
	req := msg.Body.(*proto.GardenSchemeFurnitureListReq)
	rsp := &proto.GardenSchemeFurnitureListRsp{
		Status:          proto.StatusCode_StatusCode_Ok,
		SchemeId:        req.SchemeId,
		FurnitureItemId: make([]uint32, 0),
	}
	defer g.send(s, msg.PacketId, rsp)
}

func (g *Game) GardenFurnitureSave(s *model.Player, msg *alg.GameMsg) {
	req := msg.Body.(*proto.GardenFurnitureSaveReq)
	rsp := &proto.GardenFurnitureSaveRsp{
		Status:   proto.StatusCode_StatusCode_Ok,
		SchemeId: req.SchemeId,
	}
	defer g.send(s, msg.PacketId, rsp)
}

func (g *Game) GardenFurnitureRemoveAll(s *model.Player, msg *alg.GameMsg) {
	// req := msg.Body.(*proto.GardenFurnitureRemoveAllReq)
	rsp := &proto.GardenFurnitureRemoveAllRsp{
		Status: proto.StatusCode_StatusCode_Ok,
	}
	defer g.send(s, msg.PacketId, rsp)
}
