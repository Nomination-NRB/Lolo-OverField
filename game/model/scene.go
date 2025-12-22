package model

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"

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
	PrivateChannelStart = 1000000
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
	DBID                        uint64                           `json:"-"`                         // 数据库主键
	SceneId                     uint32                           `json:"sceneId,omitempty"`         // 场景id
	UserId                      uint32                           `json:"userId,omitempty"`          // 玩家id
	GardenName                  string                           `json:"gardenName,omitempty"`      // 花园名称
	LikesNum                    int64                            `json:"likesNum,omitempty"`        // 点赞数
	AccessPlayerNum             int64                            `json:"accessPlayerNum,omitempty"` // 访问数
	LeftLikeNum                 uint32                           `json:"leftLikeNum,omitempty"`
	IsOpen                      bool                             `json:"isOpen,omitempty"`                 // 是否开放
	Password                    string                           `json:"password"`                         // 密码
	PasswordExpireTime          int64                            `json:"passwordExpireTime,omitempty"`     // 下次改密时间
	GardenFurnitureInfoMap      map[int64]*FurnitureDetailsInfo  `json:"gardenFurnitureInfoMap,omitempty"` // 主人家具信息
	OtherPlayerFurnitureInfoMap map[uint32]*FurnitureDetailsInfo `json:"-"`                                // 客人家具信息
}

func GetSceneGardenData(userId, sceneId uint32) *SceneGardenData {
	if sceneId != 9999 {
		return &SceneGardenData{
			SceneId:                     sceneId,
			GardenFurnitureInfoMap:      make(map[int64]*FurnitureDetailsInfo),
			OtherPlayerFurnitureInfoMap: make(map[uint32]*FurnitureDetailsInfo),
		}
	}
	data, ok := sceneGardenCache.Get(fmt.Sprintf("%v|%v", userId, sceneId))
	if userId < PrivateChannelStart {
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
				DBID:                        home.ID,
				SceneId:                     home.SceneID,
				UserId:                      home.UserID,
				GardenName:                  home.GardenName,
				LikesNum:                    home.LikesNum,
				AccessPlayerNum:             home.AccessPlayerNum,
				LeftLikeNum:                 home.LeftLikeNum,
				IsOpen:                      home.IsOpen,
				Password:                    home.Password,
				PasswordExpireTime:          home.PasswordExpireTime,
				GardenFurnitureInfoMap:      make(map[int64]*FurnitureDetailsInfo),
				OtherPlayerFurnitureInfoMap: make(map[uint32]*FurnitureDetailsInfo),
			}
			if len(home.GardenFurnitureInfoMap) != 0 {
				if err := sonic.Unmarshal(home.GardenFurnitureInfoMap, &data.GardenFurnitureInfoMap); err != nil {
					log.Game.Errorf("UserId:%v SceneId:%v GardenFurnitureInfoMapUnmarshal err:%v",
						userId, sceneId, err)
				}
			}
			sceneGardenCache.Set(fmt.Sprintf("%v|%v", userId, sceneId), data)
		}
	}
	return data
}

type FurnitureDetailsInfo struct {
	FurnitureId     int64    `json:"furnitureId,omitempty"`     // 家具id
	FurnitureItemId uint32   `json:"furnitureItemId,omitempty"` // 家具物品id
	Pos             *Vector3 `json:"pos,omitempty"`             // 坐标
	Rotation        *Vector3 `json:"rotation,omitempty"`        // 坐标
	LayerNum        uint32   `json:"layerNum,omitempty"`
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

func (s *SceneGardenData) AddFurniture(user *Player, channelId uint32,
	furniture *proto.FurnitureDetailsInfo, save bool) {
	info := &FurnitureDetailsInfo{
		FurnitureId:     furniture.FurnitureId,
		FurnitureItemId: furniture.FurnitureItemId,
		Pos:             ToVector3(furniture.Pos),
		Rotation:        ToVector3(furniture.Rotation),
		LayerNum:        furniture.LayerNum,
	}
	if user.UserId != channelId {
		// 客人
		s.OtherPlayerFurnitureInfoMap[user.UserId] = info
	} else {
		// 主人
		user.GetItemModel().AddFurnitureItem(info.FurnitureItemId)
		s.GardenFurnitureInfoMap[furniture.FurnitureId] = info
		if save {
			s.Save()
		}
	}
}

// 删除家具
func (s *SceneGardenData) RemoveFurniture(user *Player, channelId uint32, furnitureId int64, save bool) *FurnitureDetailsInfo {
	if user.UserId != channelId {
		info, ok := s.OtherPlayerFurnitureInfoMap[user.UserId]
		if ok {
			delete(s.OtherPlayerFurnitureInfoMap, user.UserId)
		}
		return info
	} else {
		info, ok := s.GardenFurnitureInfoMap[furnitureId]
		if ok {
			user.GetItemModel().DelFurnitureItem(info.FurnitureItemId)
			delete(s.GardenFurnitureInfoMap, furnitureId)
			save = save == true
		}
		if save {
			s.Save()
		}
		return info
	}
}

func (s *SceneGardenData) Save() {
	h := &db.OFHome{
		ID:                     s.DBID,
		UserID:                 s.UserId,
		SceneID:                s.SceneId,
		GardenName:             s.GardenName,
		LikesNum:               s.LikesNum,
		AccessPlayerNum:        s.AccessPlayerNum,
		LeftLikeNum:            s.LeftLikeNum,
		IsOpen:                 s.IsOpen,
		Password:               s.Password,
		PasswordExpireTime:     s.PasswordExpireTime,
		GardenFurnitureInfoMap: nil,
	}
	bin, err := sonic.Marshal(s.GardenFurnitureInfoMap)
	if err != nil {
		log.Game.Errorf("玩家花园数据:%v序列化失败err:%s",
			s.UserId, err.Error())
		return
	}
	h.GardenFurnitureInfoMap = bin
	if err := db.SaveOFHome(h); err != nil {
		log.Game.Errorf("UserId:%v SceneId:%v db.SaveOFHome err:%s", s.UserId, s.SceneId, err.Error())
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

func (s *SceneGardenData) NewFurnitureList() []*proto.FurnitureDetailsInfo {
	list := make([]*proto.FurnitureDetailsInfo, 0)
	for _, v := range s.GardenFurnitureInfoMap {
		alg.AddList(&list, v.FurnitureDetailsInfo())
	}
	return list
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
