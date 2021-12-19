package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"idea_server/global"
	"idea_server/model/user"
	"idea_server/model/user/request"
	"idea_server/utils"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	codeExpiration = time.Minute * 10 // 验证码有效期
	codeForbidden  = time.Minute      // 验证码重复发送时间
)

var (
	emailRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type UserBaseService struct {
}

func (u *UserBaseService) Login(email, passwd string) bool {
	if err := global.IDEA_DB.Where("email = ? and passwd = ?", email, passwd).First(&user.User{}).Error; err != nil {
		global.IDEA_LOG.Error("登录失败", zap.Error(err))
		return false
	}
	return true
}

func (u *UserBaseService) Register(regInfo request.Register) (bool, error) {
	// 判断次邮箱是否已注册
	if isExist, err := u.IsExistEmail(regInfo.Email); err != nil {
		return false, errors.New("server internal error")
	} else {
		if isExist {
			return false, errors.New("邮箱已注册")
		}
	}

	//if !u.IsCorrectEmailCode(regInfo.Email, regInfo.Code) {
	//	return false, errors.New("验证码不存在或者已过期")
	//}
	var maxId int
	err := global.IDEA_DB.Select("id").Last(&user.User{}).Row().Scan(&maxId)
	if err != nil {
		return false, err
	}

	avatarInfo := "{\"userId\": \"" + strconv.Itoa(maxId+1) + "\"}"
	res, _ := http.Post("http://127.0.0.1:9998/get_avatar_url", "application/json", bytes.NewBuffer([]byte(avatarInfo)))
	if res != nil {
		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)
		if res.StatusCode != 200 {
			return false, errors.New(string(data))
		}
		// 新建用户
		err = global.IDEA_DB.Create(&user.User{Email: regInfo.Email, Username: regInfo.Username, Passwd: regInfo.Passwd, Avatar: string(data), Weight: 1}).Error
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (u *UserBaseService) GetEmailCode(email, ip string) (code string, err error) {
	if !u.IsAllowedIP(ip) {
		return "", errors.New("请求太频繁")
	}
	// 生成验证码
	code = fmt.Sprintf("%06v", emailRand.Int31n(1000000))

	// 发送邮箱
	emails := make([]string, 0, 1)
	emails = append(emails, email)
	err = utils.SendMail("idea 邮箱验证码", "你的邮箱验证码："+code, emails)
	if err != nil {
		global.IDEA_LOG.Error("邮箱发送失败", zap.Error(err))
		return "", errors.New("邮箱发送失败")
	}

	ctx := context.Background()
	global.IDEA_REDIS.Set(ctx, email, code, codeExpiration)
	global.IDEA_REDIS.Set(ctx, ip, ip, codeForbidden)
	return code, nil
}

func (u *UserBaseService) GetID(email string) string {
	var u2 user.User
	global.IDEA_DB.Where("email = ?", email).First(&u2)
	return strconv.Itoa(int(u2.ID))
}

func (u *UserBaseService) IsExistEmail(email string) (bool, error) {
	if err := global.IDEA_DB.Where("email = ?", email).First(&user.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (u *UserBaseService) IsCorrectEmailCode(email, code string) bool {
	re, err := global.IDEA_REDIS.Get(context.Background(), email).Result()
	if err != nil {
		return false
	}
	if re != code {
		return false
	}
	return true
}

func (u *UserBaseService) IsAllowedIP(ip string) bool {
	_, err := global.IDEA_REDIS.Get(context.Background(), ip).Result()
	if err == redis.Nil {
		return true
	} else {
		return false
	}
}
