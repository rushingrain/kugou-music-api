package sdk

import "testing"

func TestDecodeRegisterDevResponsePlainJSON(t *testing.T) {
	decoded := decodeRegisterDevResponse([]byte(`{"data":{"scheme":0,"dfid":"349A4N1sIipW2kDpcc0ysWeR"},"status":1,"error_code":0}`), "unused")
	if decoded == nil {
		t.Fatal("decodeRegisterDevResponse() returned nil")
	}
	if decoded["status"] != float64(1) {
		t.Fatalf("status = %v, want 1", decoded["status"])
	}
	data, _ := decoded["data"].(map[string]any)
	if data["dfid"] != "349A4N1sIipW2kDpcc0ysWeR" {
		t.Fatalf("dfid = %v, want expected value", data["dfid"])
	}
}
