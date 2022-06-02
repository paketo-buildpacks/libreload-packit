package watchexec

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/paketo-buildpacks/libreload-packit"
	"github.com/paketo-buildpacks/packit/v2"
)

type WatchexecReloader struct{}

func NewWatchexecReloader() WatchexecReloader {
	return WatchexecReloader{}
}

func (WatchexecReloader) ShouldEnableLiveReload() (bool, error) {
	if reload, found := os.LookupEnv(libreload.LiveReloadEnabledEnvVar); found {
		if shouldEnableReload, err := strconv.ParseBool(reload); err != nil {
			return false, fmt.Errorf("failed to parse %s value %s: %w", libreload.LiveReloadEnabledEnvVar, reload, err)
		} else if shouldEnableReload {
			return true, nil
		}
	}
	return false, nil
}

func (WatchexecReloader) TransformReloadableProcesses(originalProcess packit.Process, spec libreload.ReloadableProcessSpec) (nonReloadable packit.Process, reloadable packit.Process) {
	nonReloadable = originalProcess
	nonReloadable.Default = false

	reloadable = originalProcess
	reloadable.Type = fmt.Sprintf("reload-%s", originalProcess.Type)
	reloadable.Command = "watchexec"
	reloadable.Args = buildArgs(originalProcess, spec)

	return nonReloadable, reloadable
}

// libreload.ReloadableProcessSpec contains information to be translated directly into watchexec arguments
// - spec.WatchPaths will translate into --watch args
// Optional. If len == 0, then no --watch args will be provided
// - spec.IgnorePaths will translate into --ignore args
// Optional. If len == 0, then no --ignore args will be provided
// - spec.Shell will translate into --shell
// Optional. If not provided, will use "none"
// - spec.VerbosityLevel will translate into -v
// If spec.VerbosityLevel is 0, no -v arg will be provided.
// If spec.VerbosityLevel is greater than 0, the appropriate number of v's will be added
// E.g. spec.VerbosityLevel = 2 => -vv
// E.g. spec.VerbosityLevel = 4 => -vvvv
func buildArgs(originalProcess packit.Process, spec libreload.ReloadableProcessSpec) []string {
	args := []string{
		"--restart",
	}

	for _, watchPath := range spec.WatchPaths {
		args = append(args, "--watch", watchPath)
	}

	for _, ignorePath := range spec.IgnorePaths {
		args = append(args, "--ignore", ignorePath)
	}

	shell := "none"
	if spec.Shell != "" {
		shell = spec.Shell
	}
	args = append(args, "--shell", shell)

	if spec.VerbosityLevel > 0 {
		args = append(args, "-"+strings.Repeat("v", spec.VerbosityLevel))
	}

	args = append(args,
		"--",
		originalProcess.Command)
	args = append(args, originalProcess.Args...)
	return args
}
