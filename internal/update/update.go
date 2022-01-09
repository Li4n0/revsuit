package update

import (
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func CheckUpgrade(currentVersion string) (bool, *selfupdate.Release, error) {
	latest, found, err := selfupdate.DetectLatest("Li4n0/revsuit")
	if err != nil {
		return false, nil, errors.Wrap(err, "error when detecting version")
	}

	v, err := semver.Parse(strings.TrimPrefix(currentVersion, "v"))
	if err != nil {
		return false, nil, errors.Wrap(err, "error when parse version")
	}
	if !found || latest.Version.LTE(v) {
		return false, nil, nil
	}
	return true, latest, nil
}
