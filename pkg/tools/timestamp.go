package tools

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoTimestamp converts time.Time to *timestamppb.Timestamp.
// Returns nil if the input time is zero.
func ToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// ToProtoTimestampPtr converts *time.Time to *timestamppb.Timestamp.
// Returns nil if the input pointer is nil.
func ToProtoTimestampPtr(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return ToProtoTimestamp(*t)
}

// ToTime converts *timestamppb.Timestamp to time.Time.
// Returns zero time if the input is nil.
func ToTime(t *timestamppb.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.AsTime()
}

// ToTimePtr converts *timestamppb.Timestamp to *time.Time.
// Returns nil if the input is nil.
func ToTimePtr(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}
	res := t.AsTime()
	return &res
}
