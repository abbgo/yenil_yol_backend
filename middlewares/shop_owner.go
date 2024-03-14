package middlewares

import (
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// CheckShopOwner middleware ahli shop_owner - lere dostup beryar
// gelen request - in shop_owner tarapyndan gelip gelmedigini barlayar
// we shop_owner bolmasa gecirmeyar
func CheckShopOwner() gin.HandlerFunc {
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

		if err := helpers.ValidateRecordByID("shop_owners", claims.AdminID, "NULL", db); err != nil {
			c.AbortWithStatusJSON(404, gin.H{"message": "customer not found"})
			return
		}

		c.Set("shop_owner_id", claims.AdminID)
		c.Next()
	}
}
