// Package converter содержит функции преобразования между транспортными DTO
// и доменными моделями IAM-сервиса.
package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mbakhodurov/homeworks2/week7/iam/internal/model"
	commonv1 "github.com/mbakhodurov/homeworks2/week7/shared/pkg/proto/common/v1"
)

// UserToDTO конвертирует доменную модель пользователя в транспортный DTO.
func UserToDTO(u model.User) *commonv1.User {
	dto := &commonv1.User{
		Uuid: u.UUID.String(),
		Info: &commonv1.UserInfo{
			Login: u.Login,
		},
		CreatedAt: timestamppb.New(u.CreatedAt),
	}

	if u.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*u.UpdatedAt)
	}

	return dto
}

// SessionToDTO конвертирует доменную модель сессии в транспортный DTO.
func SessionToDTO(s model.Session) *commonv1.Session {
	dto := &commonv1.Session{
		Uuid:      s.UUID.String(),
		CreatedAt: timestamppb.New(s.CreatedAt),
		ExpiresAt: timestamppb.New(s.ExpiresAt),
	}

	if s.UpdatedAt != nil {
		dto.UpdatedAt = timestamppb.New(*s.UpdatedAt)
	}

	return dto
}
