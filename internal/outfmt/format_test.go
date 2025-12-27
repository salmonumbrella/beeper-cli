package outfmt

import (
	"bytes"
	"testing"
)

func TestTextTable(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTableWriter(&buf)
	tw.SetHeader([]string{"Name", "Age"})
	tw.Append([]string{"Alice", "30"})
	tw.Append([]string{"Bob", "25"})
	tw.Render()

	output := buf.String()
	if output == "" {
		t.Error("TextTable output is empty")
	}
	if !bytes.Contains(buf.Bytes(), []byte("Alice")) {
		t.Error("Output should contain 'Alice'")
	}
}

func TestJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"name": "test"}

	if err := WriteJSON(&buf, data); err != nil {
		t.Fatalf("WriteJSON() error: %v", err)
	}

	expected := `{"name":"test"}`
	got := bytes.TrimSpace(buf.Bytes())
	if string(got) != expected {
		t.Errorf("WriteJSON() = %q, want %q", got, expected)
	}
}
