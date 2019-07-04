package storage

import (
	"context"

	//"gopkg.in/mgo.v2/bson"

	"github.com/evgeny08/collection-key/types"
)

// CreateUser creates a user in storage
func (s *Storage) InsertKey(ctx context.Context, key *types.Key) error {
	_, err := s.session.Collection(collectionKey).InsertOne(context.TODO(), &key)
	return err
}

//// FindUserByLogin find user by given login
//func (s *Storage) FindUserByLogin(ctx context.Context, login string) (*types.User, error) {
//	var user *types.User
//	filter := bson.M{"login": login}
//	err := s.session.Collection(collectionUser).FindOne(context.TODO(), filter).Decode(&user)
//	return user, err
//}
//
//// CreateSession create new session in storage
//func (s *Storage) CreateSession(ctx context.Context, session *types.Session) error {
//	_, err := s.session.Collection(collectionAuth).InsertOne(context.TODO(), &session)
//	return err
//}
//
//// FindAccessToken find AccessToken by storage
//func (s *Storage) FindAccessToken(ctx context.Context, clientToken string) (*types.Session, error) {
//	var session *types.Session
//	filter := bson.M{"access_token": clientToken}
//	err := s.session.Collection(collectionAuth).FindOne(context.TODO(), filter).Decode(&session)
//	return session, err
//}
