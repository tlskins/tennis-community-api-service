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

func (s *UsersStore) AddFriendNoteToUsers(userIDs []string, note *t.FriendNote) (err error) {
	sess, c := s.C(ColUsers)
	defer sess.Close()

	if note.ID == "" {
		note.ID = uuid.NewV4().String()
	}

	err = m.UpdateAll(c, m.M{"_id": m.M{"$in": userIDs}}, m.M{"$push": m.M{"frndNotes": m.M{
		"$each":  []*t.FriendNote{note},
		"$sort":  m.M{"crAt": -1},
		"$slice": 10,
	}}})
	return
}
