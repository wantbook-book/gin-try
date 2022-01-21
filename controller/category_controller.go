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

//避免同一个包中的方法名冲突，使用结构体方法，相当命名空间
type ICategoryController interface {
	RestController
}

type CategoryController struct {
	DB *gorm.DB
}

func NewCategoryController() ICategoryController {
	db := common.GetDB()
	db.AutoMigrate(model.Category{})
	return CategoryController{DB: db}
}
func (c CategoryController) Create(ctx *gin.Context) {
	var requestCategory vo.CreateCategoryRequest
	//获取请求数据
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		//binding指定必须包含name，name=""就认为没有提供name
		//验证数据
		response.Fail(ctx, nil, "数据验证失败：分类名称必须不为空")
		return
	}
	createCategory := model.Category{Name: requestCategory.Name}
	//创建文章
	if result := c.DB.Create(&createCategory); result.RowsAffected == 0 || result.Error != nil {
		//创建失败
		response.Fail(ctx, nil, "创建失败")
		return
	} else {
		//返回响应
		response.Success(ctx, gin.H{"category": createCategory}, "文章创建成功")
	}

}

func (c CategoryController) Update(ctx *gin.Context) {
	//获取分类id
	id, _ := strconv.Atoi(ctx.Params.ByName("id"))
	//localhost:9000/categories/:id必然会有id，不然就404

	//获取分类名称更新数据
	var requestCategory vo.CreateCategoryRequest
	//使用ctx.Bind 如果报错的话,返回的就是普通文本，后续的ctx.JSON不能对header的content-type进行重写
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		//binding指定必须包含name，name=""就认为没有提供name
		//验证数据
		response.Fail(ctx, nil, "数据验证失败：分类名称必须不为空")
		return
	}

	var dbCategory model.Category
	//查询数据库,该id的文章是否存在
	if result := c.DB.First(&dbCategory, id); result.Error != nil {
		//未能找到或其他错误
		response.Fail(ctx, nil, "数据验证错误：未能找到分类")
		return
	} else {
		//存在，更新对应id的文章
		c.DB.Model(&dbCategory).Update("name", requestCategory.Name)

		//返回响应
		dbCategory.Name = requestCategory.Name
		response.Success(ctx, gin.H{"category": dbCategory}, "更新文章成功")
		return
	}

}

func (c CategoryController) Show(ctx *gin.Context) {
	//获取path中的id数据
	id, _ := strconv.Atoi(ctx.Params.ByName("id"))
	//数据库查询文章
	var dbCategory model.Category
	if result := c.DB.First(&dbCategory, id); result.Error != nil {
		//文章不存在,返回
		response.Fail(ctx, nil, "数据验证错误：未能找到分类")
		return
	} else {
		//文章存在，返回
		response.Success(ctx, gin.H{"category": dbCategory}, "成功查询到分类")
		return
	}

}

func (c CategoryController) Delete(ctx *gin.Context) {
	//获取path中的id参数
	id, _ := strconv.Atoi(ctx.Params.ByName("id"))
	//从结构体 automigrate创建的表，先传结构体model参数可以指定删除的是哪张表
	if result := c.DB.Delete(model.Category{}, id); result.Error != nil || result.RowsAffected == 0 {
		//删除失败
		log.Printf("%+v", result)
		response.Fail(ctx, nil, "删除失败")
		return
	} else {
		//删除成功
		log.Printf("%+v", result)
		response.Success(ctx, nil, "删除成功")
		return
	}

}
