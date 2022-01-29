package v1

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) initRecipesRoutes(api *gin.RouterGroup) {
	recipes := api.Group("/recipes", h.userIdentity)
	{
		recipes.GET("", h.getRecipes)
		recipes.POST("", h.createRecipe)
		recipes.GET("/:recipe_id", h.getRecipe)
		recipes.POST("/:recipe_id", h.addRecipeToRecipeBook)
		recipes.PUT("/:recipe_id", h.updateRecipe)
		recipes.DELETE("/:recipe_id", h.deleteRecipe)

		recipes.PUT("/:recipe_id/categories", h.setRecipeCategories)
		recipes.PUT("/favourites/:recipe_id", h.markRecipeFavourite)
		recipes.DELETE("/favourites/:recipe_id", h.unmarkRecipeFavourite)
		recipes.PUT("/liked/:recipe_id", h.likeRecipe)
		recipes.DELETE("/liked/:recipe_id", h.unlikeRecipe)

		recipes.POST("/:recipe_id/pictures", h.uploadRecipePicture)
		recipes.DELETE("/:recipe_id/pictures", h.deleteRecipePicture)

		recipes.GET("/:recipe_id/encryption", h.getRecipeKey)
		recipes.POST("/:recipe_id/encryption", h.uploadRecipeKey)
		recipes.DELETE("/:recipe_id/encryption", h.deleteRecipeKey)

		recipes.GET("/:recipe_id/users", h.getRecipeUsers)
		recipes.POST("/:recipe_id/users", h.setRecipePublicKey)
		recipes.PUT("/:recipe_id/users", h.setRecipePrivateKey)
		recipes.DELETE("/:recipe_id/users/:user_id", h.deleteUserAccess)
	}
}

func (h *Handler) getRecipes(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recipes, err := h.services.GetRecipesByUser(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, recipes)
}

func (h *Handler) createRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var recipe models.Recipe

	if err := c.BindJSON(&recipe); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	recipe.OwnerId = userId
	id, err := 	h.services.Recipes.CreateRecipe(recipe)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": models.RespRecipeAdded,
	})
}


func (h *Handler) addRecipeToRecipeBook(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	err = h.services.AddRecipeToRecipeBook(recipeId, userId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeAdded,
	})
}

func (h *Handler) getRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	recipe, err := h.services.GetRecipeById(recipeId, userId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, recipe)
}

func (h *Handler) updateRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	var recipe models.Recipe

	if err := c.BindJSON(&recipe); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	recipe.Id = recipeId

	if err := h.services.Recipes.UpdateRecipe(recipe, userId); err != nil {
		newResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeUpdated,
	})
}

func (h *Handler) deleteRecipe(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	if err := h.services.DeleteRecipe(recipeId, userId); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeDeleted,
	})
}

func (h *Handler) setRecipeCategories(c *gin.Context) {
	var input models.RecipeCategoriesInput
	var err error
	if err = c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetRecipeCategories(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrRecipeNotInRecipeBook.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespCategoriesUpdated,
	})
}

func (h *Handler) markRecipeFavourite(c *gin.Context) {
	var err error
	input := models.FavouriteRecipeInput {
		Favourite: true,
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.MarkRecipeFavourite(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrRecipeNotInRecipeBook.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespFavouriteStatusUpdated,
	})
}

func (h *Handler) unmarkRecipeFavourite(c *gin.Context) {
	var err error
	input := models.FavouriteRecipeInput {
		Favourite: false,
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.MarkRecipeFavourite(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrRecipeNotInRecipeBook.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespFavouriteStatusUpdated,
	})
}

func (h *Handler) likeRecipe(c *gin.Context) {
	var err error
	input := models.RecipeLikeInput {
		Liked: true,
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetRecipeLike(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeLikeSet,
	})
}

func (h *Handler) unlikeRecipe(c *gin.Context) {
	var err error
	input := models.RecipeLikeInput {
		Liked: false,
	}
	input.UserId, err = getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	input.RecipeId, err = strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetRecipeLike(input)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipeLikeSet,
	})
}

func (h *Handler) uploadRecipePicture(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidFileInput.Error())
		return
	}
	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	buffer := make([]byte, fileHeader.Size)
	_, err = file.Read(buffer)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fileType := http.DetectContentType(buffer)

	fileBytes := bytes.NewReader(buffer)

	url, err := h.services.UploadRecipePicture(c.Request.Context(), recipeId, userId, fileBytes, fileHeader.Size, fileType)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"link": url,
	})
}

func (h *Handler) deleteRecipePicture(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	var picture models.RecipeDeletePictureInput
	if err := c.BindJSON(&picture); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.DeleteRecipePicture(c.Request.Context(), recipeId, userId, picture.PictureName)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, models.ErrUnableDeleteRecipePicture.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespRecipePictureDeleted,
	})
}

func (h *Handler) getRecipeKey(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	url, err := h.services.GetRecipeKey(recipeId, userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"link": url,
	})
}

func (h *Handler) uploadRecipeKey(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidFileInput.Error())
		return
	}
	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	buffer := make([]byte, fileHeader.Size)
	_, err = file.Read(buffer)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fileType := http.DetectContentType(buffer)
	fileBytes := bytes.NewReader(buffer)

	url, err := h.services.UploadRecipeKey(c.Request.Context(), recipeId, userId, fileBytes, fileHeader.Size, fileType)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"link": url,
	})
}

func (h *Handler) deleteRecipeKey(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.DeleteRecipeKey(c.Request.Context(), recipeId, userId)
	if err != nil {
		if err == models.ErrNotOwner {
			newResponse(c, http.StatusBadRequest, models.ErrNotOwner.Error())
		} else {
			newResponse(c, http.StatusInternalServerError, models.ErrUnableDeleteRecipeKey.Error())
		}
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespKeyDeleted,
	})
}

func (h *Handler) setRecipePublicKey(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	var keys models.RecipeKeys
	if err := c.BindJSON(&keys); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetUserPublicKeyForRecipe(recipeId, userId, keys.PublicKey)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespKeySet,
	})
}

func (h *Handler) setRecipePrivateKey(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	var keys models.RecipeKeys
	if err := c.BindJSON(&keys); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetUserPrivateKeyForRecipe(recipeId, userId, keys.PrivateKey)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespKeySet,
	})
}

func (h *Handler) getRecipeUsers(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	userList, err := h.services.GetRecipeUserList(recipeId, userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, userList)
}

func (h *Handler) deleteUserAccess(c *gin.Context) {
	requesterId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	recipeId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}
	userId, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.DeleteUserAccessToRecipe(recipeId, userId, requesterId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespKeySet,
	})
}