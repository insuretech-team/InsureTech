package db

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"gorm.io/gorm/schema"
)

// ProtoEnumSerializer handles conversion between proto enums and DB string/varchar columns.
//
// In this codebase, DB enum columns often store values without the enum type prefix, e.g.:
//   DB:  "B2C_CUSTOMER"
//   Proto enum String(): "USER_TYPE_B2C_CUSTOMER"
//
// This serializer supports both forms for reading, and writes the DB-friendly form by default.
// It is intentionally generic and does not hardcode enum names.
type ProtoEnumSerializer struct{}

func (ProtoEnumSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	if dbValue == nil {
		return nil
	}

	var raw string
	switch v := dbValue.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("failed to scan proto enum: unexpected value type %T", dbValue)
	}

	normDB := normalizeEnumToken(raw)
	dstType := field.FieldType

	// Proto enums are int32 aliases and support String() on the value receiver.
	// We resolve by brute force over a bounded range.
	for i := int32(0); i < 256; i++ {
		cand := reflect.ValueOf(i).Convert(dstType)
		m := cand.MethodByName("String")
		if !m.IsValid() {
			return fmt.Errorf("proto_enum serializer used on type %s without String()", dstType.Name())
		}
		out := m.Call(nil)
		if len(out) == 0 {
			continue
		}

		protoStr := out[0].String() // e.g. USER_TYPE_B2C_CUSTOMER
		normProto := normalizeEnumToken(protoStr)

		// Match either exact or suffix form.
		if normProto == normDB || strings.HasSuffix(normProto, "_"+normDB) {
			field.ReflectValueOf(ctx, dst).Set(cand)
			return nil
		}
	}

	return fmt.Errorf("failed to map string %q to enum %s", raw, dstType.Name())
}

func (ProtoEnumSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if fieldValue == nil {
		return nil, nil
	}

	v := reflect.ValueOf(fieldValue)
	m := v.MethodByName("String")
	if !m.IsValid() {
		return nil, fmt.Errorf("proto_enum serializer used on value %T without String()", fieldValue)
	}
	out := m.Call(nil)
	if len(out) == 0 {
		return nil, fmt.Errorf("failed to call String() on enum value %v", fieldValue)
	}

	protoStr := out[0].String() // e.g. USER_TYPE_B2C_CUSTOMER
	prefix := enumPrefixFromTypeName(field.FieldType.Name())
	if strings.HasPrefix(protoStr, prefix) {
		return strings.TrimPrefix(protoStr, prefix), nil
	}

	// Fallback: if it looks prefixed (X_Y_), return token after the prefix-like part if possible.
	// Prefer DB style tokens where possible.
	return protoStr, nil
}

func normalizeEnumToken(s string) string {
	s = strings.TrimSpace(strings.ToUpper(s))
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

// enumPrefixFromTypeName converts a Go enum type name (CamelCase) into the typical
// generated enum string prefix (UPPER_SNAKE_) used by protoc-gen-go.
// Example: UserType -> USER_TYPE_
func enumPrefixFromTypeName(typeName string) string {
	if typeName == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(typeName) + 8)
	for i, r := range typeName {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r)
			continue
		}
		b.WriteRune(unicode.ToUpper(r))
	}
	b.WriteByte('_')
	return b.String()
}

func init() {
	schema.RegisterSerializer("proto_enum", ProtoEnumSerializer{})
}
