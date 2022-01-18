package middleware

import (
	"github.com/gin-gonic/gin"
	"hannibal/gin-try/common"
	"hannibal/gin-try/model"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取tokenstring
		tokenString := ctx.GetHeader("Authorization")
		//检验token格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}
		//检验token
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		//解析出错，或者token已经过期,则需要重新登录
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}
		//通过验证
		userID := claims.UserId
		db := common.GetDB()
		var user model.User
		db.First(&user, userID)
		if user.ID == 0 {
			//用户不存在,用户已经被删除
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		} else {
			//存在
			ctx.Set("user", user)
		}
		ctx.Next()
	}
}
