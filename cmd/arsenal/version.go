package main

import "runtime/debug"

func buildVersion() string {
	moduleVersion := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		moduleVersion = info.Main.Version
	}
	return selectBuildVersion(version, moduleVersion)
}

func selectBuildVersion(injected, moduleVersion string) string {
	if injected != "" && injected != "dev" {
		return injected
	}
	if moduleVersion != "" && moduleVersion != "(devel)" {
		return moduleVersion
	}
	return "dev"
}
