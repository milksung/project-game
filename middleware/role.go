package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type connectDB struct {
	db *gorm.DB
}

func Role(db *gorm.DB) connectDB {
	return connectDB{db}
}

func (r connectDB) CheckAdmin(perName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		role := c.MustGet("role").(string)

		if role != "SUPER_ADMIN" {

			adminId, err := c.Get("adminId")
			if !err {
				c.AbortWithStatusJSON(401, gin.H{
					"message": "Unauthorized",
				})
				return
			}

			var result int64

			if err := r.db.Table("Admins a").
				Joins("LEFT JOIN Admin_group_permissions agp ON agp.group_id = a.admin_group_id").
				Joins("LEFT JOIN Permissions p ON p.id = agp.permission_id").
				Select("p.id").
				Where("a.id = ? AND p.permission_key = ?", adminId, perName).
				Count(&result).Error; err != nil {

				if err == gorm.ErrRecordNotFound {
					c.AbortWithStatusJSON(403, gin.H{
						"message": "Permission Denied",
					})
					return
				}

				c.AbortWithStatusJSON(500, gin.H{
					"message": "Internal Server Error",
				})
				return
			}

			if result < 1 {
				c.AbortWithStatusJSON(403, gin.H{
					"message": "Permission Denied",
				})
				return
			}
			println("result", result)
		}
		c.Next()
	}
}
