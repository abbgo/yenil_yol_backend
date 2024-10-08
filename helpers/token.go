package helpers

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaimForAdmin struct {
	PhoneNumber  string `json:"phone_number"`
	AdminID      string `json:"admin_id"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	jwt.StandardClaims
}

func GenerateAccessTokenForAdmin(phoneNumber, adminID string, isSuperAdmin bool) (string /* string,*/, error) {
	accessTokenTimeOut, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TIMEOUT"))
	if err != nil {
		return "", err
	}
	expirationTimeAccessToken := time.Now().Add(time.Duration(accessTokenTimeOut) * time.Second)

	claimsAccessToken := &JWTClaimForAdmin{
		PhoneNumber:  phoneNumber,
		AdminID:      adminID,
		IsSuperAdmin: isSuperAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeAccessToken.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccessToken)
	accessTokenString, err := accessToken.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	// refreshTokenTimeOut, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TIMEOUT"))
	// if err != nil {
	// 	return "", "", err
	// }
	// expirationTimeRefreshToken := time.Now().Add(time.Duration(refreshTokenTimeOut) * time.Second)
	// claimsRefreshToken := &JWTClaimForAdmin{
	// 	PhoneNumber:  phoneNumber,
	// 	AdminID:      adminID,
	// 	IsSuperAdmin: isSuperAdmin,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: expirationTimeRefreshToken.Unix(),
	// 	},
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)
	// refreshTokenString, err := token.SignedString(JwtKey)
	// if err != nil {
	// 	return "", "", err
	// }

	return accessTokenString /* refreshTokenString, */, nil
}

// func RefreshTokenForAdmin(c *gin.Context) {
// 	tokenStr := c.GetHeader("RefreshToken")
// 	tokenString := strings.Split(tokenStr, " ")[1]

// 	if tokenString == "" {
// 		c.JSON(401, gin.H{
// 			"message": "request does not contain an refresh token",
// 		})
// 		// c.Abort()
// 		return
// 	}

// 	token, err := jwt.ParseWithClaims(
// 		tokenString,
// 		&JWTClaimForAdmin{},
// 		func(token *jwt.Token) (interface{}, error) {
// 			return []byte(JwtKey), nil
// 		},
// 	)

// 	if err != nil {
// 		c.JSON(403, gin.H{
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	claims, ok := token.Claims.(*JWTClaimForAdmin)

// 	if !ok {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": errors.New("couldn't parse claims")})
// 		return
// 	}

// 	if claims.ExpiresAt < time.Now().Local().Unix() {
// 		c.JSON(403, gin.H{
// 			"message": errors.New("token expired"),
// 		})
// 		return
// 	}

// 	accessTokenString, refreshTokenString, err := GenerateAccessTokenForAdmin(claims.PhoneNumber, claims.AdminID, claims.IsSuperAdmin)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":        true,
// 		"access_token":  accessTokenString,
// 		"refresh_token": refreshTokenString,
// 	})
// }
