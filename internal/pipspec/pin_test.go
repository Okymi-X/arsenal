package pipspec

import "testing"

func TestParseExact(t *testing.T) {
	tests := []struct {
		value   string
		wantErr bool
	}{
		{value: "dploot==3.1.3"},
		{value: "package-name==1!2.0.0+local.1"},
		{value: "", wantErr: true},
		{value: "-package==1.0", wantErr: true},
		{value: "dploot", wantErr: true},
		{value: "dploot>=3.1.0", wantErr: true},
		{value: "--index-url==https://example.test", wantErr: true},
		{value: "dploot==3.1.3; python_version > '3'", wantErr: true},
		{value: " dploot==3.1.3", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			_, err := ParseExact(test.value)
			if (err != nil) != test.wantErr {
				t.Fatalf("ParseExact(%q) error = %v, wantErr %v", test.value, err, test.wantErr)
			}
		})
	}
}

func TestPinNormalizedName(t *testing.T) {
	pin, err := ParseExact("My.Package_name==1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if got := pin.NormalizedName(); got != "my-package-name" {
		t.Fatalf("NormalizedName = %q", got)
	}
}
