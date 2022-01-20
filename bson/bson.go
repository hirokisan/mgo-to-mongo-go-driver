package bson

import (
	"encoding/hex"
	"fmt"

	mgobson "github.com/hirokisan/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// M :
type M = primitive.M

// D :
type D = primitive.D

// RegEx :
type RegEx = primitive.Regex

// DateTime :
type DateTime = primitive.DateTime

// ObjectID :
type ObjectID primitive.ObjectID

// NewObjectID :
func NewObjectID() ObjectID {
	return ObjectID(primitive.NewObjectID())
}

// ObjectIDHex :
func ObjectIDHex(s string) ObjectID {
	pid, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		panic(fmt.Errorf("invalid input to ObjectIdHex: %q :%w", s, err))
	}
	return ObjectID(pid)
}

// IsObjectIDHex :
func IsObjectIDHex(s string) bool {
	return primitive.IsValidObjectID(s)
}

// Unmarshal :
func Unmarshal(data []byte, val interface{}) error {
	return bson.Unmarshal(data, val)
}

// Marshal :
func Marshal(val interface{}) ([]byte, error) {
	return bson.Marshal(val)
}

// MarshalBSONValue : write
func (id ObjectID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	b, err := hex.DecodeString((id).Hex())
	return bsontype.ObjectID, b, err
}

// UnmarshalBSONValue : read
func (id *ObjectID) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	// Set to nil by nullawareDecoder
	if t == bsontype.Null {
		return nil
	}
	if t == bsontype.EmbeddedDocument {
		var got struct {
			Data    *ObjectID `bson:"data"`
			SubType int       `bson:"sybtype"`
		}
		if err := bson.Unmarshal(data, &got); err != nil {
			return err
		}
		if got.Data == nil {
			return nil
		}
		*id = *got.Data
		return nil
	}
	hex := hex.EncodeToString(data)
	if !primitive.IsValidObjectID(hex) {
		return fmt.Errorf("invalid object id: %s", hex)
	}
	tmp, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return fmt.Errorf("invalid object id: %s", hex)
	}
	*id = ObjectID(tmp)
	return nil
}

// Hex :
func (id ObjectID) Hex() string {
	return hex.EncodeToString(id[:])
}

// GetBSON(mgo) : write
func (id ObjectID) GetBSON() (interface{}, error) {
	return mgobson.ObjectIdHex(id.Hex()), nil
}

// SetBSON(mgo) : read
func (id *ObjectID) SetBSON(raw mgobson.Raw) error {
	// ref : https://bsonspec.org/spec.html
	if raw.Kind == 10 && len(raw.Data) == 0 {
		return mgobson.SetZero
	}
	var v []byte
	if err := raw.Unmarshal(&v); err != nil {
		return fmt.Errorf("unmarshal bson: %w", err)
	}
	hex := hex.EncodeToString(v)
	if !primitive.IsValidObjectID(hex) {
		return fmt.Errorf("invalid object id: %s", hex)
	}
	tmp, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return fmt.Errorf("from hex: %s", hex)
	}
	*id = ObjectID(tmp)
	return nil
}
