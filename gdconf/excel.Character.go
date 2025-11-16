package gdconf

import (
	"gucooing/lolo/protocol/excel"
)

type Character struct {
	all             *excel.AllCharacterDatas
	CharacterAllMap map[uint32]*CharacterAllInfo
	GrowthLevelMap  map[int32]map[uint32]*excel.CharacterLevelInfo
}

type CharacterAllInfo struct {
	CharacterId   uint32
	CharacterInfo *excel.CharacterConfigure
	LevelRules    []*excel.CharacterLevelRuleInfo
}

func (g *GameConfig) loadCharacter() {
	info := &Character{
		all:             new(excel.AllCharacterDatas),
		CharacterAllMap: make(map[uint32]*CharacterAllInfo),
		GrowthLevelMap:  make(map[int32]map[uint32]*excel.CharacterLevelInfo),
	}
	g.Excel.Character = info
	name := "Character.json"
	ReadJson(g.excelPath, name, &info.all)

	getCharacterAllMap := func(id int32) *CharacterAllInfo {
		if info.CharacterAllMap[uint32(id)] == nil {
			info.CharacterAllMap[uint32(id)] = &CharacterAllInfo{
				CharacterId: uint32(id),
			}
		}
		return info.CharacterAllMap[uint32(id)]
	}
	getLevelMap := func(id int32) map[uint32]*excel.CharacterLevelInfo {
		if info.GrowthLevelMap[id] == nil {
			info.GrowthLevelMap[id] = make(map[uint32]*excel.CharacterLevelInfo)
		}
		return info.GrowthLevelMap[id]
	}

	for _, v := range info.all.GetCharacter().GetDatas() {
		getCharacterAllMap(v.ID).CharacterInfo = v
	}
	for _, v := range info.all.GetLevelRule().GetDatas() {
		getCharacterAllMap(v.ID).LevelRules = v.LevelRuleInfo
	}
	for _, v := range info.all.GetLevel().GetDatas() {
		levelMap := getLevelMap(v.ID)
		for _, v2 := range v.GetLevelInfo() {
			levelMap[uint32(v2.Level)] = v2
		}
	}
}

func GetCharacterAllMap() map[uint32]*CharacterAllInfo {
	return cc.Excel.Character.CharacterAllMap
}

func GetCharacterAll(id uint32) *CharacterAllInfo {
	return cc.Excel.Character.CharacterAllMap[id]
}

func GetCharacterLevelMap(id int32) map[uint32]*excel.CharacterLevelInfo {
	return cc.Excel.Character.GrowthLevelMap[id]
}

func AddCharacterExp(levelId, oldExp int32, oldLevel, maxLevel uint32) (newLevel, newExp uint32) {
	levelMap := GetCharacterLevelMap(levelId)
	for {
		if oldLevel >= maxLevel {
			return oldLevel, 0
		}
		conf, ok := levelMap[oldLevel]
		if !ok || oldExp < conf.NeedExp {
			return oldLevel, uint32(oldExp)
		}
		oldExp -= conf.NeedExp
		oldLevel++
	}
}
