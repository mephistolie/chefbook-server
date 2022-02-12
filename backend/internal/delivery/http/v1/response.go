package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/pkg/logger"
	"net/http"
)

var (
	RespActivationLink      = "profile activation link has been sent to email"
	RespProfileActivated    = "profile is activated"
	RespSignOutSuccessfully = "signed out successfully"

	RespUsernameChanged = "username successfully changed"
	RespAvatarDeleted   = "avatar has been deleted"
	RespKeySet          = "encrypted key set"
	RespKeyDeleted      = "encrypted key deleted"

	RespRecipeAdded            = "recipe has been added"
	RespRecipeUpdated          = "recipe has been updated"
	RespRecipeDeleted          = "recipe has been deleted"
	RespCategoriesUpdated      = "categories has been updated"
	RespFavouriteStatusUpdated = "favourite status has been updated"
	RespRecipeLikeSet          = "recipe like status has been set"
	RespRecipePictureDeleted   = "picture has been deleted"

	RespCategoryAdded   = "category has been added"
	RespCategoryUpdated = "category has been updated"
	RespCategoryDeleted = "category has been deleted"

	RespShoppingListUpdated = "shopping list has been updated"
)

type idResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type messageResponse struct {
	Message string `json:"message"`
}

type linkResponse struct {
	Link string `json:"link"`
}

func newIdResponse(c *gin.Context, id int, msg string) {
	c.JSON(http.StatusOK,idResponse{id, msg})
}

func newMessageResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, messageResponse{message})
}

func newLinkResponse(c *gin.Context, link string) {
	c.JSON(http.StatusOK, linkResponse{link})
}

func newErrorResponse(c *gin.Context, statusCode int, errorMessage string) {
	logger.Errorf(errorMessage)
	c.AbortWithStatusJSON(statusCode, messageResponse{errorMessage})
}
