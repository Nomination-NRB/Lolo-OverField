package db

// 家园数据
type OFHome struct {
	UserID          uint32 `gorm:"index;index:idx_id"`
	SceneID         uint32 `gorm:"index;index:idx_id"`
	GardenName      string
	LikesNum        int64
	AccessPlayerNum int64
	LeftLikeNum     uint32
	IsOpen          bool
	Password        string
}

func GetOFHome(userId, sceneId uint32) (*OFHome, error) {
	h := &OFHome{
		UserID:  userId,
		SceneID: sceneId,
	}
	err := db.Where("user_id = ? AND scene_id = ?", userId, sceneId).
		FirstOrCreate(&h).Error
	return h, err
}
