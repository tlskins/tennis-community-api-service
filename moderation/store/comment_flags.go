package store

import (
	"time"

	uuid "github.com/satori/go.uuid"

	t "github.com/tennis-community-api-service/moderation/types"
	m "github.com/tennis-community-api-service/pkg/mongo"
)

func (s *ModerationStore) CreateCommentFlag(data *t.CommentFlag) (flag *t.CommentFlag, err error) {
	sess, c := s.C(ColFlaggedComments)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	flag = &t.CommentFlag{}
	err = m.Upsert(c, flag, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *ModerationStore) UpdateCommentFlag(data *t.UpdateCommentFlag) (flag *t.CommentFlag, err error) {
	sess, c := s.C(ColFlaggedComments)
	defer sess.Close()

	flag = &t.CommentFlag{}
	err = m.Update(c, flag, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *ModerationStore) RecentFlaggedComments(start, end time.Time, resolved *bool, limit, offset int) (flags []*t.CommentFlag, err error) {
	sess, c := s.C(ColFlaggedComments)
	defer sess.Close()

	flags = []*t.CommentFlag{}
	query := m.M{"crAt": m.M{"$gte": start, "$lt": end}}

	if resolved != nil {
		query["res"] = *resolved
	}

	if limit > 0 {
		err = c.Find(query).Skip(offset).Limit(limit).All(&flags)
	} else {
		err = m.Find(c, &flags, query, nil)
	}
	return
}
