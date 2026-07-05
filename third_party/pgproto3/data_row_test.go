package pgproto3

import "testing"

func TestDataRowDecodeRejectsNegativeFieldLength(t *testing.T) {
	// DataRow with one field and an invalid negative field length other than -1 (NULL).
	src := []byte{0, 1, 0xff, 0xff, 0xff, 0xfe}

	var row DataRow
	if err := row.Decode(src); err == nil {
		t.Fatal("expected invalid message error")
	}
}
