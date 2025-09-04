package dtos

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Validate is a function to validate struct using validator package
func Validate(c *gin.Context, obj any) {
	if err := c.ShouldBindJSON(obj); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, e := range validationErrs {
				errors[e.Field()] = e.Tag()
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
