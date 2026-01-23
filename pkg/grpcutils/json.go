package grpcutils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Uint64FromStringDecoder struct{}

func (d *Uint64FromStringDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	// В unsafe.Pointer мы не лезем, просто используем API итератора
	// для чтения значения. Он вернет его как interface{}
	value := iter.Read()
	switch v := value.(type) {
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			iter.Error = err
			return
		}

		*(*uint64)(ptr) = parsed
	case float64: // jsoniter читает числа как float64
		*(*uint64)(ptr) = uint64(v)
	case json.Number:
		parsed, err := v.Int64()
		if err != nil {
			iter.Error = err
			return
		}

		*(*uint64)(ptr) = uint64(parsed)
	default:
		iter.Error = fmt.Errorf("cannot decode %T into uint64", v)
	}
}

type TimestampDecoder struct{}

func (d *TimestampDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	// В unsafe.Pointer мы не лезем, просто используем API итератора
	// для чтения значения. Он вернет его как interface{}
	value := iter.Read()

	var ts *timestamppb.Timestamp

	switch v := value.(type) {
	case string:
		// Пробуем распарсить как RFC3339
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			// Если не получилось, пробуем как Unix timestamp в строке
			unixSec, parseErr := strconv.ParseInt(v, 10, 64)
			if parseErr != nil {
				iter.Error = fmt.Errorf("cannot parse timestamp string: %w", err)
				return
			}

			ts = timestamppb.New(time.Unix(unixSec, 0))
		} else {
			ts = timestamppb.New(t)
		}
	case float64: // jsoniter читает числа как float64
		// Интерпретируем как Unix timestamp в секундах
		ts = timestamppb.New(time.Unix(int64(v), 0))
	case json.Number:
		unixSec, err := v.Int64()
		if err != nil {
			iter.Error = err
			return
		}

		ts = timestamppb.New(time.Unix(unixSec, 0))
	case nil:
		// Если значение nil, оставляем указатель nil
		ts = nil
	default:
		iter.Error = fmt.Errorf("cannot decode %T into *timestamppb.Timestamp", v)
		return
	}

	*(**timestamppb.Timestamp)(ptr) = ts
}

func UnmarshalJSON(data []byte, v interface{}) error {
	jsoniter.RegisterTypeDecoder("uint64", &Uint64FromStringDecoder{})
	jsoniter.RegisterTypeDecoder("*timestamppb.Timestamp", &TimestampDecoder{})

	var json = jsoniter.ConfigDefault

	return json.Unmarshal(data, v)
}

func MarshalJSON(v interface{}) ([]byte, error) {
	jsoniter.RegisterTypeDecoder("uint64", &Uint64FromStringDecoder{})
	jsoniter.RegisterTypeDecoder("*timestamppb.Timestamp", &TimestampDecoder{})

	var json = jsoniter.ConfigDefault

	return json.Marshal(v)
}
