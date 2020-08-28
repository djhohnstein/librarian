package osinfo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"syscall"
)

type Release struct {
	Major int
	Minor int
}
type OSInfo struct {
	Sysname    string
	Nodename   string
	Release    Release
	Version    string
	Machine    string
	DomainName string
}

// A utility to convert the values to proper strings.
func int8ToStr(arr []int8) string {
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

func New() (*OSInfo, error) {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		return nil, err
	}
	strRelease := int8ToStr(uname.Release[:])
	partsRelease := strings.Split(strRelease, ".")
	if len(partsRelease) < 2 {
		return nil, errors.New(fmt.Sprintf("Unexpected OS Release: %s\n", strRelease))
	}
	major, err := strconv.Atoi(partsRelease[0])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(partsRelease[1])
	if err != nil {
		return nil, err
	}
	return &OSInfo{
		Sysname:    int8ToStr(uname.Sysname[:]),
		Nodename:   int8ToStr(uname.Nodename[:]),
		Release:    Release{Major: major, Minor: minor},
		Version:    int8ToStr(uname.Version[:]),
		Machine:    int8ToStr(uname.Machine[:]),
		DomainName: int8ToStr(uname.Domainname[:]),
	}, nil
}
