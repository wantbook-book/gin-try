package controller

import (
	"github.com/gin-gonic/gin"
	"hannibal/gin-try/model"
	"hannibal/gin-try/repository"
	"hannibal/gin-try/response"
	"hannibal/gin-try/vo"
	"strconv"
)

//避免同一个包中的方法名冲突，使用结构体方法，相当命名空间
type ICategoryController interface {
	RestController
}

type CategoryController struct {
	//DB *gorm.DB
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	//db := common.GetDB()
	categoryRepository := repository.NewCategoryRepository()
	categoryRepository.DB.AutoMigrate(model.Category{})
	return CategoryController{Repository: categoryRepository}
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

	//创建文章
	if createCategory, err := c.Repository.Create(requestCategory.Name); err != nil {
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

	//查询数据库,该id的文章是否存在
	category := model.Category{ID: uint(id)}
	if dbCategory, err := c.Repository.Update(category, requestCategory.Name); err != nil {
		//更新失败
		response.Fail(ctx, nil, "数据验证错误：分类不存在")
		return
	} else {
		response.Success(ctx, gin.H{"category": dbCategory}, "更新文章成功")

		return
	}

}

func (c CategoryController) Show(ctx *gin.Context) {
	//获取path中的id数据
	id, _ := strconv.Atoi(ctx.Params.ByName("id"))
	//数据库查询文章
	if dbCategory, err := c.Repository.SelectById(id); err != nil {
		//文章不存在,返回
		response.Fail(ctx, nil, "数据验证错误：分类不存在")
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
	if err := c.Repository.DeleteById(id); err != nil {
		//删除失败
		response.Fail(ctx, nil, "删除失败")
		return
	} else {
		//删除成功
		response.Success(ctx, nil, "删除成功")
		return
	}

}
