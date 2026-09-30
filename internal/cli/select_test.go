package cli

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/Okymi-X/arsenal/internal/registry"
)

func TestSelectOne(t *testing.T) {
	var output bytes.Buffer
	got, err := selectOne(bufio.NewReader(strings.NewReader("2\n")), &output, "Choose:", []selectOption{
		{label: "one", value: "first"},
		{label: "two", value: "second"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != "second" {
		t.Fatalf("selectOne() = %q; want second", got)
	}
	if !strings.Contains(output.String(), "Select [1-2]:") {
		t.Fatalf("selection prompt missing from %q", output.String())
	}
}

func TestCuratedGitHubVersionMatchesCommonTagForms(t *testing.T) {
	tool := registry.Tool{Versions: []registry.Version{
		{Tag: "1.5.1", Commit: "v1.5.1", Tested: true},
	}}
	for _, ref := range []string{"1.5.1", "v1.5.1"} {
		version, ok := curatedGitHubVersion(tool, ref)
		if !ok || !version.Tested {
			t.Fatalf("curatedGitHubVersion(%q) = %#v, %v", ref, version, ok)
		}
	}
	if _, ok := curatedGitHubVersion(tool, "v1.5.0"); ok {
		t.Fatal("unexpected match for uncurated ref")
	}
}

func TestVerifySelectedGitHubCommitRejectsMovedTag(t *testing.T) {
	expected := "0123456789abcdef0123456789abcdef01234567"
	resolved := "fedcba9876543210fedcba9876543210fedcba98"
	if err := verifySelectedGitHubCommit("v1.0.0", expected, expected); err != nil {
		t.Fatalf("unchanged tag rejected: %v", err)
	}
	err := verifySelectedGitHubCommit("v1.0.0", expected, resolved)
	if err == nil || !strings.Contains(err.Error(), "moved from 0123456789ab to fedcba987654") {
		t.Fatalf("moved tag error = %v", err)
	}
}

func TestSelectOneRejectsInvalidChoice(t *testing.T) {
	_, err := selectOne(bufio.NewReader(strings.NewReader("3\n")), &bytes.Buffer{}, "Choose:", []selectOption{
		{label: "one", value: "first"},
		{label: "two", value: "second"},
	})
	if err == nil || !strings.Contains(err.Error(), "invalid selection") {
		t.Fatalf("selectOne() error = %v; want invalid selection", err)
	}
}

func TestSelectOneSanitizesDisplayedMetadata(t *testing.T) {
	var output bytes.Buffer
	_, err := selectOne(bufio.NewReader(strings.NewReader("1\n")), &output, "Choose:\x1b", []selectOption{
		{label: "unsafe\nlabel", value: "safe"},
		{label: "other", value: "other"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "\x1b") || strings.Contains(output.String(), "unsafe\nlabel") {
		t.Fatalf("control character was not sanitized in %q", output.String())
	}
}

func TestSelectOneReturnsOnlyChoiceWithoutReading(t *testing.T) {
	got, err := selectOne(nil, nil, "Choose:", []selectOption{{label: "one", value: "first"}})
	if err != nil || got != "first" {
		t.Fatalf("selectOne() = %q, %v; want first", got, err)
	}
}

func TestSelectOnePreservesFollowingSelection(t *testing.T) {
	input := bufio.NewReader(strings.NewReader("2\n1\n"))
	options := []selectOption{{label: "one", value: "first"}, {label: "two", value: "second"}}
	first, err := selectOne(input, &bytes.Buffer{}, "Choose:", options)
	if err != nil {
		t.Fatal(err)
	}
	second, err := selectOne(input, &bytes.Buffer{}, "Choose:", options)
	if err != nil {
		t.Fatal(err)
	}
	if first != "second" || second != "first" {
		t.Fatalf("selections = %q, %q; want second, first", first, second)
	}
}
