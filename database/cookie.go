package database

import "time"

// cookie表
type Cookie struct {
	Email        string     `json:"email" gorm:"column:email;not null"`
	CookieValue  string     `json:"cookieValue" gorm:"column:cookieValue;not null"`
	GenerateDate *LocalTime `json:"generateDate" gorm:"column:generateDate;not null;default:CURRENT_TIMESTAMP"`
}

func (ck Cookie) TableName() string {
	return "cookie"
}

// 将cookie加入数据库
func AddUserCookie(cookie Cookie) error {
	return db.Table(cookie.TableName()).Create(cookie).Error
}

// 移除所有过期的cookie
func RemoveExpireCookie() error {
	return db.Where("generateDate<?", time.Now().AddDate(0, 0, -1)).Delete(&Cookie{}).Error
}
