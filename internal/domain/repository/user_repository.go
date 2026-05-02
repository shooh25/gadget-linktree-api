package repository

import "gadget-linktree-api/internal/domain/model/user"

type UserRepository interface {
	ExistsByGoogleId(googleId string) (bool, error) // 登録済みか判定
	Save(user user.User) error                      // ユーザー保存
}