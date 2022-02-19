package v1

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/model"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func getRecipesRequestParamsByCtx(c *gin.Context) *model.RecipesRequestParams {
	var params model.RecipesRequestParams
	if owned, ok := c.GetQuery("owned"); ok {
		params.Owned = owned == "true"
	}
	if minTime, ok := c.GetQuery("min_time"); ok {
		params.MinTime, _ = strconv.Atoi(minTime)
		if params.MinTime < 0 {
			params.MinTime = 0
		}
	}
	if maxTime, ok := c.GetQuery("max_time"); ok {
		params.MaxTime, _ = strconv.Atoi(maxTime)
		if params.MaxTime < 0 {
			params.MaxTime = 0
		}
	}
	if minServings, ok := c.GetQuery("min_servings"); ok {
		params.MinServings, _ = strconv.Atoi(minServings)
		if params.MinServings < 0 {
			params.MinServings = 0
		}
	}
	if maxServings, ok := c.GetQuery("max_servings"); ok {
		params.MaxServings, _ = strconv.Atoi(maxServings)
		if params.MaxServings < 0 {
			params.MaxServings = 0
		}
	}
	if minCalories, ok := c.GetQuery("min_calories"); ok {
		params.MinCalories, _ = strconv.Atoi(minCalories)
		if params.MinCalories < 0 {
			params.MinCalories = 0
		}
	}
	if maxCalories, ok := c.GetQuery("max_calories"); ok {
		params.MaxCalories, _ = strconv.Atoi(maxCalories)
		if params.MaxCalories < 0 {
			params.MaxCalories = 0
		}
	}
	if authorId, ok := c.GetQuery("author_id"); ok {
		params.AuthorId, _ = strconv.Atoi(authorId)
		if params.AuthorId < 0 {
			params.AuthorId = 0
		}
	}
	if sortBy, ok := c.GetQuery("sort_by"); ok {
		if sortBy == "likes" || sortBy == "time" || sortBy == "servings" || sortBy == "calories" {
			params.SortBy = sortBy
		}
	}
	if params.SortBy == "" {
		params.SortBy = "recipe_id"
	}
	if search, ok := c.GetQuery("search"); ok {
		params.Search = search
	}
	if page, ok := c.GetQuery("page"); ok {
		params.Page, _ = strconv.Atoi(page)
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if pageSize, ok := c.GetQuery("page_size"); ok {
		params.PageSize, _ = strconv.Atoi(pageSize)
		if params.PageSize > 50 {
			params.PageSize = 50
		}
	} else {
		params.PageSize = 20
	}
	if languages, ok := c.GetQueryArray("language"); ok {
		params.Languages = languages
	}
	return &params
}

func getUserAndRecipeIdByCtx(c *gin.Context) (int, int, error) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 0, 0, os.ErrInvalid
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return userId, 0, os.ErrInvalid
	}
	return userId, recipeId, nil
}

func getRecipeByCtx(c *gin.Context) (model.Recipe, error) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return model.Recipe{}, os.ErrInvalid
	}
	recipeId, _ := strconv.Atoi(c.Param("recipe_id"))

	if userId < 1 {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return model.Recipe{}, os.ErrInvalid
	}

	var recipe model.Recipe
	if err := c.BindJSON(&recipe); err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return model.Recipe{}, os.ErrInvalid
	}
	recipe.Id = recipeId
	recipe.OwnerId = userId

	return recipe, nil
}

func getUserAndCategoryIdByCtx(c *gin.Context) (int, int, error) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return 0, 0, os.ErrInvalid
	}

	categoryId, err := strconv.Atoi(c.Param("category_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return 0, 0, os.ErrInvalid
	}
	return userId, categoryId, nil
}

func getCategoryByCtx(c *gin.Context) (model.Category, error) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return model.Category{}, os.ErrInvalid
	}
	categoryId, _ := strconv.Atoi(c.Param("category_id"))

	var category model.Category
	if err := c.BindJSON(&category); err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return model.Category{}, os.ErrInvalid
	}
	category.Id = categoryId
	category.UserId = userId

	return category, nil
}

func getFileByCtx(c *gin.Context) (model.MultipartFileInfo, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidFileInput.Error())
		return model.MultipartFileInfo{}, os.ErrInvalid
	}

	defer closeMultipartFile(c, file)

	buffer := make([]byte, header.Size)
	_, err = file.Read(buffer)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return model.MultipartFileInfo{}, os.ErrInvalid
	}
	return model.MultipartFileInfo{
		Name: uuid.NewString(),
		File: bytes.NewReader(buffer),
		Size: header.Size,
		ContentType: http.DetectContentType(buffer),
	}, nil
}

func closeMultipartFile(c *gin.Context, file multipart.File) {
	err := file.Close()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
}

func processKeyError(c *gin.Context, err error)  {
	if err == model.ErrNotOwner {
		newErrorResponse(c, http.StatusBadRequest, model.ErrNotOwner.Error())
	} else {
		newErrorResponse(c, http.StatusInternalServerError, model.ErrUnableDeleteUserKey.Error())
	}
}