package store

import (
	m "github.com/tennis-community-api-service/pkg/mongo"
	t "github.com/tennis-community-api-service/users/types"

	uuid "github.com/satori/go.uuid"
)

func (s *UsersStore) GetConfirmation(ID string) (conf *t.UserConfirmation, err error) {
	sess, c := s.C(ColConfirmations)
	defer sess.Close()

	conf = &t.UserConfirmation{}
	err = m.FindOne(c, conf, m.M{"_id": ID})
	return
}

func (s *UsersStore) CreateConfirmation(data *t.UserConfirmation) (conf *t.UserConfirmation, err error) {
	sess, c := s.C(ColConfirmations)
	defer sess.Close()

	if data.ID == "" {
		data.ID = uuid.NewV4().String()
	}
	conf = &t.UserConfirmation{}
	err = m.Upsert(c, conf, m.M{"_id": data.ID}, m.M{"$set": data})
	return
}

func (s *UsersStore) DeleteUserConfirmations(userID, email string) (err error) {
	sess, c := s.C(ColConfirmations)
	defer sess.Close()

	or := []m.M{}
	if userID != "" {
		or = append(or, m.M{"usrId": userID})
	}
	if email != "" {
		or = append(or, m.M{"em": email})
	}

	return m.Remove(c, m.M{"$or": or})
}
