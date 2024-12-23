package helper

import (
	"encoding/json"
	"fmt"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// DateStrToUnixNano convert date string formated to epoch unix nano seconds.
func DateStrToUnixNano(value string, layout ...string) int64 {
	format := "2006-01-02"
	if len(layout) > 0 {
		format = layout[0]
	}

	t, err := time.Parse(format, value)
	if err != nil {
		return 0
	}

	return t.UnixNano()
}

// DateStrToUnixNanoStrict same like DateStrToUnixNano, but will return error if got any error.
func DateStrToUnixNanoStrict(value string, layout ...string) (int64, error) {
	format := "2006-01-02"
	if len(layout) > 0 {
		format = layout[0]
	}

	t, err := time.Parse(format, value)
	if err != nil {
		return 0, err
	}

	return t.UnixNano(), nil
}

// NewPBStruct constructs a Struct from a general-purpose Go map.
// The map keys must be valid UTF-8.
// The map values are converted using NewValue.
func NewPBStruct(v map[string]interface{}) (map[string]*structpb.Value, error) {
	x := make(map[string]*structpb.Value, len(v))
	for k, v := range v {
		if !utf8.ValidString(k) {
			return nil, fmt.Errorf("invalid UTF-8 in string: %q", k)
		}
		var err error
		x[k], err = structpb.NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

// convert *structpb.Struct to map[string]interface{}
func PBStructToMap(s *structpb.Struct) (*string, error) {
	b, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	str := fmt.Sprintf("%s", m["opt"])
	return &str, nil
}

// convert map[string]interface{} to *structpb.Struct
func MapToProtobufStruct(m map[string]interface{}) (*structpb.Struct, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	s := &structpb.Struct{}
	err = protojson.Unmarshal(b, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
