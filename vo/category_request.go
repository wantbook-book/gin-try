package vo

//对应前台字段
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
