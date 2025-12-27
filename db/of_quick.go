package db

// 快游戏sdk
type OFQuick struct {
	ID        uint32 `gorm:"primarykey;autoIncrement;index"`
	Username  string `gorm:"unique;index"`
	Password  string `gorm:"not null"`
	RegDevice string `gorm:"not null"`
	UserToken string // 网关登录
	AuthToken string `gorm:"unique;not null"` // 自动登录
}

func CreateOFQuick(username string, password string) (*OFQuick, error) {
	q := &OFQuick{
		Username: username,
		Password: password,
	}
	err := db.Create(q).Error
	return q, err
}

func GetOFQuick(username, password string) (*OFQuick, error) {
	q := &OFQuick{
		Username: username,
		Password: password,
	}
	err := db.Where("username = ?", username).FirstOrCreate(q).Error
	if err != nil {
		return nil, err
	}
	return q, err
}

func GetOFQuickByAuthToken(authToken string) (*OFQuick, error) {
	q := new(OFQuick)
	err := db.Where("auth_token = ?", authToken).First(q).Error
	return q, err
}

func UpOFQuick(q *OFQuick) error {
	return db.Save(q).Error
}
