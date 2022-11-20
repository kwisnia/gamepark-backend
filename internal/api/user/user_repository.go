package user

import (
	"fmt"
	"github.com/kwisnia/inzynierka-backend/internal/api/user/userschema"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config/database"
	"gorm.io/gorm"
	"strings"
)

func SaveNewUser(u *userschema.User) {
	database.DB.Create(u)
}

func UpdateUser(u *userschema.User) error {
	return database.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(u).Error
}

func GetByUsername(username string) *userschema.User {
	var u userschema.User
	if err := database.DB.Preload("UserProfile").Preload("UserFeatureUnlock").Preload("Lists").Where("LOWER(username) = ?", strings.ToLower(username)).First(&u).Error; err != nil {
		return nil
	}
	return &u
}

func GetByID(id uint) *userschema.User {
	var u userschema.User
	if err := database.DB.Preload("UserProfile").Preload("UserFeatureUnlock").Preload("Lists").Where("id = ?", id).First(&u).Error; err != nil {
		return nil
	}
	return &u
}

func GetBySearch(pageSize int, offset int, search string) ([]userschema.User, error) {
	var users []userschema.User
	query := GetPageQuery(pageSize, offset, strings.ToLower(search))
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetPageQuery(pageSize int, offset int, search string) *gorm.DB {
	fmt.Println(offset)
	// search by username or display name
	query := database.DB.Preload("UserProfile").Model(&userschema.User{}).
		Limit(pageSize).Offset(offset).Order("created_at DESC").Where("LOWER(username) LIKE ?", "%"+search+"%")
	return query
}
