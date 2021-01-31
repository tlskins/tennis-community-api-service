package albums

import (
	"context"
	"time"

	t "github.com/tennis-community-api-service/albums/types"
	"github.com/tennis-community-api-service/pkg/enums"
)

func (u *AlbumsService) GetUserAlbums(ctx context.Context, userID string) ([]*t.Album, error) {
	return u.Store.GetAlbumsByUser(userID)
}

func (u *AlbumsService) GetPublicAlbums(ctx context.Context, homeApproved *bool) ([]*t.Album, error) {
	return u.Store.GetPublicAlbums(homeApproved)
}

func (u *AlbumsService) GetFriendsAlbums(ctx context.Context, userID string) ([]*t.Album, error) {
	return u.Store.GetFriendsAlbums(userID)
}

func (u *AlbumsService) GetAlbum(ctx context.Context, id string) (*t.Album, error) {
	return u.Store.GetAlbum(id)
}

func (u *AlbumsService) DeleteAlbum(ctx context.Context, id string) error {
	return u.Store.DeleteAlbum(id)
}

func (u *AlbumsService) CreateAlbumFromUpload(ctx context.Context, userID, uploadKey, albumName string, isPublic, isFriends bool, friendIDs []string) (*t.Album, error) {
	now := time.Now()
	name := albumName
	if name == "" {
		name = uploadKey
	}
	return u.Store.CreateAlbum(&t.Album{
		Name:                name,
		UploadKey:           uploadKey,
		UserID:              userID,
		CreatedAt:           now,
		UpdatedAt:           now,
		Status:              enums.AlbumStatusProcessing,
		IsPublic:            isPublic,
		IsViewableByFriends: isFriends,
		FriendIDs:           friendIDs,
	})
}

func (u *AlbumsService) CreateAlbum(ctx context.Context, data *t.Album) (*t.Album, error) {
	now := time.Now()
	data.CreatedAt = now
	data.UpdatedAt = now
	return u.Store.CreateAlbum(data)
}

func (u *AlbumsService) UpdateAlbum(ctx context.Context, data *t.UpdateAlbum) (*t.Album, error) {
	data.UpdatedAt = time.Now()
	return u.Store.UpdateAlbum(data)
}

func (u *AlbumsService) AddVideosToAlbum(ctx context.Context, userID, uploadKey string, swingVideos []*t.SwingVideo) (*t.Album, error) {
	return u.Store.AddVideosToAlbum(userID, uploadKey, swingVideos)
}

func (u *AlbumsService) PostCommentToAlbum(ctx context.Context, albumID string, comment *t.Comment) (*t.Album, error) {
	return u.Store.PostCommentToAlbum(albumID, comment)
}

func (u *AlbumsService) PostCommentToSwing(ctx context.Context, albumID, swingID string, comment *t.Comment) (*t.Album, error) {
	return u.Store.PostCommentToSwing(albumID, swingID, comment)
}

func (u *AlbumsService) RecentAlbums(ctx context.Context, start, end time.Time, limit, offset int) ([]*t.Album, error) {
	return u.Store.RecentAlbums(start, end, limit, offset)
}

func (u *AlbumsService) RecentAlbumComments(ctx context.Context, start, end time.Time, limit, offset int) ([]*t.Comment, error) {
	return u.Store.RecentAlbumComments(start, end, limit, offset)
}

func (u *AlbumsService) RecentSwingComments(ctx context.Context, start, end time.Time, limit, offset int) ([]*t.Comment, error) {
	return u.Store.RecentSwingComments(start, end, limit, offset)
}
