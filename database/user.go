package database

type User struct {
	Email      string `json:"email";gorm:"type:varchar(255);not null;primary key"`
	Password   string `json:"password";gorm:"type:varchar(255);not null"`
	RegistDate string `json:"registDate";gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	Avatar     string `json:"avatar";gorm:"type:varchar(255)"`
	Vip        int    `json:"vip";gorm:"type:int;not null;default:0"`
}

func (user *User) TableName() string {
	return "users"
}

// 向用户表添加用户
func AddUser(email string, password string) error {
	var user User
	user.Email = email
	user.Password = password
	return db.Table("users").Create(&user).Error
}

// 获取指定用户
func GetUser(email string) (User, error) {
	var user User
	res := db.Where("email=?", email).First(&user)
	return user, res.Error
}
