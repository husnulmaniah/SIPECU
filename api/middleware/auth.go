package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "sipecut_super_secret_key_12345"
	}
	jwtSecretKey = []byte(secret)
}

// GenerateToken creates an access token and a refresh token
func GenerateToken(nip string, role string, menus []string) (string, string, error) {
	// 1. Access Token (expires in 24 hours for easier local testing/use)
	accessTokenClaims := jwt.MapClaims{
		"nip":               nip,
		"role":              role,
		"daftar_menu_akses": menus,
		"exp":               time.Now().Add(24 * time.Hour).Unix(),
		"iat":               time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessStr, err := accessToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	// 2. Refresh Token (expires in 7 days)
	refreshTokenClaims := jwt.MapClaims{
		"nip":  nip,
		"role": role,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshStr, err := refreshToken.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}

// ValidateToken parses and validates a JWT token string
func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// AuthMiddleware intercepts requests and validates JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var tokenStr string

		if authHeader == "" {
			// Fallback: check token query parameter for file downloads/links
			tokenStr = c.Query("token")
			if tokenStr == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
				c.Abort()
				return
			}
		} else {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be Bearer <token>"})
				c.Abort()
				return
			}
			tokenStr = parts[1]
		}

		claims, err := ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid token: %v", err)})
			c.Abort()
			return
		}

		// Set variables to Gin Context
		c.Set("nip", claims["nip"].(string))
		c.Set("role", claims["role"].(string))

		// Parse menus
		if menusVal, ok := claims["daftar_menu_akses"]; ok {
			if menusInterface, ok := menusVal.([]interface{}); ok {
				var menus []string
				for _, m := range menusInterface {
					if s, ok := m.(string); ok {
						menus = append(menus, s)
					}
				}
				c.Set("menus", menus)
			}
		}

		c.Next()
	}
}

// RoleMiddleware checks if the authenticated user has the required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not identified"})
			c.Abort()
			return
		}

		userRole := roleVal.(string)
		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// MenuAccessCheck memvalidasi hak akses modul menu pegawai secara real-time
func MenuAccessCheck(requiredMenuCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, existsRole := c.Get("role")
		if existsRole && role.(string) == "admin" {
			// Admin has access to all menus
			c.Next()
			return
		}

		menus, exists := c.Get("menus")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak. Tidak memiliki menu akses."})
			c.Abort()
			return
		}

		allowed := false
		for _, m := range menus.([]string) {
			if m == requiredMenuCode {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke modul ini"})
			c.Abort()
			return
		}
		c.Next()
	}
}
