package store

import (
	"fmt"

	uuid "github.com/satori/go.uuid"

	m "github.com/tennis-community-api-service/pkg/mongo"
	t "github.com/tennis-community-api-service/users/types"
)

func (s *UsersStore) SendFriendRequest(req *t.FriendRequest, note *t.FriendNote) (err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	if req.ID == "" {
		req.ID = uuid.NewV4().String()
	}
	if note.ID == "" {
		note.ID = uuid.NewV4().String()
	}

	query := m.M{"_id": m.M{"$in": []string{req.FromUserID, req.ToUserID}}}
	err = m.UpdateAll(c, query, m.M{"$push": m.M{"frndReqs": m.M{
		"$each": []*t.FriendRequest{req},
		"$sort": m.M{"crAt": -1},
	}}})
	if err != nil {
		return
	}

	return m.Update(c, nil, m.M{"_id": req.ToUserID}, m.M{"$push": m.M{"frndNotes": note}})
}

func (s *UsersStore) AcceptFriendRequest(acceptorID, reqID string, accept bool, note *t.FriendNote) (user *t.User, err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	req := &t.FriendRequest{}
	if err = m.AggregateOne(c, req, []m.M{
		{"$match": m.M{"_id": acceptorID}},
		{"$unwind": "$frndReqs"},
		{"$match": m.M{"frndReqs._id": reqID}},
		{"$replaceRoot": m.M{"newRoot": "$frndReqs"}},
	}); err != nil {
		return
	}

	user = &t.User{}
	acceptorUpdate := m.M{"$pull": m.M{"frndReqs": m.M{"_id": reqID}}}
	targetUpdate := m.M{"$pull": m.M{"frndReqs": m.M{"_id": reqID}}}

	if accept {
		if note.ID == "" {
			note.ID = uuid.NewV4().String()
		}
		acceptorUpdate["$addToSet"] = m.M{"friendIds": req.FromUserID}
		targetUpdate["$addToSet"] = m.M{"friendIds": req.ToUserID}
		targetUpdate["$push"] = m.M{"frndNotes": note}
	}

	if err = m.Update(c, user, m.M{"_id": acceptorID}, acceptorUpdate); err != nil {
		return
	}
	err = m.Update(c, nil, m.M{"_id": req.FromUserID}, targetUpdate)
	return
}

func (s *UsersStore) Unfriend(sourceID, targetID string) (err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	ids := []string{sourceID, targetID}
	return m.UpdateAll(c, m.M{"_id": m.M{"$in": ids}}, m.M{"$pull": m.M{"friendIds": m.M{"$in": ids}}})
}

func (s *UsersStore) SearchFriends(search *string, IDs *[]string, limit, offset int) (friends []*t.Friend, err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	friends = []*t.Friend{}
	query := m.M{"status": "Created"}

	if search != nil {
		query["$or"] = []m.M{
			{"usrNm": m.M{"$regex": fmt.Sprintf("(?i)%s", *search)}},
			{"fnm": m.M{"$regex": fmt.Sprintf("(?i)%s", *search)}},
		}
	}

	if IDs != nil {
		query["_id"] = m.M{"$in": IDs}
	}

	if limit > 0 {
		err = c.Find(query).Skip(offset).Limit(limit).All(&friends)
	} else {
		err = m.Find(c, &friends, query, nil)
	}
	return
}
