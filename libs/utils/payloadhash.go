package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// ComputePayloadHash returns a short hash string for the provided request payload.
// For proto.Message it uses deterministic proto marshaling. For other structs it falls back to JSON.
func ComputePayloadHash(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("nil payload")
	}
	if m, ok := v.(proto.Message); ok {
		b, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
		if err != nil {
			return "", fmt.Errorf("marshal proto: %w", err)
		}
		h := sha256.Sum256(b)
		return hex.EncodeToString(h[:8]), nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:8]), nil
}
