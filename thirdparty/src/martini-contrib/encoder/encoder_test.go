package encoder

import (
	"encoding/json"
	"testing"
)

type Sample struct {
	Visible string `json:"visible"`
	Hidden  string `json:"hidden" out:"false"`
}

func TestEncoder(t *testing.T) {
	src := &Sample{Visible: "visible", Hidden: "this field won't be exported"}
	dst := &Sample{}

	enc := &JsonEncoder{}
	result, err := enc.Encode(src)
	if err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(result, dst); err != nil {
		t.Fatal("Unmarshal error:", err)
	}

	if dst.Hidden != "" {
		t.Fatalf("Expected empty field 'Hidden', got %v\n", dst.Hidden)
	}
}
