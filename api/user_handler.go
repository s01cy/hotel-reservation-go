package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"hotel-reservation/db"
	"hotel-reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
	}
	return c.JSON(u)
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		ctx = context.Background()
		id  = c.Params("id")
	)
	u, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(u)
}
