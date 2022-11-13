package userschema

import (
	"github.com/kwisnia/inzynierka-backend/internal/api/achievements"
	"github.com/kwisnia/inzynierka-backend/internal/api/dashboard/activity"
	"github.com/kwisnia/inzynierka-backend/internal/api/schema"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email                 string `gorm:"unique"`
	Password              string
	Username              string                               `gorm:"unique"`
	FollowerCount         uint                                 `gorm:"default:0"`
	FollowingCount        uint                                 `gorm:"default:0"`
	UserProfile           UserProfile                          `gorm:"foreignkey:UserID"`
	UserFeatureUnlock     UserFeatureUnlock                    `gorm:"foreignkey:UserID"`
	Lists                 []schema.GameList                    `gorm:"foreignkey:Owner;references:ID"`
	Reviews               []schema.GameReview                  `gorm:"foreignkey:Creator;references:ID"`
	Helpfuls              []schema.ReviewHelpful               `gorm:"foreignkey:UserID;references:ID"`
	Discussions           []schema.GameDiscussion              `gorm:"foreignkey:CreatorID;references:ID"`
	Posts                 []schema.DiscussionPost              `gorm:"foreignkey:CreatorID;references:ID"`
	CompletedAchievements []achievements.AchievementCompletion `gorm:"foreignkey:UserID;references:ID"`
	Followers             []Following                          `gorm:"foreignkey:Followed;references:ID"`
	Following             []Following                          `gorm:"foreignkey:UserID;references:ID"`
	Activities            []activity.UserActivity              `gorm:"foreignkey:UserID;references:ID"`
}

type UserProfile struct {
	gorm.Model
	DisplayName    string
	UserID         uint
	Avatar         *string
	Bio            string `gorm:"default:''"`
	Banner         *string
	BannerPosition float32 `gorm:"default:50"`
	UserScore      int     `gorm:"default:0"`
}

type UserFeatureUnlock struct {
	gorm.Model
	UserID         uint `gorm:"unique" json:"-"`
	Banner         bool `gorm:"default:false" json:"banner"`
	Avatar         bool `gorm:"default:false" json:"avatar"`
	AnimatedAvatar bool `gorm:"default:false" json:"animatedAvatar"`
	AnimatedBanner bool `gorm:"default:false" json:"animatedBanner"`
}

//func (userProfile *User) AfterUpdate(tx *gorm.DB) (err error) {
//	// check if score passed the threshholds
//	userUnlocks := UserFeatureUnlock{}
//	tx.Where("user_id = ?", userProfile.ID).First(&userUnlocks)
//	fmt.Println("am I here?")
//	fmt.Println(userUnlocks.UserID)
//	fmt.Println("userScore", userProfile.UserProfile.UserScore)
//	if userProfile.UserProfile.UserScore >= BannerUnlock {
//		userUnlocks.Banner = true
//	}
//	if userProfile.UserProfile.UserScore >= AnimatedAvatarUnlock {
//		userUnlocks.AnimatedAvatar = true
//	}
//	if userProfile.UserProfile.UserScore >= AnimatedBannerUnlock {
//		userUnlocks.AnimatedBanner = true
//	}
//	return tx.Save(&userUnlocks).Error
//}
