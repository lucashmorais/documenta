package controllers

import "github.com/gofiber/fiber/v2"

type PermissionIDValue int

const (
	FinishReviewSequence   PermissionIDValue = 17
	FinishApprovalSequence PermissionIDValue = 18
)

func UserCanFinishReviewSequence(c *fiber.Ctx) bool {
	return len(CoreGetLoggedUserPermissions(c, int(FinishReviewSequence))) > 0
}

func UserCanFinishApprovalSequence(c *fiber.Ctx) bool {
	return len(CoreGetLoggedUserPermissions(c, int(FinishApprovalSequence))) > 0
}
