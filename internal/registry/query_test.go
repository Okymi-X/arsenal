package registry

import "testing"

func fixture() *Registry {
	return &Registry{
		Tools: []Tool{
			{
				Name:        "netexec",
				Aliases:     []string{"nxc"},
				Description: "Network execution tool",
				Category:    CategoryAD,
				Tags:        []string{"smb", "ldap"},
				Versions:    []Version{{Tag: "1.4.0"}},
			},
			{
				Name:        "ffuf",
				Description: "Fast web fuzzer",
				Category:    CategoryWeb,
				Tags:        []string{"fuzzing"},
				Versions:    []Version{{Tag: "2.1.0"}},
			},
		},
	}
}

func TestFindTool(t *testing.T) {
	r := fixture()
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{"by name", "netexec", true},
		{"by alias", "nxc", true},
		{"case insensitive", "NetExec", true},
		{"missing", "ghost", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := r.FindTool(tt.query)
			if ok != tt.want {
				t.Fatalf("FindTool(%q) ok = %v, want %v", tt.query, ok, tt.want)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	r := fixture()
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"by tag", "smb", []string{"netexec"}},
		{"by category", "ad", []string{"netexec"}},
		{"by description", "fuzzer", []string{"ffuf"}},
		{"empty returns all sorted", "", []string{"ffuf", "netexec"}},
		{"no match", "zzz", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.Search(tt.query)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d results, want %d", len(got), len(tt.want))
			}
			for i, name := range tt.want {
				if got[i].Name != name {
					t.Fatalf("result %d = %s, want %s", i, got[i].Name, name)
				}
			}
		})
	}
}

func TestByCategory(t *testing.T) {
	r := fixture()
	ad := r.ByCategory(CategoryAD)
	if len(ad) != 1 || ad[0].Name != "netexec" {
		t.Fatalf("ByCategory(ad) = %+v", ad)
	}
}
