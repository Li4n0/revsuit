package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/li4n0/revsuit/internal/update"
	"github.com/li4n0/revsuit/pkg/server"
	"github.com/pkg/errors"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	log "unknwon.dev/clog/v2"
)

func checkUpdateFromCli() {
	log.Info("Checking for updates...")
	if upgradeable, release, err := update.CheckUpgrade(server.VERSION); err != nil {
		log.Warn(err.Error())
	} else if upgradeable {
		err := updateFromCli(release)
		if err != nil {
			log.Warn(err.Error())
		}
	} else {
		log.Info("The current version is the latest")
	}
}

func updateFromCli(release *selfupdate.Release) error {
	releaseNotes := release.ReleaseNotes
	releaseNotes = strings.Replace(releaseNotes, "\n\n", "\n", -1)
	releaseNotes = strings.Replace(releaseNotes, "\r\n\r\n", "\n", -1)
	releaseNotes = strings.Replace(releaseNotes, "**", "", -1)

	log.Info("New version detected: %s\n%s", release.URL, releaseNotes)
	log.Warn("Do you want to update to v%s?[Y/n]:", release.Version)

	var upgrade string
	_, _ = fmt.Scanln(&upgrade)
	if upgrade != "n" {
		exe, err := os.Executable()
		if err != nil {
			log.Error("error when locate executable path, err: %v", err)
		}

		assetURL := "https://upgrade.revsuit.pro" + strings.TrimPrefix(release.AssetURL, "https://github.com/Li4n0/revsuit/releases") + "?from=" + server.VERSION

		log.Info("Downloading from %s", assetURL)
		if err := selfupdate.UpdateTo(assetURL, exe); err != nil {
			return errors.Wrap(err, "error when download asset")
		}
		log.Info("Successfully updated to v%s", release.Version)
		return nil
	}
	return nil
}
