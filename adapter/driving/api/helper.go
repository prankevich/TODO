package http

import (
	"TODO/pkg"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) extractTokenFromHeader(c *gin.Context, headerKey string) (string, error) {
	header := c.GetHeader(headerKey)

	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid authorization header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("empty token")
	}

	return headerParts[1], nil
}

func (s *Server) generateNewTokenPair(userID int) (string, string, error) {
	// сгенерировать токен (браслет)
	accessToken, err := pkg.GenerateToken(userID,
		s.cfg.AuthParams.AccessTokenTllMinutes,
		false)
	if err != nil {
		return "", err
	}

	refreshToken, err := pkg.GenerateToken(userID,
		s.cfg.AuthParams.RefreshTokenTllDays,
		true)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
