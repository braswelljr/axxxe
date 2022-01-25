package helper

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// MatchUserTypeToUID : Match Role to userid
func MatchUserTypeToUID(ctx *fiber.Ctx, userId string) error {
	// get user type and id
	userType := ctx.Get("role")
	uid := ctx.Get("user_id")

	// check for user type before access is granted
	if userType == "USER" && uid != userId {
		return errors.New("unauthorised to access this resource")
	}

	return CheckUserType(ctx, userType)
}

func CheckUserType(ctx *fiber.Ctx, role string) (err error) {
	// get user type and id
	userType := ctx.Get("role")
	err = nil

	// check for type
	if userType != role {
		errors.New("unauthorised to access this resource")
	}

	return err
}
