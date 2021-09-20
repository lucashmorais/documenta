package controllers

import "github.com/gofiber/fiber/v2"

func UserCanFinishReviewSequence(c *fiber.Ctx) bool {
	return len(CoreGetLoggedUserPermissions(c, 17)) > 0
}

func UserCanFinishApprovalSequence(c *fiber.Ctx) bool {
	return len(CoreGetLoggedUserPermissions(c, 18)) > 0
}
