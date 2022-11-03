package userschema

import (
	"fmt"
	"gorm.io/gorm"
)

type Following struct {
	UserID   uint
	Followed uint
}

// after create hook
func (f *Following) AfterCreate(tx *gorm.DB) error {
	err := tx.Model(&User{}).Where("id = ?", f.UserID).Updates(
		map[string]interface{}{
			"following_count": gorm.Expr("following_count + ?", 1),
		}).Error
	if err != nil {
		return err
	}
	err = tx.Model(&User{}).Where("id = ?", f.Followed).Updates(
		map[string]interface{}{
			"follower_count": gorm.Expr("follower_count + ?", 1),
		}).Error
	return err
}

// after delete hook
func (f *Following) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("Halo hook, czemu nie wchodzisz")
	err = tx.Model(&User{}).Where("id = ?", f.UserID).Updates(
		map[string]interface{}{
			"following_count": gorm.Expr("following_count - ?", 1),
		}).Error
	if err != nil {
		return err
	}
	err = tx.Model(&User{}).Where("id = ?", f.Followed).Updates(
		map[string]interface{}{
			"follower_count": gorm.Expr("follower_count - ?", 1),
		}).Error
	return err

}
