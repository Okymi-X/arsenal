package registry

import (
	"strings"
	"testing"
)

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

func TestParseRejectsDuplicateVersionTags(t *testing.T) {
	data := []byte(`
version = "1"
[[tool]]
name = "duplicate"
install_method = "pip"
  [[tool.version]]
  tag = "1.0.0"
  [[tool.version]]
  tag = "1.0.0"
`)
	if _, err := Parse(data); err == nil {
		t.Fatal("expected duplicate version tags to be rejected")
	}
}

func TestParseRequiresTargetsForNativeInstallers(t *testing.T) {
	data := []byte(`
version = "1"
[[tool]]
name = "native"
install_method = "gobin"
binary = "native"
  [[tool.version]]
  tag = "1.0.0"
  commit = "v1.0.0"
`)
	if _, err := Parse(data); err == nil {
		t.Fatal("expected missing install target to be rejected")
	}
}

func TestParseAcceptsCompleteNativeTargets(t *testing.T) {
	data := []byte(`
version = "1"
[[tool]]
name = "native"
install_method = "gobin"
binary = "native"
  [[tool.version]]
  tag = "1.0.0"
  commit = "v1.0.0"
  install_targets = { native = "github.com/owner/native" }
`)
	if _, err := Parse(data); err != nil {
		t.Fatalf("Parse: %v", err)
	}
}

func TestParseRejectsUnknownFields(t *testing.T) {
	data := []byte(`
version = "1"
unknown = true
[[tool]]
name = "example"
install_method = "pip"
binary = "example"
  [[tool.version]]
  tag = "1.0.0"
`)
	if _, err := Parse(data); err == nil {
		t.Fatal("expected unknown registry field to be rejected")
	}
}

func TestParseRejectsInvalidInstallMethod(t *testing.T) {
	data := []byte(`
version = "1"
[[tool]]
name = "example"
install_method = "typo"
binary = "example"
  [[tool.version]]
  tag = "1.0.0"
`)
	if _, err := Parse(data); err == nil {
		t.Fatal("expected invalid install method to be rejected")
	}
}

func TestParseRejectsUnsafeAssetMetadata(t *testing.T) {
	tests := []string{
		`name = "../escape"`,
		"name = \"asset\"\nrepo = \"https://example.com/owner/repo\"",
		"name = \"asset\"\nrepo = \"https://github.com/owner/repo\"\ndir = \"../escape\"",
		"name = \"asset\"\nrepo = \"https://github.com/owner/repo\"\nbuilds = [\"../escape\"]",
	}
	for _, fields := range tests {
		data := []byte("version = \"1\"\n[[asset]]\n" + fields + "\nsource = \"github-raw\"\n")
		if _, err := Parse(data); err == nil {
			t.Fatalf("expected asset metadata to be rejected:\n%s", data)
		}
	}
}

func TestParseValidatesExactPipDependencies(t *testing.T) {
	valid := []byte(`
version = "1"
[[tool]]
name = "example"
install_method = "gitpip"
repo = "https://github.com/owner/example"
binary = "example"
  [[tool.version]]
  tag = "1.0.0"
  commit = "v1.0.0"
  pip_dependencies = ["dependency==2.0.0"]
`)
	if _, err := Parse(valid); err != nil {
		t.Fatalf("Parse valid dependency: %v", err)
	}
	duplicate := []byte(strings.ReplaceAll(string(valid), `pip_dependencies = ["dependency==2.0.0"]`, `pip_dependencies = ["dependency-name==2.0.0", "Dependency_name==2.0.0"]`))
	if _, err := Parse(duplicate); err == nil {
		t.Fatal("expected normalized duplicate dependency to be rejected")
	}

	for _, requirement := range []string{"dependency>=2", "--index-url==example", "dependency=="} {
		data := []byte(strings.ReplaceAll(string(valid), "dependency==2.0.0", requirement))
		if _, err := Parse(data); err == nil {
			t.Fatalf("expected dependency %q to be rejected", requirement)
		}
	}
}
