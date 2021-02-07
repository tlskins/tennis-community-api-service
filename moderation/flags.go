package albums

import (
	"context"
	"time"

	t "github.com/tennis-community-api-service/moderation/types"
)

func (u *ModerationService) CreateCommentFlag(ctx context.Context, data *t.CommentFlag) (*t.CommentFlag, error) {
	return u.Store.CreateCommentFlag(data)
}

func (u *ModerationService) UpdateCommentFlag(ctx context.Context, data *t.UpdateCommentFlag) (*t.CommentFlag, error) {
	return u.Store.UpdateCommentFlag(data)
}

func (u *ModerationService) RecentFlaggedComments(ctx context.Context, start, end time.Time, resolved *bool, limit, offset int) ([]*t.CommentFlag, error) {
	return u.Store.RecentFlaggedComments(start, end, resolved, limit, offset)
}

func (u *ModerationService) CreateAlbumFlag(ctx context.Context, data *t.AlbumFlag) (*t.AlbumFlag, error) {
	return u.Store.CreateAlbumFlag(data)
}

func (u *ModerationService) UpdateAlbumFlag(ctx context.Context, data *t.UpdateAlbumFlag) (*t.AlbumFlag, error) {
	return u.Store.UpdateAlbumFlag(data)
}

func (u *ModerationService) RecentFlaggedAlbums(ctx context.Context, start, end time.Time, resolved *bool, limit, offset int) ([]*t.AlbumFlag, error) {
	return u.Store.RecentFlaggedAlbums(start, end, resolved, limit, offset)
}
