package model

import (
	"fmt"
	"time"

	"gucooing/lolo/db"
	"gucooing/lolo/gdconf"
	"gucooing/lolo/pkg/alg"
	"gucooing/lolo/pkg/cache"
	"gucooing/lolo/pkg/log"
	"gucooing/lolo/protocol/proto"
)

const (
	ChannelTypePublic = iota
	ChannelTypePrivate
)

var (
	sceneGardenCache   = cache.New[string, *SceneGardenData](5 * time.Second) // 私人花园仅读
	furnitureSnowflake = alg.NewSnowflakeWorker(1, 1765900800000)             // 家具唯一id生成器
)

func NextFurnitureId() int64 {
	return furnitureSnowflake.GenId()
}

// 房间花园数据
type SceneGardenData struct {
	SceneId                     uint32                           `json:"sceneId;omitempty"`         // 场景
	GardenName                  string                           `json:"gardenName;omitempty"`      // 花园名称
	LikesNum                    int64                            `json:"likesNum;omitempty"`        // 点赞数
	AccessPlayerNum             int64                            `json:"accessPlayerNum;omitempty"` // 访问数
	LeftLikeNum                 uint32                           `json:"leftLikeNum;omitempty"`
	IsOpen                      bool                             `json:"isOpen;omitempty"`                 // 是否开放
	Password                    string                           `json:"password"`                         // 密码
	GardenFurnitureInfoMap      map[int64]*FurnitureDetailsInfo  `json:"gardenFurnitureInfoMap;omitempty"` // 主人家具信息
	OtherPlayerFurnitureInfoMap map[uint32]*FurnitureDetailsInfo `json:"-"`                                // 客人家具信息

}

func GetSceneGardenData(userId, sceneId uint32, channelType int) *SceneGardenData {
	if sceneId != 9999 {
		return &SceneGardenData{
			SceneId:                     sceneId,
			GardenFurnitureInfoMap:      make(map[int64]*FurnitureDetailsInfo),
			OtherPlayerFurnitureInfoMap: make(map[uint32]*FurnitureDetailsInfo),
		}
	}
	data, ok := sceneGardenCache.Get(fmt.Sprintf("%v|%v", userId, sceneId))
	if channelType == ChannelTypePublic {
		if !ok {
			data = &SceneGardenData{
				SceneId:                     sceneId,
				GardenFurnitureInfoMap:      make(map[int64]*FurnitureDetailsInfo),
				OtherPlayerFurnitureInfoMap: make(map[uint32]*FurnitureDetailsInfo),
			}
			sceneGardenCache.Set(fmt.Sprintf("%v|%v", userId, sceneId), data)
		}
	} else {
		if !ok {
			home, err := db.GetOFHome(userId, sceneId)
			if err != nil {
				log.Game.Errorf("UserId:%v SceneId:%v func db.GetOFHome err:%v", userId, sceneId, err)
				return nil
			}
			data = &SceneGardenData{
				SceneId:                     home.SceneID,
				GardenName:                  home.GardenName,
				LikesNum:                    home.LikesNum,
				AccessPlayerNum:             home.AccessPlayerNum,
				LeftLikeNum:                 home.LeftLikeNum,
				IsOpen:                      home.IsOpen,
				Password:                    home.Password,
				GardenFurnitureInfoMap:      make(map[int64]*FurnitureDetailsInfo),
				OtherPlayerFurnitureInfoMap: make(map[uint32]*FurnitureDetailsInfo),
			}
			sceneGardenCache.Set(fmt.Sprintf("%v|%v", userId, sceneId), data)
		}
	}
	return data
}

type FurnitureDetailsInfo struct {
	FurnitureId     int64    `json:"furnitureId;omitempty"`     // 家具id
	FurnitureItemId uint32   `json:"furnitureItemId;omitempty"` // 家具物品id
	Pos             *Vector3 `json:"pos;omitempty"`             // 坐标
	Rotation        *Vector3 `json:"rotation;omitempty"`        // 坐标
	LayerNum        uint32   `json:"layerNum;omitempty"`
}

