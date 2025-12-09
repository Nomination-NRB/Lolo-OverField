package game

type ChatInfo struct {
	noticeChan    *ChatChannel            // 通知频道
	allSystemChat map[uint32]*ChatChannel // 系统频道
	privateChat   map[uint64]*ChatChannel // 私聊频道
}

func (g *Game) getChatInfo() *ChatInfo {
	if g.chatInfo == nil {
		chatInfo := &ChatInfo{
			noticeChan:    nil,
			allSystemChat: make(map[uint32]*ChatChannel),
		}
		g.chatInfo = chatInfo
	}
	return g.chatInfo
}

// 获取通知频道
func (c *ChatInfo) getNoticeChan() *ChatChannel {
	return c.noticeChan
}

// ChatChannel 聊天房间对象
type ChatChannel struct {
}

func newChatChannel() *ChatChannel {
	return &ChatChannel{}
}
