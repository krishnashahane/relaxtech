package output

import "testing"

func TestSelectColumns(t *testing.T) {
	data := []map[string]any{{"x": 1, "y": 2}, {"x": 3, "y": 4}}
	result := SelectColumns(data, []string{"y"})
	if len(result) != 2 {
		t.Fatalf("expected 2 records, got %d", len(result))
	}
	if _, present := result[0]["x"]; present {
		t.Fatalf("column x should have been excluded")
	}
	if result[0]["y"].(int) != 2 || result[1]["y"].(int) != 4 {
		t.Fatalf("unexpected column values: %+v", result)
	}
}
