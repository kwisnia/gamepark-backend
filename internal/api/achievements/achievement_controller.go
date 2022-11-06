package achievements

import "github.com/gin-gonic/gin"

func GetAllAchievementsHandler(c *gin.Context) {
	achievements := GetAllAchievements()
	c.JSON(200, achievements)
}
