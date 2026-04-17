package response

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FetchErrorOrEmpty converts not-found fetch errors into a successful empty response.
// Returns true when a response has been written and the caller should return.
func FetchErrorOrEmpty(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		Success(c, "")
		return true
	}

	Error(c, 500, err.Error())
	return true
}
