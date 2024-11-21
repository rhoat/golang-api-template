package version

import (
	"runtime"
	"strings"

	"github.com/rhoat/go-exercise/pkg/system"
)

type Info struct {
	ApplicationName string
	ApplicationID   string
	CommitHash      string
	BuildDate       string
	BuildTag        string
	BuildOS         string
	BuildArch       string
	BuildVersion    string
	GoVersion       string
}

func NewInfo() Info {
	return Info{
		ApplicationName: system.ApplicationName,
		ApplicationID:   system.ApplicationID,
		BuildVersion:    system.BuildVersion,
		BuildTag:        system.BuildTag,
		CommitHash:      system.CommitHash,
		BuildDate:       system.BuildDate,
		BuildOS:         runtime.GOOS,
		BuildArch:       runtime.GOARCH,
		GoVersion:       runtime.Version(),
	}
}

func (i Info) String() string {
	var builder strings.Builder

	builder.WriteString(i.ApplicationName)
	builder.WriteString(" ApplicationID: ")
	builder.WriteString(i.ApplicationID)
	builder.WriteString(" version: ")
	builder.WriteString(i.BuildVersion)
	builder.WriteString(" Build:")
	builder.WriteString(i.BuildTag)
	builder.WriteString("(")
	builder.WriteString(i.CommitHash)
	builder.WriteString(") - ")
	builder.WriteString(i.BuildDate)
	builder.WriteString(" Runtime: ")
	builder.WriteString(i.BuildOS)
	builder.WriteString("/")
	builder.WriteString(i.BuildArch)
	builder.WriteString(" GoVersion: ")
	builder.WriteString(i.GoVersion)

	return builder.String()
}
