package storage

import (
	"context"


	"gopkg.in/mgo.v2/bson"

	"github.com/evgeny08/collection-key/types"
)

// CreateUser creates a user in storage
func (s *Storage) InsertKey(ctx context.Context, key *types.Key) error {
	_, err := s.session.Collection(collectionKey).InsertOne(context.TODO(), &key)
	return err
}

// GetKey returns an unreleased key
func (s *Storage) GetKey(ctx context.Context) (*types.Key, error) {
	var key *types.Key
	filter := bson.M{"issued": false}
	err := s.session.Collection(collectionKey).FindOne(context.TODO(), filter).Decode(&key)
	if err != nil {
		return nil, err
	}

	_, err = s.session.Collection(collectionKey).UpdateOne(context.TODO(),bson.M{"id": key.ID}, bson.M{"$set": bson.M{"issued": true}})
	if err != nil {
		return nil, err
	}
	key.Issued = true

	return key, nil
}

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
