package database

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// 定义时间结构体
type LocalTime time.Time

// MarshalJSON JSON序列化时调用，转换为指定格式的字符串
func (t *LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	// 格式化为 "2006-01-02 15:04:05" 字符串
	return []byte(fmt.Sprintf("\"%s\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

// Value 数据库存储时调用，转换为time.Time类型
func (t LocalTime) Value() (driver.Value, error) {
	tTime := time.Time(t)
	// 零值处理：如果是默认零时间，返回nil（数据库会存为NULL）
	if tTime.IsZero() {
		return nil, nil
	}
	return tTime, nil
}

// Scan 数据库读取时调用，将数据库返回值转换为LocalTime
func (t *LocalTime) Scan(v interface{}) error {
	switch value := v.(type) {
	case time.Time:
		// 直接是time.Time类型，直接转换
		*t = LocalTime(value)
		return nil
	case []byte:
		// 数据库返回字节数组（最常见情况），转换为字符串后解析
		timeStr := string(value)
		// 解析字符串为time.Time（格式要与数据库存储的一致）
		parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return fmt.Errorf("解析时间失败: %v, 原始值: %s", err, timeStr)
		}
		*t = LocalTime(parsedTime)
		return nil
	case string:
		// 少数情况可能返回字符串，同样解析
		parsedTime, err := time.Parse("2006-01-02 15:04:05", value)
		if err != nil {
			return fmt.Errorf("解析时间失败: %v, 原始值: %s", err, value)
		}
		*t = LocalTime(parsedTime)
		return nil
	default:
		// 不支持的类型
		return fmt.Errorf("不支持的类型 %T，无法转换为LocalTime", v)
	}
}

// 用户
type User struct {
	Email      string     `json:"email" gorm:"column:email;type:varchar(255);not null;primary key"`
	NickName   string     `json:"nickname" gorm:"column:nickname;type:varchar(30);not null;"`
	Password   string     `json:"password" gorm:"column:password;type:varchar(255);not null"`
	RegistDate *LocalTime `json:"registDate" gorm:"column:registDate;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	Avatar     string     `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;default:Default/avatar.png"`
	Vip        int        `json:"vip" gorm:"column:vip;type:int;not null;default:0"`
}

func (user *User) TableName() string {
	return "users"
}

// 向用户表添加用户
func AddUser(user User) error {
	return db.Table(user.TableName()).Create(&user).Error
}

// 获取指定用户
func GetUser(email string) (User, error) {
	var user User
	res := db.Where("email=?", email).First(&user)
	return user, res.Error
}
