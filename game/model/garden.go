package model

import (
	"fmt"

	"gucooing/lolo/gdconf"
)

var (
	gardenModelKey = func(index uint32) string {
		return fmt.Sprintf("预设%d", index+1)
	}
)

type GardenModel struct {
	SchemeMap map[uint32]*GardenSchemeInfo
}

type GardenSchemeInfo struct {
	SchemeName string
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
	list := s.GetSchemeMap()
	for index := range gdconf.GetConstant().GardenSchemeNum {
		info, ok := list[index]
		if !ok {
			info = &GardenSchemeInfo{
				SchemeName: gardenModelKey(index),
			}
			list[index] = info
		}
		schemes[index] = info.SchemeName
	}
	return schemes
}
