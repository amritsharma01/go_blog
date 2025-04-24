package validators

import (
	"crud_api/models"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	UserValidator *validator.Validate
)

// Call this once during app startup (like in main.go)
func InitValidators(db *gorm.DB) {

	UserValidator = validator.New()

	// --- Register unique email validator
	UserValidator.RegisterValidation("uniqueEmail", func(fl validator.FieldLevel) bool {
		user := fl.Parent().Interface().(models.User)
		var count int64
		db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
		return count == 0
	})

}