type Vector3 struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	Z int32 `json:"z"`
}

func (v *Vector3) Vector3() *proto.Vector3 {
	return &proto.Vector3{X: v.X, Y: v.Y, Z: v.Z}
}

func ToVector3(v *proto.Vector3) *Vector3 {
	return &Vector3{X: v.X, Y: v.Y, Z: v.Z}
}

func (f *FurnitureDetailsInfo) FurnitureDetailsInfo() *proto.FurnitureDetailsInfo {
	return &proto.FurnitureDetailsInfo{
		FurnitureId:     f.FurnitureId,
		FurnitureItemId: f.FurnitureItemId,
		Pos:             f.Pos.Vector3(),
		Rotation:        f.Rotation.Vector3(),
		LayerNum:        f.LayerNum,
	}
}

func (s *SceneGardenData) AddFurniture(userId, channelId uint32, furniture *proto.FurnitureDetailsInfo) {
	info := &FurnitureDetailsInfo{
		FurnitureId:     furniture.FurnitureId,
		FurnitureItemId: furniture.FurnitureItemId,
		Pos:             ToVector3(furniture.Pos),
		Rotation:        ToVector3(furniture.Rotation),
		LayerNum:        furniture.LayerNum,
	}
	if userId != channelId {
		// 客人
		s.OtherPlayerFurnitureInfoMap[userId] = info
	} else {
		// 主人
		s.GardenFurnitureInfoMap[furniture.FurnitureId] = info
	}
}

// 删除家具
func (s *SceneGardenData) RemoveFurniture(userId, channelId uint32, furnitureId int64) *FurnitureDetailsInfo {
	if userId != channelId {
		info, ok := s.OtherPlayerFurnitureInfoMap[userId]
		if ok {
			delete(s.OtherPlayerFurnitureInfoMap, userId)
		}
		return info
	} else {
		info, ok := s.GardenFurnitureInfoMap[furnitureId]
		if ok {
			delete(s.GardenFurnitureInfoMap, furnitureId)
		}
		return info
	}
}

func (s *SceneGardenData) SceneGardenData() *proto.SceneGardenData {
	info := &proto.SceneGardenData{
		GardenFurnitureInfoMap:        make(map[int64]*proto.FurnitureDetailsInfo), // 主人家具
		LikesNum:                      s.LikesNum,
		AccessPlayerNum:               s.AccessPlayerNum,
		LeftLikeNum:                   s.LeftLikeNum,
		GardenName:                    s.GardenName,
		FurniturePlayerMap:            make(map[int64]uint32),
		OtherPlayerFurnitureInfoMap:   make(map[int64]*proto.SceneGardenOtherPlayerData), // 客人家具
		FurnitureCurrentPointNum:      0,
		PlayerHandingFurnitureInfoMap: make(map[int64]*proto.SceneGardenOtherPlayerData),
	}
	for _, v := range s.GardenFurnitureInfoMap {
		info.GardenFurnitureInfoMap[v.FurnitureId] = v.FurnitureDetailsInfo()
	}
	for k, v := range s.OtherPlayerFurnitureInfoMap {
		info.OtherPlayerFurnitureInfoMap[v.FurnitureId] = &proto.SceneGardenOtherPlayerData{
			PlayerId:      k,
			FurnitureInfo: v.FurnitureDetailsInfo(),
		}
	}
	return info
}

func (s *SceneGardenData) GardenBaseInfo() *proto.GardenBaseInfo {
	info := &proto.GardenBaseInfo{
		LikesNum:           s.LikesNum,
		AccessNum:          s.AccessPlayerNum,
		FurnitureNum:       0, // 家具数量
		FurnitureLimitNum:  gdconf.GetConstant().FurnitureLimitNum,
		IsOpen:             s.IsOpen,
		Password:           s.Password,
		PasswordExpireTime: 0,
	}
	return info
}
