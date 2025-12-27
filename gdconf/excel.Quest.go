package gdconf

import (
	"gucooing/lolo/protocol/excel"
)

type Quest struct {
	all        *excel.AllQuestDatas
	QuestInfos map[uint32]*QuestInfo
}

type QuestInfo struct {
	Config            *excel.QuestConfigure
	ConditionSetGroup *excel.ConditionSetGroupConfigure
}

func (g *GameConfig) loadQuest() {
	info := &Quest{
		all:        new(excel.AllQuestDatas),
		QuestInfos: make(map[uint32]*QuestInfo),
	}
	g.Excel.Quest = info
	name := "Quest.json"
	ReadJson(g.excelPath, name, &info.all)

	for _, v := range info.all.GetQuest().GetDatas() {
		questInfo := &QuestInfo{
			Config:            v,
			ConditionSetGroup: nil,
		}
		info.QuestInfos[uint32(v.ID)] = questInfo
		for _, v2 := range info.all.GetConditionSetGroup().GetDatas() {
			if v.ConditionSetGroupID == v2.ID {
				questInfo.ConditionSetGroup = v2
				break
			}
		}
	}
}

func GetQuestInfos() map[uint32]*QuestInfo {
	return cc.Excel.Quest.QuestInfos
}
