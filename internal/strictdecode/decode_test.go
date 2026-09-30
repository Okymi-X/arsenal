package strictdecode

import "testing"

type document struct {
	Name string `json:"name" toml:"name"`
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{name: "valid", data: `{"name":"value"}`},
		{name: "unknown", data: `{"name":"value","extra":true}`, wantErr: true},
		{name: "trailing", data: `{"name":"value"} {"name":"other"}`, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got document
			err := JSON([]byte(test.data), &got)
			if (err != nil) != test.wantErr {
				t.Fatalf("JSON() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestTOML(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{name: "valid", data: `name = "value"`},
		{name: "unknown", data: "name = \"value\"\nextra = true", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got document
			err := TOML([]byte(test.data), &got)
			if (err != nil) != test.wantErr {
				t.Fatalf("TOML() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
