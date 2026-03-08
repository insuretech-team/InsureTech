package db

// MoneySerializer handles *commonv1.Money ↔ BIGINT (amount in paisa).
// The DB stores only the integer amount; currency lives in a separate column.
// When GORM reads a row, it passes the raw int64 from total_payable to Scan().
// We reconstruct a Money with amount only; the caller sets Currency separately.
//
// Usage in struct tag: gorm:"column:total_payable;not null;serializer:proto_money"
//
// NOTE: Because currency is a separate column, we only manage the amount here.
// The repository manually sets .Currency after a Find/First.

import (
	"context"
	"fmt"
	"reflect"

	commonv1 "github.com/newage-saint/insuretech/gen/go/insuretech/common/v1"
	"gorm.io/gorm/schema"
)

// ProtoMoneySerializer converts *commonv1.Money ↔ int64 (paisa).
type ProtoMoneySerializer struct{}

func (ProtoMoneySerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {
	if dbValue == nil {
		return nil
	}

	var amount int64
	switch v := dbValue.(type) {
	case int64:
		amount = v
	case int32:
		amount = int64(v)
	case float64:
		amount = int64(v)
	case []byte:
		// pgx may return numeric as []byte
		var n int64
		if _, err := fmt.Sscanf(string(v), "%d", &n); err != nil {
			return fmt.Errorf("proto_money: parse []byte %q: %w", v, err)
		}
		amount = n
	default:
		return fmt.Errorf("proto_money: unexpected db type %T", dbValue)
	}

	money := &commonv1.Money{
		Amount:        amount,
		Currency:      "BDT",
		DecimalAmount: float64(amount) / 100.0,
	}
	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(money))
	return nil
}

func (ProtoMoneySerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if fieldValue == nil {
		return int64(0), nil
	}
	m, ok := fieldValue.(*commonv1.Money)
	if !ok {
		return nil, fmt.Errorf("proto_money: expected *commonv1.Money, got %T", fieldValue)
	}
	if m == nil {
		return int64(0), nil
	}
	return m.Amount, nil
}

func init() {
	schema.RegisterSerializer("proto_money", ProtoMoneySerializer{})
}
