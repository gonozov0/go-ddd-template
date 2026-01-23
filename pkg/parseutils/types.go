package parseutils

import (
	"encoding/json"
	"time"

	gofrs "github.com/gofrs/uuid"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"fmt"
)

var moscowLocation *time.Location

const readableTimeFormat = "02.01.2006 15:04"

func init() {
	var err error

	moscowLocation, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic("failed to load timezone Europe/Moscow: " + err.Error())
	}
}

func GRPCTimestampToDate(timestamp *timestamppb.Timestamp) time.Time {
	if timestamp == nil || !timestamp.IsValid() {
		return time.Time{}
	}

	return timestamp.AsTime()
}

func GRPCTimestampToDateWithTruncateMinute(timestamp *timestamppb.Timestamp) time.Time {
	return GRPCTimestampToDate(timestamp).Truncate(time.Minute)
}

func DateToGRPCTimestamp(t time.Time) *timestamppb.Timestamp {
	if !t.IsZero() {
		return timestamppb.New(t)
	}

	return nil
}

func Int32ToIntPointer(p *int32) *int {
	if p == nil {
		return nil
	}

	v := int(*p)

	return &v
}

func JsonToGRPCStruct(jsonData json.RawMessage) (*structpb.Struct, error) {
	rpcData := &structpb.Struct{}

	if jsonData != nil && string(jsonData) != "null" {
		err := protojson.Unmarshal(jsonData, rpcData)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal protojson data: %w", err)
		}
	}

	return rpcData, nil
}

func GRPCListToStrings(rpcData *structpb.ListValue) ([]string, error) {
	result := make([]string, 0, len(rpcData.Values))

	for _, value := range rpcData.GetValues() {
		result = append(result, value.GetStringValue())
	}

	return result, nil
}

func GoogleUUIDsToGofrsUUIDs(ids uuid.UUIDs) []gofrs.UUID {
	gofrsUUIDs := make([]gofrs.UUID, 0, len(ids))

	for _, id := range ids {
		gofrsUUIDs = append(gofrsUUIDs, gofrs.UUID(id))
	}

	return gofrsUUIDs
}

func GofrsUUIDsToGoogleUUIDs(ids []gofrs.UUID) uuid.UUIDs {
	uuids := make(uuid.UUIDs, 0, len(ids))

	for _, id := range ids {
		uuids = append(uuids, uuid.UUID(id))
	}

	return uuids
}

func ConvertTimeToMskString(utcTime time.Time) string {
	moscowTime := utcTime.In(moscowLocation)
	formattedTime := moscowTime.Format(readableTimeFormat)

	return formattedTime
}
