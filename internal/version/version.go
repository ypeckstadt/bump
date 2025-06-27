package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

func Parse(versionStr string) (*Version, error) {
	re := regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) != 4 {
		return nil, fmt.Errorf("invalid version format: %s", versionStr)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		Raw:   versionStr,
	}, nil
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) BumpPatch() *Version {
	return &Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch + 1,
	}
}

func (v *Version) BumpMinor() *Version {
	return &Version{
		Major: v.Major,
		Minor: v.Minor + 1,
		Patch: 0,
	}
}

func (v *Version) BumpMajor() *Version {
	return &Version{
		Major: v.Major + 1,
		Minor: 0,
		Patch: 0,
	}
}

func (v *Version) Bump(versionType string) (*Version, error) {
	switch strings.ToLower(versionType) {
	case "patch":
		return v.BumpPatch(), nil
	case "minor":
		return v.BumpMinor(), nil
	case "major":
		return v.BumpMajor(), nil
	default:
		return nil, fmt.Errorf("invalid version type: %s (must be patch, minor, or major)", versionType)
	}
}

func NewFromString(versionStr string) *Version {
	if versionStr == "" || versionStr == "v0.0.0" {
		return &Version{Major: 0, Minor: 0, Patch: 0}
	}

	version, err := Parse(versionStr)
	if err != nil {
		return &Version{Major: 0, Minor: 0, Patch: 0}
	}

	return version
}
