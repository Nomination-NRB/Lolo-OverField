package model

import (
	"gucooing/lolo/gdconf"
	"gucooing/lolo/protocol/proto"
)

type ShopModel struct {
	ShopInfos map[uint32]*ShopInfo `json:"ShopInfos,omitempty"` // 商店信息

}

func NewShopModel() *ShopModel {
	return &ShopModel{}
}

func (s *Player) GetShopModel() *ShopModel {
	if s.Shop == nil {
		s.Shop = NewShopModel()
	}
	return s.Shop
}

type ShopInfo struct {
	ShopID    uint32               `json:"shopId"`
	ShopGrids map[uint32]*GridInfo `json:"shopGridMap,omitempty"` // 格子信息
}

func (s *ShopModel) GetShopInfos() map[uint32]*ShopInfo {
	if s.ShopInfos == nil {
		s.ShopInfos = make(map[uint32]*ShopInfo)
	}
	return s.ShopInfos
}

func (s *ShopModel) GetShopInfo(shopId uint32) *ShopInfo {
	ls := s.GetShopInfos()
	info, ok := ls[shopId]
	if !ok {
		info = &ShopInfo{
			ShopID:    shopId,
			ShopGrids: make(map[uint32]*GridInfo),
		}
		ls[shopId] = info
	}
	return info
}

func (s *ShopInfo) GetGridInfos() map[uint32]*GridInfo {
	if s.ShopGrids == nil {
		s.ShopGrids = make(map[uint32]*GridInfo)
	}
	return s.ShopGrids
}

func (s *ShopInfo) GetGridInfo(gridId uint32) *GridInfo {
	ls := s.GetGridInfos()
	info, ok := ls[gridId]
	if !ok {
		conf := gdconf.GetGrid(s.ShopID, gridId)
		if conf == nil {
			return nil
		}
		info = &GridInfo{
			ShopID: s.ShopID,
			GridID: gridId,
			PoolID: uint32(conf.ShopPoolID),
		}
		ls[gridId] = info
	}
	return info
}

type GridInfo struct {
	ShopID   uint32 `json:"shopId"`
	GridID   uint32 `json:"gridId"`
	PoolID   uint32 `json:"poolId"`
	BuyTimes uint32 `json:"buyTimes"`
}

func (g *GridInfo) ShopGrid() *proto.ShopGrid {
	info := &proto.ShopGrid{
		Id:         g.ShopID,
		GridId:     g.GridID,
		PoolId:     g.PoolID,
		PoolIndex:  1, // TODO PoolIndex 需要根据配置表来
		BuyTimes:   g.BuyTimes,
		UpdateTime: 0,
	}

	return info
}
