package version

import (
	"strings"
	"testing"
)

func Test(t *testing.T) {
	expected := "applicationName ApplicationID:  version: buildVersion " +
		"Build:buildTag(commitHash) - buildDate Runtime: buildOS/buildArch " +
		"GoVersion: goVersion"
	testinfo := Info{
		ApplicationName: "applicationName",
		CommitHash:      "commitHash",
		BuildDate:       "buildDate",
		BuildTag:        "buildTag",
		BuildOS:         "buildOS",
		BuildArch:       "buildArch",
		BuildVersion:    "buildVersion",
		GoVersion:       "goVersion",
	}
	result := testinfo.String()
	if !strings.EqualFold(expected, result) {
		t.Errorf("expected:\"%s\" result:\"%s\"", expected, result)
	}
}
