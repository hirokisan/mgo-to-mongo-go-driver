package bson

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CustomRegistry :
func CustomRegistry() *bsoncodec.RegistryBuilder {
	customValues := []interface{}{
		PtrObjectID(ObjectID("")),
		M{},
	}
	rb := bson.NewRegistryBuilder()
	for _, v := range customValues {
		t := reflect.TypeOf(v)
		defDecoder, err := bson.DefaultRegistry.LookupDecoder(t)
		if err != nil {
			panic(err)
		}
		rb.RegisterTypeDecoder(t, &customDecoder{
			defDecoder: defDecoder,
			zeroValue:  reflect.Zero(t)})
	}
	return rb
}

// customDecoder :
type customDecoder struct {
	defDecoder bsoncodec.ValueDecoder
	zeroValue  reflect.Value
}

// DecodeValue :
func (d *customDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if val.Type() == reflect.TypeOf(primitive.M{}) {
		if err := bsoncodec.NewMapCodec().DecodeValue(dctx, vr, val); err != nil {
			return err
		}
		for _, key := range val.MapKeys() {
			if val.MapIndex(key).IsNil() {
				continue
			}
			if reflect.ValueOf(val.MapIndex(key).Interface()).Type() != reflect.TypeOf(primitive.ObjectID{}) {
				continue
			}
			rawID, ok := val.MapIndex(key).Interface().(primitive.ObjectID)
			if !ok {
				continue
			}
			id := ObjectIDHex(rawID.Hex())
			val.SetMapIndex(key, reflect.ValueOf(&id).Elem())
		}
		return nil
	}

	if vr.Type() != bsontype.Null {
		return d.defDecoder.DecodeValue(dctx, vr, val)
	}

	if !val.CanSet() {
		return errors.New("value not settable")
	}
	if err := vr.ReadNull(); err != nil {
		return err
	}

	val.Set(d.zeroValue)

	return nil
}
