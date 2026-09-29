package registry

// VersionRelation describes an installed version relative to the registry's
// newest tested release. Registry order is authoritative and newest-first.
type VersionRelation uint8

const (
	VersionUntracked VersionRelation = iota
	VersionCurrent
	VersionOutdated
	VersionAhead
	VersionNoRecommendation
)

// RecommendedVersion returns the newest version marked tested.
func (t Tool) RecommendedVersion() (Version, bool) {
	for _, version := range t.Versions {
		if version.Tested {
			return version, true
		}
	}
	return Version{}, false
}

// DefaultVersion returns the recommended version, or the newest catalogued
// version when maintainers have not marked any release tested.
func (t Tool) DefaultVersion() (Version, bool) {
	if version, ok := t.RecommendedVersion(); ok {
		return version, true
	}
	if len(t.Versions) == 0 {
		return Version{}, false
	}
	return t.Versions[0], true
}

// CompareToRecommended classifies current without parsing version strings.
// This supports date, commit, and upstream-specific tags consistently.
func (t Tool) CompareToRecommended(current string) (Version, VersionRelation) {
	recommended, ok := t.RecommendedVersion()
	if !ok {
		return Version{}, VersionNoRecommendation
	}
	recommendedIndex := -1
	currentIndex := -1
	for i, version := range t.Versions {
		if version.Tag == recommended.Tag {
			recommendedIndex = i
		}
		if version.Tag == current {
			currentIndex = i
		}
	}
	if currentIndex < 0 {
		return recommended, VersionUntracked
	}
	switch {
	case currentIndex == recommendedIndex:
		return recommended, VersionCurrent
	case currentIndex > recommendedIndex:
		return recommended, VersionOutdated
	default:
		return recommended, VersionAhead
	}
}
