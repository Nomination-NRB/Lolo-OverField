package game

import (
	"gucooing/lolo/game/model"
	"gucooing/lolo/pkg/alg"
	"gucooing/lolo/pkg/log"
	"gucooing/lolo/protocol/cmd"
	"gucooing/lolo/protocol/proto"
)

func (g *Game) GetAllCharacterEquip(s *model.Player, msg *alg.GameMsg) {
	rsp := &proto.GetAllCharacterEquipRsp{
		Status: proto.StatusCode_StatusCode_OK,
		Items:  make([]*proto.ItemDetail, 0),
	}
	defer g.send(s, cmd.GetAllCharacterEquipRsp, msg.PacketId, rsp)
}

func (g *Game) GetCharacterAchievementList(s *model.Player, msg *alg.GameMsg) {
	req := msg.Body.(*proto.GetCharacterAchievementListReq)
	rsp := &proto.GetCharacterAchievementListRsp{
		Status:                  proto.StatusCode_StatusCode_OK,
		CharacterAchievementLst: make([]*proto.Achieve, 0),
		HasRewardedIds:          make([]uint32, 0),
		IsUnlockedPayment:       false,
		CharacterId:             req.CharacterId,
		RewardedIdLst:           make([]uint32, 0),
	}
	defer g.send(s, cmd.GetCharacterAchievementListRsp, msg.PacketId, rsp)
}

func (g *Game) OutfitPresetUpdate(s *model.Player, msg *alg.GameMsg) {
	req := msg.Body.(*proto.OutfitPresetUpdateReq)
	rsp := &proto.OutfitPresetUpdateRsp{
		Status: proto.StatusCode_StatusCode_OK,
		CharId: req.CharId,
		Preset: req.Preset,
	}
	defer func() {
		g.send(s, cmd.OutfitPresetUpdateRsp, msg.PacketId, rsp)
		teamInfo := s.GetTeamModel().GetTeamInfo()
		scenePlayer := g.getWordInfo().getScenePlayer(s)
		if (req.CharId == teamInfo.Char1 ||
			req.CharId == teamInfo.Char2 ||
			req.CharId == teamInfo.Char3) &&
			(scenePlayer != nil &&
				scenePlayer.channelInfo != nil) {
			scenePlayer.channelInfo.serverSceneSyncChan <- &ServerSceneSyncCtx{
				ScenePlayer: scenePlayer,
				ActionType:  proto.SceneActionType_SceneActionType_UPDATE_FASHION,
			}
		}
	}()
	characterInfo := s.GetCharacterModel().GetCharacterInfo(req.CharId)
	if characterInfo == nil {
		log.Game.Warnf("保存角色预设装扮失败,角色%v不存在", req.CharId)
		return
	}
	outfitPreset := characterInfo.GetOutfitPreset(req.Preset.PresetIndex)

	outfitPreset.Hat = req.Preset.Hat
	outfitPreset.HatDyeSchemeIndex = req.Preset.HatDyeSchemeIndex
	outfitPreset.Hair = req.Preset.Hair
	outfitPreset.HairDyeSchemeIndex = req.Preset.HairDyeSchemeIndex
	outfitPreset.Clothes = req.Preset.Clothes
	outfitPreset.ClothesDyeSchemeIndex = req.Preset.ClothesDyeSchemeIndex
	outfitPreset.Ornament = req.Preset.Ornament
	outfitPreset.OrnamentDyeSchemeIndex = req.Preset.OrnamentDyeSchemeIndex
	outfitPreset.OutfitHideInfo = &model.OutfitHideInfo{
		HideOrn:   req.Preset.OutfitHideInfo.HideOrn,
		HideBraid: req.Preset.OutfitHideInfo.HideBraid,
	}
	outfitPreset.PendTop = req.Preset.PendTop
	outfitPreset.PendTopDyeSchemeIndex = req.Preset.PendTopDyeSchemeIndex
	outfitPreset.PendChest = req.Preset.PendChest
	outfitPreset.PendChestDyeSchemeIndex = req.Preset.PendChestDyeSchemeIndex
	outfitPreset.PendPelvis = req.Preset.PendPelvis
	outfitPreset.PendPelvisDyeSchemeIndex = req.Preset.PendPelvisDyeSchemeIndex
	outfitPreset.PendUpFace = req.Preset.PendUpFace
	outfitPreset.PendUpFaceDyeSchemeIndex = req.Preset.PendUpFaceDyeSchemeIndex
	outfitPreset.PendDownFace = req.Preset.PendDownFace
	outfitPreset.PendDownFaceDyeSchemeIndex = req.Preset.PendDownFaceDyeSchemeIndex
	outfitPreset.PendLeftHand = req.Preset.PendLeftHand
	outfitPreset.PendLeftHandDyeSchemeIndex = req.Preset.PendLeftHandDyeSchemeIndex
	outfitPreset.PendRightHand = req.Preset.PendRightHand
	outfitPreset.PendRightHandDyeSchemeIndex = req.Preset.PendRightHandDyeSchemeIndex
	outfitPreset.PendLeftFoot = req.Preset.PendLeftFoot
	outfitPreset.PendLeftFootDyeSchemeIndex = req.Preset.PendLeftFootDyeSchemeIndex
	outfitPreset.PendRightFoot = req.Preset.PendRightFoot
	outfitPreset.PendRightFootDyeSchemeIndex = req.Preset.PendRightFootDyeSchemeIndex
}

func (g *Game) CharacterEquipUpdate(s *model.Player, msg *alg.GameMsg) {
	req := msg.Body.(*proto.CharacterEquipUpdateReq)
	rsp := &proto.CharacterEquipUpdateRsp{
		Status:    proto.StatusCode_StatusCode_OK,
		Character: make([]*proto.Character, 0),
		Items:     make([]*proto.ItemDetail, 0),
	}
	defer g.send(s, cmd.CharacterEquipUpdateRsp, msg.PacketId, rsp)
	characterInfo := s.GetCharacterModel().GetCharacterInfo(req.CharId)
	if characterInfo == nil {
		log.Game.Warnf("保存角色装备失败,角色%v不存在", req.CharId)
		return
	}
	defer alg.AddList(&rsp.Character, characterInfo.Character())

	equipmentPreset := characterInfo.GetEquipmentPreset(req.EquipmentPreset.PresetIndex)
	// 更新武器
	if req.EquipmentPreset.Weapon != equipmentPreset.Weapon {
		oldEquipmentInfo := s.GetItemModel().GetItemWeaponInfo(equipmentPreset.Weapon)
		newEquipmentInfo := s.GetItemModel().GetItemWeaponInfo(req.EquipmentPreset.Weapon)
		if newEquipmentInfo != nil &&
			oldEquipmentInfo != nil {
			oldEquipmentInfo.WearerId = 0
			alg.AddList(&rsp.Items, oldEquipmentInfo.ItemDetail())

			if oldCharacterInfo := s.GetCharacterModel().GetCharacterInfo(newEquipmentInfo.WearerId); oldCharacterInfo != nil {
				// 移除装备上的角色
			}
			newEquipmentInfo.WearerId = req.CharId
			equipmentPreset.Weapon = newEquipmentInfo.InstanceId
			alg.AddList(&rsp.Items, newEquipmentInfo.ItemDetail())
		}
	}
	// 更新盔甲
	// 更新海报
}
