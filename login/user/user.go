package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Ctn            string    `json:"ctn"`
	LastLoginDate  time.Time `gorm:"default:null"`
	PasswordUpdate time.Time `gorm:"default:null"`
	FailCount      int
	Status         int
	Password       string `json:"password"`
	OldPassword    string
	Email          string `json:"email"`
	AuthId         string `json:"auth_id"`
}

const (
	StandBy = iota
	Active
	Stop
)

func (u *User) TableName() string {
	return "user"
}

func (u *User) FindDuplicateDateUser(ctn string, db *gorm.DB) int64 {
	result := db.Find(&u, "ctn = ?", ctn)
	return result.RowsAffected
}

func (u *User) AddFailCount(db *gorm.DB) (int, error) {
	addCount := u.FailCount + 1
	err := db.Model(&u).Update("fail_count", addCount).Error
	if err != nil {
		return addCount, err
	}
	return addCount, nil
}

func ModifyStatus(id, newStatus string, db *gorm.DB) error {
	err := db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": newStatus,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ModifyPassword(id, newPassword string, db *gorm.DB) error {
	result := db.Find(&u, "id = ?", id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	tmpPassword := u.Password

	err := db.Model(&u).Where("id = ?", id).Updates(map[string]interface{}{
		"password":     newPassword,
		"old_password": tmpPassword,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
