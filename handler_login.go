package main

import (
	"log"
	"net/http"
	"piszcodziennie/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (cfg *Config) handlerLoginUser(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	accessToken, refreshToken, err := cfg.loginUser(req.Email, req.Password)

	if err != nil {
		jsonErr, err := utils.ExtractSupabaseError(err.Error())
		if err == nil {
			c.JSON(jsonErr.Code, gin.H{"error": jsonErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
func (cfg *Config) loginUser(email, password string) (string, string, error) {
	session, err := cfg.supabaseClient.Auth.SignInWithEmailPassword(email, password)
	if err != nil {
		log.Println("Failed to login user:", err)
		return "", "", err
	}
	println("userID: ", session.User.ID.String())
	return session.AccessToken, session.RefreshToken, nil
}
