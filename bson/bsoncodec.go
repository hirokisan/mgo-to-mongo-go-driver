package bson

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// CustomRegistry :
func CustomRegistry() *bsoncodec.RegistryBuilder {
	customValues := []interface{}{
		&ObjectID{},
	}
	rb := bson.NewRegistryBuilder()
	for _, v := range customValues {
		t := reflect.TypeOf(v)
		defDecoder, err := bson.DefaultRegistry.LookupDecoder(t)
		if err != nil {
			panic(err)
		}
		rb.RegisterTypeDecoder(t, &nullawareDecoder{
			defDecoder: defDecoder,
			zeroValue:  reflect.Zero(t)})
	}
	return rb
}

// nullawareDecoder : ObjectIDのnullをnilとして扱うためのdecoder
type nullawareDecoder struct {
	defDecoder bsoncodec.ValueDecoder
	zeroValue  reflect.Value
}

// DecodeValue :
func (d *nullawareDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
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
