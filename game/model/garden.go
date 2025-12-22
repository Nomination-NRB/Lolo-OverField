package model

import (
	"gucooing/lolo/gdconf"
	"gucooing/lolo/pkg/alg"
	"gucooing/lolo/protocol/proto"
)

type GardenModel struct {
	SchemeMap map[uint32]*GardenSchemeInfo `json:"schemeMap,omitempty"`
}

type GardenSchemeInfo struct {
	SchemeName       string                          `json:"schemeName,omitempty"`
	FurnitureInfoMap map[int64]*FurnitureDetailsInfo `json:"furnitureInfoMap,omitempty"` // 家具信息
}

func (s *Player) GetGardenModel() *GardenModel {
	if s.Garden == nil {
		s.Garden = new(GardenModel)
	}
	return s.Garden
}

func (s *GardenModel) GetSchemeMap() map[uint32]*GardenSchemeInfo {
	if s.SchemeMap == nil {
		s.SchemeMap = make(map[uint32]*GardenSchemeInfo)
	}
	return s.SchemeMap
}

func (s *GardenModel) Schemes() map[uint32]string {
	schemes := make(map[uint32]string, gdconf.GetConstant().GardenSchemeNum)
	for index := range gdconf.GetConstant().GardenSchemeNum {
		info := s.GetGardenSchemeInfo(index)
		schemes[index] = info.SchemeName
	}
	return schemes
}

func (s *GardenModel) GetGardenSchemeInfo(schemeId uint32) *GardenSchemeInfo {
	list := s.GetSchemeMap()
	info, ok := list[schemeId]
	if !ok {
		info = &GardenSchemeInfo{
			SchemeName: "",
		}
		list[schemeId] = info
	}
	return info
}

func (s *GardenSchemeInfo) FurnitureItemId() []uint32 {
	furnitureItemId := make([]uint32, len(s.FurnitureInfoMap))
	for _, v := range s.FurnitureInfoMap {
		alg.AddLists(&furnitureItemId, v.FurnitureItemId)
	}
	return furnitureItemId
}

func (s *GardenSchemeInfo) FurnitureDetailsInfos() []*proto.FurnitureDetailsInfo {
	list := make([]*proto.FurnitureDetailsInfo, len(s.FurnitureInfoMap))
	for _, v := range s.FurnitureInfoMap {
		alg.AddList(&list, v.FurnitureDetailsInfo())
	}
	return list
}
