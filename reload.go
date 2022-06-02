package libreload

import "github.com/paketo-buildpacks/packit/v2"

type ReloadableProcessSpec struct {
	WatchPaths     []string
	IgnorePaths    []string
	Shell          string
	VerbosityLevel int
}

var LiveReloadEnabledEnvVar = "BP_LIVE_RELOAD_ENABLED"

type Reloader interface {
	ShouldEnableLiveReload() (bool, error)
	TransformReloadableProcesses(originalProcess packit.Process, spec ReloadableProcessSpec) (nonReloadable packit.Process, reloadable packit.Process)
}
