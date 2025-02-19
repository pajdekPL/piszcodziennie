package main

import "fmt"

// TODO: to be done
func (cfg *Config) refreshToken(refreshToken string) (string, string, error) {
	session, err := cfg.supabaseClient.Auth.RefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	newAccessToken := session.AccessToken
	newRefreshToken := session.RefreshToken

	fmt.Println("New Access Token:", newAccessToken)
	fmt.Println("New Refresh Token:", newRefreshToken)

	return newAccessToken, newRefreshToken, nil
}
