package version

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"
	"strings"
	"time"
)

const tsFormat = "2006-01-02_15:04:05"

var (
	release   = "-"
	revision  = "-"
	buildUser = "-"
	buildDate = "-"
	goVersion = runtime.Version()
)

// Release is the semantic version of the current build.
func Release() string {
	return release
}

// Revision is the last git commit hash of the source repository at the moment
// the binary was built.
func Revision() string {
	return revision
}

// BuildUser is the username of the user who performed the build.
func BuildUser() string {
	return buildUser
}

// BuildDate is a timestamp of the moment when the binary was built.
func BuildDate() (time.Time, error) {
	return time.ParseInLocation(tsFormat, buildDate, time.UTC)
}

// BuildDateString is BuildDate formatted as RFC3339
func BuildDateString() string {
	ts, err := BuildDate()
	if err != nil {
		return buildDate
	}
	return ts.UTC().Format(time.RFC3339)
}

// GoVersion is the go version the build utilizes.
func GoVersion() string {
	return goVersion
}

// String returns the version and build information.
func String() string {
	return fmt.Sprintf("release=%s, revision=%s, user=%s, build=%s, go=%s",
		release, revision, buildUser, buildDate, goVersion)
}

// versionInfoTmpl contains the template used by Info.
var versionInfoTmpl = `
{{.program}}, release {{.release}} (revision: {{.revision}})
  build user:       {{.buildUser}}
  build date:       {{.buildDate}}
  go version:       {{.goVersion}}
`

// Print returns version and build information in a user friendly, formatted
// string.
func Print(program string) string {
	m := map[string]string{
		"program":   program,
		"release":   release,
		"revision":  revision,
		"buildUser": buildUser,
		"buildDate": buildDate,
		"goVersion": goVersion,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		return ""
	}
	return strings.TrimSpace(buf.String())
}
