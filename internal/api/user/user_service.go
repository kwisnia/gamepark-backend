package user

type BasicUserDetails struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	ID          uint   `json:"id"`
}

func GetUserDetails(userName string) *DetailsResponse {
	user := GetByUsername(userName)
	if user == nil {
		return nil
	}
	return &DetailsResponse{
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   user.UserProfile.FirstName,
		LastName:    user.UserProfile.LastName,
		DisplayName: user.UserProfile.DisplayName,
		ID:          user.ID,
		Lists:       user.Lists,
	}
}

func GetBasicUserDetailsByID(userID uint) *BasicUserDetails {
	user := GetByID(userID)
	if user == nil {
		return nil
	}
	return &BasicUserDetails{
		Username:    user.Username,
		DisplayName: user.UserProfile.DisplayName,
		ID:          user.ID,
	}
}
