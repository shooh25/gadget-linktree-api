package repository

import "gadget-linktree-api/internal/domain/model/user"

type TokenGenerator interface {
	Generate(u *user.User) (string, error)
}
