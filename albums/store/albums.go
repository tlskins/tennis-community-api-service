package store

import (
	"time"

	uuid "github.com/satori/go.uuid"

	t "github.com/tennis-community-api-service/albums/types"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

func (s *AlbumsStore) GetAlbum(id string) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	album = &t.Album{}
	err = m.FindOne(c, album, m.M{"_id": id})
	return
}

func (s *AlbumsStore) DeleteAlbum(id string) (err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	return m.Remove(c, m.M{"_id": id})
}

func (s *AlbumsStore) SearchAlbums(userIDs, friendIDs []string, public, friends, homeApproved *bool, limit, offset int) (albums []*t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	query := m.M{}
	albums = []*t.Album{}

	if len(userIDs) > 0 {
		query["userId"] = m.M{"$in": userIDs}
	}
	if len(friendIDs) > 0 {
		query["frndIds"] = m.M{"$elemMatch": m.M{"$in": friendIDs}}
	}
	if public != nil {
		query["public"] = *public
	}
	if friends != nil {
		query["frndView"] = *friends
	}
	if homeApproved != nil {
		query["home"] = *homeApproved
	}

	if limit > 0 {
		err = c.Find(query).Skip(offset).Sort("-updAt").Limit(limit).All(&albums)
	} else {
		err = m.Find(c, &albums, query, nil, []string{"-updAt"})
	}
	return
}

func (s *AlbumsStore) CreateAlbum(data *t.Album) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	album = &t.Album{}
	err = m.Upsert(c, album, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *AlbumsStore) UpdateAlbum(data *t.UpdateAlbum) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	album = &t.Album{}
	err = m.Update(c, album, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *AlbumsStore) AddVideosToAlbum(userId, uploadKey string, swings []*t.SwingVideo) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	for _, swing := range swings {
		if swing.ID == "" {
			swing.ID = uuid.NewV4().String()
		}
	}

	album = &t.Album{}
	err = m.Update(c, album, m.M{"userId": userId, "upKey": uploadKey}, m.M{
		"$set":  m.M{"updAt": time.Now()},
		"$push": m.M{"swingVids": m.M{"$each": swings}},
	})
	return
}

func (s *AlbumsStore) PostCommentToAlbum(albumID string, comment *t.Comment) (album *t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	if comment.ID == "" {
		comment.ID = uuid.NewV4().String()
	}

	album = &t.Album{}
	err = m.Update(c, album, m.M{"_id": albumID}, m.M{
		"$set":  m.M{"updAt": time.Now()},
		"$push": m.M{"cmnts": comment},
	})
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

func (s *AlbumsStore) RecentAlbums(start, end time.Time, limit, offset int) (albums []*t.Album, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	albums = []*t.Album{}
	query := m.M{"crAt": m.M{"$gte": start, "$lt": end}}

	if limit > 0 {
		err = c.Find(query).Sort("-crAt").Skip(offset).Limit(limit).All(&albums)
	} else {
		err = m.Find(c, &albums, query, nil)
	}
	return
}

func (s *AlbumsStore) RecentAlbumComments(start, end time.Time, limit, offset int) (comments []*t.Comment, err error) {
	sess, c := s.C(ColAlbums)
	defer sess.Close()

	comments = []*t.Comment{}

	m.Aggregate(c, &comments, []m.M{
		{"$match": m.M{"updAt": m.M{"$gte": start}}},
		{"$unwind": "$cmnts"},
		{"$match": m.M{"cmnts.crAt": m.M{"$gte": start, "$lt": end}}},
		{"$sort": m.M{"cmnts.crAt": -1}},
		{"$skip": offset},
		{"$limit": limit},
		{"$addFields": m.M{"cmnts.albumId": "$_id"}},
		{"$replaceRoot": m.M{"newRoot": "$cmnts"}},
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
