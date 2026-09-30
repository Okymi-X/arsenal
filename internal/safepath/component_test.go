package safepath

import "testing"

func TestValidateComponent(t *testing.T) {
	tests := []struct {
		value   string
		wantErr bool
	}{
		{value: "redteam-q3"},
		{value: "tool_1.2+test"},
		{value: "", wantErr: true},
		{value: ".", wantErr: true},
		{value: "..", wantErr: true},
		{value: "../escape", wantErr: true},
		{value: "nested/name", wantErr: true},
		{value: "name with space", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			err := ValidateComponent("name", test.value)
			if (err != nil) != test.wantErr {
				t.Fatalf("ValidateComponent(%q) error = %v, wantErr %v", test.value, err, test.wantErr)
			}
		})
	}
}
