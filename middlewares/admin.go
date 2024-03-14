package middlewares

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CheckAdmin middleware ahli admin - lere dostup beryar
// gelen request - in admin tarapyndan gelip gelmedigini barlayar
// we admin bolmasa gecirmeyar
func CheckAdmin() gin.HandlerFunc {

	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Token is required")
			return
		}
		var tokenString string

		splitToken := strings.Split(tokenStr, "Bearer ")
		if len(splitToken) > 1 {
			tokenString = splitToken[1]
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid token")
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
			c.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
			return
		}
		claims, ok := token.Claims.(*helpers.JWTClaimForAdmin)
		if !ok {
			c.AbortWithStatusJSON(400, gin.H{"message": "couldn't parse claims"})
			return
		}
		if claims.ExpiresAt < time.Now().Local().Unix() {
			c.AbortWithStatusJSON(403, gin.H{"message": "token expired"})
			return
		}

		db, err := config.ConnDB()
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
			return
		}
		defer db.Close()

		var shop_owner_id string
		db.QueryRow(context.Background(), "SELECT id FROM admins WHERE id = $1 AND deleted_at IS NULL", claims.AdminID).Scan(&shop_owner_id)
		if shop_owner_id == "" {
			c.AbortWithStatusJSON(404, gin.H{"message": "admin not found"})
			return
		}

		c.Set("admin_id", claims.AdminID)
		c.Next()
	}

}

// IsSuperAdmin middleware dine super adminlere dostup beryar
// adminleri gecirmeyar
func IsSuperAdmin() gin.HandlerFunc {

	return func(context *gin.Context) {

		tokenStr := context.GetHeader("Authorization")
		if tokenStr == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, "Token is required")
			return
		}
		var tokenString string

		splitToken := strings.Split(tokenStr, "Bearer ")
		if len(splitToken) > 1 {
			tokenString = splitToken[1]
		} else {
			context.AbortWithStatusJSON(http.StatusBadRequest, "Invalid token")
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
			context.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
			return
		}
		claims, ok := token.Claims.(*helpers.JWTClaimForAdmin)
		if !ok {
			context.AbortWithStatusJSON(400, gin.H{"message": "couldn't parse claims"})
			return
		}
		if claims.ExpiresAt < time.Now().Local().Unix() {
			context.AbortWithStatusJSON(403, gin.H{"message": "token expired"})
			return
		}
		// context.Set("admin_id", claims.AdminID)

		if !claims.IsSuperAdmin {
			context.AbortWithStatusJSON(400, gin.H{"message": "only super_admin can perform this task"})
			return
		}

		context.Next()
	}

}
