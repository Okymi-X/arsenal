package fetcher

import "testing"

func TestNormalizeName(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Rubeus.exe", "rubeus"},
		{"linpeas.sh", "linpeas"},
		{"winPEAS.bat", "winpeas"},
		{"pspy64", "pspy64"},
	}
	for _, tt := range tests {
		if got := normalizeName(tt.in); got != tt.want {
			t.Errorf("normalizeName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestMatchName(t *testing.T) {
	sharp := []string{"Rubeus.exe", "Seatbelt.exe", "Certify.exe"}
	pspy := []string{"pspy32", "pspy32s", "pspy64", "pspy64s"}
	tests := []struct {
		name       string
		want       string
		candidates []string
		out        string
		ok         bool
	}{
		{"loose case and suffix", "rubeus", sharp, "Rubeus.exe", true},
		{"exact wins over substring", "pspy64", pspy, "pspy64", true},
		{"substring fallback", "seat", sharp, "Seatbelt.exe", true},
		{"no match", "mimikatz", sharp, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := matchName(tt.want, tt.candidates)
			if got != tt.out || ok != tt.ok {
				t.Fatalf("matchName(%q) = (%q, %v), want (%q, %v)", tt.want, got, ok, tt.out, tt.ok)
			}
		})
	}
}
