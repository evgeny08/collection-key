package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
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

	_, err = s.session.Collection(collectionKey).UpdateOne(context.TODO(), bson.M{"issued": key.Issued}, bson.M{"$set": bson.M{"issued": true}})
	if err != nil {
		return nil, err
	}
	key.Issued = true

	return key, nil
}

// CanceledKey updates key Redemption with given id
func (s *Storage) CanceledKey(ctx context.Context, id string) error {
	var key *types.Key
	err := s.session.Collection(collectionKey).FindOne(context.TODO(), bson.M{"id": id}).Decode(&key)
	if err != nil {
		return err
	}
	if !key.Issued {
		return errors.New("the key was not issued")
	}
	if key.Canceled {
		return errors.New("the key has already been canceled")
	}
	_, err = s.session.Collection(collectionKey).UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": bson.M{"canceled": true}})
	if err != nil {
		return err
	}
	return nil
}

// VerificationKey return key info
func (s *Storage) VerificationKey(ctx context.Context, id string) (*types.Key, error) {
	var key *types.Key
	err := s.session.Collection(collectionKey).FindOne(context.TODO(), bson.M{"id": id}).Decode(&key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// UnreleasedKey return list unreleased key
func (s *Storage) UnreleasedKey(ctx context.Context) ([]*types.Key, error) {
	filter := bson.M{"issued": false}
	cursor, err := s.session.Collection(collectionKey).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var listKey []*types.Key
	for cursor.Next(context.TODO()) {
		var key *types.Key
		err := cursor.Decode(&key)
		if err != nil {
			return nil, err
		}
		listKey = append(listKey, key)
	}
	if len(listKey) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return listKey, nil
}
