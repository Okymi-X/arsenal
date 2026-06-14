package store

// Find returns the installed tool matching name and version.
func (m *Manifest) Find(name, version string) (InstalledTool, bool) {
	for _, t := range m.Tools {
		if t.Name == name && t.Version == version {
			return t, true
		}
	}
	return InstalledTool{}, false
}

// Versions returns every installed version of the named tool.
func (m *Manifest) Versions(name string) []InstalledTool {
	var out []InstalledTool
	for _, t := range m.Tools {
		if t.Name == name {
			out = append(out, t)
		}
	}
	return out
}

// Active returns the active installation of the named tool.
func (m *Manifest) Active(name string) (InstalledTool, bool) {
	for _, t := range m.Tools {
		if t.Name == name && t.Active {
			return t, true
		}
	}
	return InstalledTool{}, false
}

// Upsert inserts or replaces an installation, preserving uniqueness by
// name and version.
func (m *Manifest) Upsert(t InstalledTool) {
	for i := range m.Tools {
		if m.Tools[i].Name == t.Name && m.Tools[i].Version == t.Version {
			m.Tools[i] = t
			return
		}
	}
	m.Tools = append(m.Tools, t)
}

// Delete removes an installation by name and version, reporting success.
func (m *Manifest) Delete(name, version string) bool {
	for i := range m.Tools {
		if m.Tools[i].Name == name && m.Tools[i].Version == version {
			m.Tools = append(m.Tools[:i], m.Tools[i+1:]...)
			return true
		}
	}
	return false
}

// SetActive marks one version active and clears the flag on the tool's
// other versions. It reports whether the target was found.
func (m *Manifest) SetActive(name, version string) bool {
	found := false
	for i := range m.Tools {
		if m.Tools[i].Name != name {
			continue
		}
		active := m.Tools[i].Version == version
		m.Tools[i].Active = active
		if active {
			found = true
		}
	}
	return found
}
