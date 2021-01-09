package store

import (
	uuid "github.com/satori/go.uuid"

	m "github.com/tennis-community-api-service/pkg/mongo"
	t "github.com/tennis-community-api-service/users/types"
)

func (s *UsersStore) GetUser(userID string) (user *t.User, err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	user = &t.User{}
	err = m.FindOne(c, user, m.M{"_id": userID})
	return
}

func (s *UsersStore) GetFriends(userIDs []string) (friends []*t.Friend, err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	friends = []*t.Friend{}
	err = m.Find(c, &friends, m.M{"_id": m.M{"$in": userIDs}})
	return
}

func (s *UsersStore) GetUserByEmail(email string) (user *t.User, err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	user = &t.User{}
	err = m.FindOne(c, user, m.M{"em": email})
	return
}

func (s *UsersStore) CreateUser(data *t.User) (user *t.User, err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	user = &t.User{}
	err = m.Upsert(c, user, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *UsersStore) UpdateUser(data *t.UpdateUser) (user *t.User, err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	user = &t.User{}
	err = m.Update(c, user, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *UsersStore) AddUploadNote(userID string, note *t.UploadNote) (user *t.User, err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	user = &t.User{}
	err = m.Update(c, user, m.M{"_id": userID}, m.M{"$push": m.M{"upNotes": m.M{
		"$each":  []*t.UploadNote{note},
		"$sort":  m.M{"crAt": -1},
		"$slice": 10,
	}}})
	return
}

func (s *UsersStore) SendFriendRequest(req *t.FriendRequest) (err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	if req.ID == "" {
		req.ID = uuid.NewV4().String()
	}

	query := m.M{"_id": m.M{"$in": []string{req.FromUserID, req.ToUserID}}}
	return m.UpdateAll(c, query, m.M{"$push": m.M{"frndReqs": m.M{
		"$each": []*t.FriendRequest{req},
		"$sort": m.M{"crAt": -1},
	}}})
}

func (s *UsersStore) AcceptFriendRequest(acceptorID, reqID string, accept bool) (user *t.User, err error) {
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
	acceptorUpdate := m.M{"$pull": m.M{"frndReqs": m.M{"$elemMatch": m.M{"_id": reqID}}}}
	targetUpdate := m.M{"$pull": m.M{"frndReqs": m.M{"$elemMatch": m.M{"_id": reqID}}}}

	if accept {
		acceptorUpdate["$push"] = m.M{"friendIds": req.FromUserID}
		targetUpdate["$push"] = m.M{"friendIds": req.ToUserID}
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
