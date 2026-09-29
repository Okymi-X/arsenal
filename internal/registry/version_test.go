package registry

import "testing"

func TestCompareToRecommended(t *testing.T) {
	tool := Tool{Versions: []Version{
		{Tag: "next"},
		{Tag: "2.0.0", Tested: true},
		{Tag: "1.0.0", Tested: true},
	}}
	tests := []struct {
		current string
		want    VersionRelation
	}{
		{current: "next", want: VersionAhead},
		{current: "2.0.0", want: VersionCurrent},
		{current: "1.0.0", want: VersionOutdated},
		{current: "custom", want: VersionUntracked},
	}
	for _, test := range tests {
		gotVersion, gotRelation := tool.CompareToRecommended(test.current)
		if gotVersion.Tag != "2.0.0" || gotRelation != test.want {
			t.Fatalf("CompareToRecommended(%q) = %q, %v", test.current, gotVersion.Tag, gotRelation)
		}
	}
}

func TestCompareToRecommendedWithoutTestedVersion(t *testing.T) {
	tool := Tool{Versions: []Version{{Tag: "1.0.0"}}}
	if _, relation := tool.CompareToRecommended("1.0.0"); relation != VersionNoRecommendation {
		t.Fatalf("relation = %v, want VersionNoRecommendation", relation)
	}
}
