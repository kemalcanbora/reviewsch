package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		role    string
		wantErr bool
	}{
		{
			name:    "valid token generation",
			userID:  "123",
			role:    "admin",
			wantErr: false,
		},
		{
			name:    "empty userID",
			userID:  "",
			role:    "admin",
			wantErr: false,
		},
		{
			name:    "empty role",
			userID:  "123",
			role:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.role)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token)

			// Verify token contents
			claims := &Claims{}
			parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			assert.NoError(t, err)
			assert.True(t, parsedToken.Valid)
			assert.Equal(t, tt.userID, claims.UserID)
			assert.Equal(t, tt.role, claims.Role)
			assert.NotNil(t, claims.ExpiresAt)
			assert.True(t, time.Now().Before(claims.ExpiresAt.Time))
		})
	}
}

func TestAdminAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupAuth      func() string
		expectedCode   int
		expectedUserID string
		expectedRole   string
	}{
		{
			name: "valid token",
			setupAuth: func() string {
				token, _ := GenerateToken("123", "admin")
				return "Bearer " + token
			},
			expectedCode:   http.StatusOK,
			expectedUserID: "123",
			expectedRole:   "admin",
		},
		{
			name: "missing auth header",
			setupAuth: func() string {
				return ""
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid token format",
			setupAuth: func() string {
				return "InvalidFormat token"
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid token",
			setupAuth: func() string {
				return "Bearer invalid.token.here"
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "expired token",
			setupAuth: func() string {
				claims := &Claims{
					UserID: "123",
					Role:   "admin",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						NotBefore: jwt.NewNumericDate(time.Now()),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString(secretKey)
				return "Bearer " + signedToken
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			router := gin.New()
			router.Use(AdminAuth())

			var capturedUserID, capturedRole string
			router.GET("/test", func(c *gin.Context) {
				// Properly handle type assertions
				if userID, exists := c.Get("userID"); exists {
					capturedUserID = userID.(string)
				}
				if role, exists := c.Get("role"); exists {
					capturedRole = role.(string)
				}
				c.Status(http.StatusOK)
			})

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			if auth := tt.setupAuth(); auth != "" {
				req.Header.Set("Authorization", auth)
			}

			// Perform request
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, tt.expectedUserID, capturedUserID)
				assert.Equal(t, tt.expectedRole, capturedRole)
			}
		})
	}
}
