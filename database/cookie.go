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
	return db.Where("generateDate<?", time.Now().UTC().AddDate(0, 0, -1)).Delete(&Cookie{}).Error
}

// 查询cookie是否在数据库出现且没有过期
func CheckCookieExist(cookie string) bool {
	var cookie_ Cookie
	err := db.Where("cookieValue=? and generateDate>?", cookie, time.Now().UTC().AddDate(0, 0, -1)).First(&cookie_).Error
	return err == nil
}
