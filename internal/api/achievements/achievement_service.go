package achievements

func CompleteAchievement(userID uint, achievementID uint) error {
	_, err := GetByID(achievementID)
	if err != nil {
		return err
	}
	return CreateAchievementCompletion(userID, achievementID)
}
