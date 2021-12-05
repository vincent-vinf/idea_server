package utils

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"idea_server/utils/constant"
	"strconv"
)

func GetJwtId(c *gin.Context) uint {
	claims := jwt.ExtractClaims(c)
	id, _ := strconv.Atoi(claims[constant.IdentityKey].(string))
	return uint(id)
}
