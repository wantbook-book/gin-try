package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"hannibal/gin-try/common"
	"hannibal/gin-try/dto"
	"hannibal/gin-try/model"
	"hannibal/gin-try/response"
	"hannibal/gin-try/utils"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {
	db := common.GetDB()

	//获取请求数据
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//验证数据合法性
	if !(len(name) > 0) {
		name = utils.RandomString(10)
	}
	if len(telephone) != 11 {
		log.Println("telephone:", telephone)
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "电话号码需要是11位")
		return
	} else {
		//判断手机号是否已经存在
		if IsExistTelephone(telephone, db) {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号已经被注册")
			return
		}
	}

	if len(password) < 6 {
		log.Println("password:", password)
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码需要大于6位")

		return
	}
	//加密密码
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "服务器加密错误")
		return
	}

	//返回结果
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	result := db.Create(&newUser)
	if result.Error != nil {
		log.Println("failed to create user, err: ", result.Error.Error())
	}
	log.Println("name", name)
	response.Success(ctx, nil, "注册成功")

}
func Login(ctx *gin.Context) {
	//获取数据
	db := common.GetDB()
	//name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	//核对数据
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	//手机号是否存在
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在或密码错误")
		return
	} else {
		//密码是否正确
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在或密码错误")
			return
		}
	}

	//生成返回token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "服务器生成jwt错误")
		return
	}
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func IsExistTelephone(telephone string, db *gorm.DB) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func Info(ctx *gin.Context) {
	//已通过token中间件
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": dto.ToUserDto(user.(model.User))}, "成功获取用户信息")
}
