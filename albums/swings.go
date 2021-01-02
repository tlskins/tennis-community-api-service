package albums

import (
	"context"

	t "github.com/tennis-community-api-service/albums/types"
)

func (u *AlbumsService) CreateSwing(ctx context.Context, data *t.SwingVideo) (*t.SwingVideo, error) {
	return u.Store.CreateSwing(data)
}
