package store

import (
	"time"

	uuid "github.com/satori/go.uuid"

	t "github.com/tennis-community-api-service/albums/types"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

func (s *AlbumsStore) UpdateSwing(albumID string, data *t.UpdateSwingVideo) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	album = &t.Album{}
	err = m.Update(c, album, m.M{"_id": albumID, "swingVids._id": data.ID}, m.M{"$set": m.M{"swingVids.$": data}})
	return
}

func (s *AlbumsStore) PostCommentToSwing(albumID, swingID string, comment *t.Comment) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	if comment.ID == "" {
		comment.ID = uuid.NewV4().String()
	}

	album = &t.Album{}
	err = m.Update(c, album, m.M{"_id": albumID, "swingVids._id": swingID}, m.M{
		"$set":  m.M{"updAt": time.Now()},
		"$push": m.M{"swingVids.$.cmnts": comment},
	})
	return
}

func (s *AlbumsStore) RecentSwingComments(start, end time.Time, limit, offset int) (comments []*t.Comment, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	comments = []*t.Comment{}

	m.Aggregate(c, &comments, []m.M{
		{"$match": m.M{"updAt": m.M{"$gte": start}}},
		{"$unwind": "$swingVids"},
		{"$unwind": "$swingVids.cmnts"},
		{"$match": m.M{"swingVids.cmnts.crAt": m.M{"$gte": start, "$lt": end}}},
		{"$sort": m.M{"swingVids.cmnts.crAt": -1}},
		{"$skip": offset},
		{"$limit": limit},
		{"$addFields": m.M{"swingVids.cmnts.albumId": "$_id", "swingVids.cmnts.swingId": "$swingVids._id"}},
		{"$replaceRoot": m.M{"newRoot": "$swingVids.cmnts"}},
	})
	return
}
