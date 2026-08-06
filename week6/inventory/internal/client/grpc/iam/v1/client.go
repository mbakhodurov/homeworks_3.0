package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	authv1 "github.com/mbakhodurov/homeworks2/week6/shared/pkg/proto/auth/v1"
)

// Client — обёртка над proto-клиентом AuthService (IAMService).
type Client struct {
	authClient authv1.AuthServiceClient
}

// New создаёт новый клиент AuthService.
func New(authClient authv1.AuthServiceClient) *Client {
	return &Client{authClient: authClient}
}

// Whoami проверяет сессию и возвращает UUID её владельца.
func (c *Client) Whoami(ctx context.Context, sessionUUID string) (uuid.UUID, error) {
	resp, err := c.authClient.Whoami(ctx, &authv1.WhoamiRequest{SessionUuid: sessionUUID})
	if err != nil {
		return uuid.Nil, err
	}

	userUUID, err := uuid.Parse(resp.GetUser().GetUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("распарсить uuid пользователя: %w", err)
	}

	return userUUID, nil
}
