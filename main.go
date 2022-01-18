package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

var DB *gorm.DB

func main() {
	DB = InitDB()
	//注册路由
	r := gin.Default()
	r.POST("/api/auth/register", Register)

	//
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
func GetDB() *gorm.DB {
	return DB
}
func Register(ctx *gin.Context) {
	db := GetDB()

	//获取请求数据
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//验证数据合法性
	if !(len(name) > 0) {
		name = RandomString(10)
	}
	if len(telephone) != 11 {
		log.Println("telephone:", telephone)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "电话号码需要是11位",
		})
		return
	} else {
		//判断手机号是否已经存在
		if IsExistTelephone(telephone, db) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "手机号已经被注册"})
			return
		}
	}
	if len(password) < 6 {
		log.Println("password:", password)
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "密码需要大于6位",
		})
		return
	}
	//返回结果
	newUser := User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	result := db.Create(&newUser)
	if result.Error != nil {
		log.Println("failed to create user, err: ", result.Error.Error())
	}
	log.Println("name", name)
	ctx.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}
func IsExistTelephone(telephone string, db *gorm.DB) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
func RandomString(n int) string {

	strarr := []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	s := make([]byte, n)
	s1 := rand.NewSource(time.Now().Unix())
	r1 := rand.New(s1)
	for i := 0; i < n; i++ {
		s[i] = strarr[r1.Intn(len(strarr))]
	}
	return string(s)
}

func InitDB() *gorm.DB {
	//初始化mysql数据库
	//driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "gin_try"
	username := "root"
	password := "123"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset:%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}
