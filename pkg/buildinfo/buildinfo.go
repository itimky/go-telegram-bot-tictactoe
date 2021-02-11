package buildinfo

import "time"

// These variables will be filled while compiling the app. Check Makefile.
var (
	appName     = unknownValue
	appVersion  = unknownValue
	buildTime   = unknownValue
	buildNumber = unknownValue
	gitHash     = unknownValue
	gitBranch   = unknownValue
)

const (
	buildTimeFormat = time.RFC3339
	unknownValue    = "unknown"
)

func AppName() string {
	return appName
}

func AppVersion() string {
	return appVersion
}

var startTime = time.Now()

func Uptime() string {
	return time.Since(startTime).String()
}

func BuildTime() time.Time {
	if buildTime == unknownValue {
		return startTime
	}

	t, err := time.Parse(buildTimeFormat, buildTime)
	if err != nil {
		panic("wrong format of buildTime: " + err.Error())
	}

	return t
}

func BuildNumber() string {
	return buildNumber
}

func GitHash() string {
	return gitHash
}

func GitBranch() string {
	return gitBranch
}
