package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm/schema"
)

// ProtoTimestampSerializer handles conversion between timestamppb.Timestamp and DB timestamp
type ProtoTimestampSerializer struct{}

// Scan implements the schema.Serializer interface
func (ProtoTimestampSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	fieldValue := reflect.New(field.FieldType)

	if dbValue != nil {
		var t time.Time
		switch v := dbValue.(type) {
		case time.Time:
			t = v
		case []byte:
			// Handle byte conversion if needed (though GORM usually gives time.Time for Postgres)
			// This might vary by driver
			return fmt.Errorf("byte scan not implemented for proto timestamp")
		default:
			return fmt.Errorf("failed to scan proto timestamp: unexpected value type %T", dbValue)
		}

		if !t.IsZero() {
			pbTs := timestamppb.New(t)
			// dst is the field value (pointer to *timestamppb.Timestamp)
			// field.FieldType is *timestamppb.Timestamp
			// We need to set dst to pbTs

			// Handle direct assignment
			if !fieldValue.Elem().CanSet() {
				// If we can't set, it might be weird refelction state, but usually:
				// dst.Set(reflect.ValueOf(pbTs))
			}
			// We scan INTO the destination
			// dst is a pointer to the field.
			// field is *timestamppb.Timestamp

			// Simple assignment:
			// *dst = pbTs
			field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(pbTs))
			return nil
		}
	}
	// If nil, set nil
	return nil
}

// Value implements the schema.Serializer interface
func (ProtoTimestampSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if fieldValue == nil {
		return nil, nil
	}

	ts, ok := fieldValue.(*timestamppb.Timestamp)
	if !ok {
		return nil, fmt.Errorf("expected *timestamppb.Timestamp, got %T", fieldValue)
	}

	if ts == nil {
		return nil, nil
	}

	return ts.AsTime(), nil
}

// RegisterSerializer registers the serializer with GORM
// Note: It seems GORM registers serializers by name in the tag, but we need to ensure the schema knows it.
// Actually, GORM looks up serializer by name from the `gorm:"serializer:name"` tag.
// We manually register it in an init function for global availability.

func init() {
	schema.RegisterSerializer("proto_timestamp", ProtoTimestampSerializer{})
}
