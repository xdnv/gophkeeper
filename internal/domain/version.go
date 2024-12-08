// common configuration module filling info version info from runtime debug data
package domain

import "runtime/debug"

type Version struct {
	Version    string
	Commit     string
	LongCommit string
	Date       string
}

func GetVersion() *Version {
	av := new(Version)

	if info, ok := debug.ReadBuildInfo(); ok {
		av.Version = info.Main.Version

		for _, setting := range info.Settings {
			//vcs.revision: the revision identifier for the current commit or checkout
			if setting.Key == "vcs.revision" {
				av.Commit = setting.Value
				if len(av.Commit) > 7 {
					av.Commit = av.Commit[:7]
				}
			}
			//vcs.time: the modification time associated with vcs.revision, in RFC3339 format
			if setting.Key == "vcs.time" {
				av.Date = setting.Value
			}
		}

	}

	return av
}
