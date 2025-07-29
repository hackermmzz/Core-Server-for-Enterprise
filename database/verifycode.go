package database

import "time"

// 用途
const (
	VerifyCode_Regist = 0
)

type VerifyCode struct {
	Email        string     `json:"email" gorm:"column:email;not null;type:varchar(255)"`
	Code         string     `json:"code" gorm:"column:code;not null;type:varhcar(10)"`
	GenerateDate *LocalTime `json:"generateDate" gorm:"column:generateDate;not null;type:timestamp;default:CURRENT_TIMESTAMP"`
	Usage        int        `json:"usage" gorm:"column:usage_;not null;type:int"`
}

func (VerifyCode) TableName() string {
	return "verifycode"
}

// 获取所有email和usage均匹配的未过期的验证码(五分钟)
func GetTargetVerifyCode(email string, usage int) ([]VerifyCode, error) {
	var ret []VerifyCode
	err := db.Select("code").Where("email=? and usage_=? and generateDate>?", email, usage, time.Now().UTC().Add(time.Duration(-5*time.Minute))).Find(&ret).Error
	return ret, err
}

// 添加验证码
func AddVerifyCode(verifyCode *VerifyCode) error {
	return db.Table(verifyCode.TableName()).Create(verifyCode).Error
}

// 获取最早得指定验证码记录
func GetEarliestVerifyCode(email string, usage int) (VerifyCode, error) {
	var ret VerifyCode
	err := db.Where("email=? and usage_=?", email, usage).Order("generateDate DESC").First(&ret).Error
	return ret, err
}
