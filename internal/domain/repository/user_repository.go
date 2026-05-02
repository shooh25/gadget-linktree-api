package repository

import "gadget-linktree-api/internal/domain/model/user"

type UserRepository interface {
	// 登録済みか判定
	ExistsByGoogleId(googleId string) (bool, error)

	// ユーザー保存
	Save(u user.User) error

	// ユーザー取得
	FindByUserId(userId user.UserId) (*user.User, error)
}