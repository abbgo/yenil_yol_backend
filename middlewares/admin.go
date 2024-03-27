package middlewares

import (
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CheckAdmin middleware ahli admin - lere dostup beryar
// gelen request - in admin tarapyndan gelip gelmedigini barlayar
// we admin bolmasa gecirmeyar
func CheckToken(position string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Token is required", "status": false})
			return
		}
		var tokenString string

		splitToken := strings.Split(tokenStr, "Bearer ")
		if len(splitToken) > 1 {
			tokenString = splitToken[1]
		} else {
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid token", "status": false})
			return
		}

		token, err := jwt.ParseWithClaims(
			tokenString,
			&helpers.JWTClaimForAdmin{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(helpers.JwtKey), nil
			},
		)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"message": err.Error(), "status": false})
			return
		}
		claims, ok := token.Claims.(*helpers.JWTClaimForAdmin)
		if !ok {
			c.AbortWithStatusJSON(400, gin.H{"message": "couldn't parse claims", "status": false})
			return
		}
		if claims.ExpiresAt < time.Now().Local().Unix() {
			c.AbortWithStatusJSON(403, gin.H{"message": "token expired", "status": false})
			return
		}

		db, err := config.ConnDB()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": err.Error(), "status": false})
			return
		}
		defer db.Close()

		var tableName = ""
		switch position {
		case "admin":
			tableName = "admins"
		case "shop_owner":
			tableName = "shop_owners"
		case "customer":
			tableName = "customers"
		default:
			c.AbortWithStatusJSON(400, gin.H{"message": "position not found", "status": false})
			return
		}

		if err := helpers.ValidateRecordByID(tableName, claims.AdminID, "NULL", db); err != nil {
			c.AbortWithStatusJSON(404, gin.H{"message": position + " not found", "status": false})
			return
		}

		c.Set(position+"_id", claims.AdminID)
		c.Next()
	}
}

func CheckTokenAdminOrShopOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Token is required", "status": false})
			return
		}
		var tokenString string

		splitToken := strings.Split(tokenStr, "Bearer ")
		if len(splitToken) > 1 {
			tokenString = splitToken[1]
		} else {
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid token", "status": false})
			return
		}

		token, err := jwt.ParseWithClaims(
			tokenString,
			&helpers.JWTClaimForAdmin{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(helpers.JwtKey), nil
			},
		)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"message": err.Error(), "status": false})
			return
		}
		claims, ok := token.Claims.(*helpers.JWTClaimForAdmin)
		if !ok {
			c.AbortWithStatusJSON(400, gin.H{"message": "couldn't parse claims", "status": false})
			return
		}
		if claims.ExpiresAt < time.Now().Local().Unix() {
			c.AbortWithStatusJSON(403, gin.H{"message": "token expired", "status": false})
			return
		}

		db, err := config.ConnDB()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": err.Error(), "status": false})
			return
		}
		defer db.Close()

		// eger database - de admin ya-da customer yok bolsa onda error return edilyar
		countOfError := 0
		if err := helpers.ValidateRecordByID("admins", claims.AdminID, "NULL", db); err != nil {
			countOfError++
		}
		if err := helpers.ValidateRecordByID("shop_owners", claims.AdminID, "NULL", db); err != nil {
			countOfError++
		}
		if countOfError > 1 {
			c.AbortWithStatusJSON(404, gin.H{"message": "record not found", "status": false})
			return
		}

		c.Set("id", claims.AdminID)
		c.Next()
	}
}

// IsSuperAdmin middleware dine super adminlere dostup beryar
// adminleri gecirmeyar
func IsSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Token is required", "status": false})
			return
		}
		var tokenString string

		splitToken := strings.Split(tokenStr, "Bearer ")
		if len(splitToken) > 1 {
			tokenString = splitToken[1]
		} else {
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid token", "status": false})
			return
		}

		token, err := jwt.ParseWithClaims(
			tokenString,
			&helpers.JWTClaimForAdmin{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(helpers.JwtKey), nil
			},
		)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"message": err.Error(), "status": false})
			return
		}
		claims, ok := token.Claims.(*helpers.JWTClaimForAdmin)
		if !ok {
			c.AbortWithStatusJSON(400, gin.H{"message": "couldn't parse claims", "status": false})
			return
		}
		if claims.ExpiresAt < time.Now().Local().Unix() {
			c.AbortWithStatusJSON(403, gin.H{"message": "token expired", "status": false})
			return
		}

		if !claims.IsSuperAdmin {
			c.AbortWithStatusJSON(400, gin.H{"message": "only super_admin can perform this task", "status": false})
			return
		}

		c.Next()
	}
}
