package main

import (
	"net/http"
	"piszcodziennie/utils"

	"github.com/gin-gonic/gin"
	"github.com/supabase-community/gotrue-go/types"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (cfg *Config) handlerRegisterUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := cfg.registerUser(req.Email, req.Password)

	if err != nil {
		jsonErr, err := utils.ExtractSupabaseError(err.Error())
		if err == nil {
			c.JSON(jsonErr.Code, gin.H{"error": jsonErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered", "user_id": user.User.ID})
}

func (cfg *Config) registerUser(email, password string) (*types.SignupResponse, error) {
	user, err := cfg.supabaseClient.Auth.Signup(types.SignupRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return &types.SignupResponse{}, err
	}

	return user, nil
}
