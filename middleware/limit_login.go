package middleware

import (
	"context"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"idea_server/global"
	"idea_server/model/user/response"
	"idea_server/utils/constant"
	"math"
	"strconv"
	"time"
)

// 登录数量限制

const (
	spaceName = "login_status:"
	limitNum  = 1
)

func LimitLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取过期时间
		claims := jwt.ExtractClaims(c)
		// 可知毫秒为 0
		sec, _ := math.Modf(claims["exp"].(float64))
		// 由于支持 redis.z score 类型为 float64
		//exp := claims["exp"].(float64)
		// 获取 id
		tokenInfo, _ := c.Get(constant.IdentityKey)
		key := spaceName + tokenInfo.(*response.UserResponse).Email

		ctx := context.Background()

		// 删除过期元素
		now := time.Now().Unix()
		global.IDEA_REDIS.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now,10))


		// 判断是否在登录名单中
		if _, err := global.IDEA_REDIS.ZScore(ctx, key, strconv.Itoa(int(sec))).Result(); err == redis.Nil {
			if cnt, _ := global.IDEA_REDIS.ZCard(ctx, key).Result(); cnt < limitNum { // 如果元素小于限制数直接添加
				err := global.IDEA_REDIS.ZAdd(ctx, key, &redis.Z{Score: sec, Member: sec}).Err()
				if err != nil {
					global.IDEA_LOG.Error("添加有序集合失败", zap.Error(err))
				}
			} else {
				// 判断首个过期时间是否大于当前 token 过期时间
				result, _ := global.IDEA_REDIS.ZRange(ctx, key, 0, 1).Result()
				min, _ := strconv.Atoi(result[0])
				if int(sec) < min {
					global.IDEA_LOG.Error("此帐号已在别处登录")
					c.JSON(401, gin.H{
						"code":    401,
						"message": "此帐号已在别处登录",
					})
					c.Abort()
				} else {
					global.IDEA_REDIS.ZRemRangeByScore(ctx, key, result[0], result[0])
					err := global.IDEA_REDIS.ZAdd(ctx, key, &redis.Z{Score: sec, Member: sec}).Err()
					fmt.Println("添加有序集合成功")
					if err != nil {
						global.IDEA_LOG.Error("添加有序集合失败", zap.Error(err))
					}
				}
			}
		}

		// 继续往下处理
		c.Next()
	}
}
