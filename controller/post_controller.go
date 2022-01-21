package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hannibal/gin-try/common"
	"hannibal/gin-try/model"
	"hannibal/gin-try/response"
	"hannibal/gin-try/vo"
	"log"
	"strconv"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{
		DB: db,
	}
}

func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	//获取文章数据
	//检验数据合法性，交由binding完成
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "请求数据格式不正确")
		return
	}
	//获取用户信息，合成完整的数据库文章记录
	user, _ := ctx.Get("user")
	//创建记录
	createPost := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}
	if err := p.DB.Create(&createPost).Error; err != nil {
		//创建失败
		panic(err)
		//response.Fail(ctx, nil, "创建失败")
		return
	} else {
		//创建成功
		response.Success(ctx, nil, "创建成功")
		return
	}

}

func (p PostController) Update(ctx *gin.Context) {
	//获取path的id
	requestId := ctx.Params.ByName("id")
	//判断文章是否存在
	var dbPost model.Post
	if err := p.DB.Where("id", requestId).First(&dbPost).Error; err == gorm.ErrRecordNotFound {
		//不存在返回
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	//判断登录用户是否为文章作者
	user, _ := ctx.Get("user")
	if user.(model.User).ID != dbPost.UserId {
		//文章不属于当前用户
		response.Fail(ctx, nil, "修改失败：文章不属于当前登录用户")
		return
	}

	var requestPost vo.CreatePostRequest
	//获取文章数据
	//检验数据合法性，交由binding完成
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "请求数据格式不正确")
		return
	}

	//更新数据
	dbPost.Content = requestPost.Content
	dbPost.CategoryId = requestPost.CategoryId
	dbPost.HeadImg = requestPost.HeadImg
	dbPost.Title = requestPost.Title
	if err := p.DB.Save(&dbPost).Error; err != nil {
		//更新失败
		panic(err)
		return
	} else {
		//更新成功
		response.Success(ctx, gin.H{"post": dbPost}, "修改成功")
		return
	}
}

func (p PostController) Show(ctx *gin.Context) {
	//获取path的id
	requestId := ctx.Params.ByName("id")
	//判断文章是否存在
	var dbPost model.Post
	if err := p.DB.Preload("Category").Where("id = ?", requestId).First(&dbPost).Error; err == gorm.ErrRecordNotFound {
		//不存在返回
		response.Fail(ctx, nil, "文章不存在")
		return
	} else {
		response.Success(ctx, gin.H{"post": dbPost}, "查询成功")
	}
}

func (p PostController) Delete(ctx *gin.Context) {
	//获取path的id
	requestId := ctx.Params.ByName("id")
	log.Println(requestId)
	//判断文章是否存在
	var dbPost model.Post
	if err := p.DB.Where("id = ?", requestId).First(&dbPost).Error; err == gorm.ErrRecordNotFound {
		//不存在返回
		response.Fail(ctx, nil, "文章不存在")
		return
	}
	//判断当前登录用户是否为文章作者
	user, _ := ctx.Get("user")
	if user.(model.User).ID != dbPost.UserId {
		//文章不属于当前登录用户
		response.Fail(ctx, nil, "删除失败：文章不属于当前登录用户")
		return
	}

	if err := p.DB.Where("id = ?", requestId).Delete(model.Post{}).Error; err != nil {
		//删除失败
		panic(err)
		return
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

func (p PostController) PageList(ctx *gin.Context) {
	//获取分页参数
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("pageNumber", "1"))
	log.Println("pageSize: ", pageSize)
	//查询posts
	var total int64
	var posts []model.Post
	p.DB.Model(model.Post{}).Order("created_at desc").Offset((pageNumber - 1) * pageSize).Limit(pageSize).Find(&posts)
	//返回posts数量
	p.DB.Model(model.Post{}).Count(&total)
	response.Success(ctx, gin.H{"data": posts, "total": total}, "查询数量成功")
}
