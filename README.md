# Paketo libreload

This library exists to easily support live reloadable processes for Paketo-style buildpacks built using [`packit v2`](https://github.com/paketo-buildpacks/packit).

## Detect

Use `ShouldEnableLiveReload` to decide when to require the necessary dependency.

```go
if shouldEnableReload, err := reload.ShouldEnableLiveReload(); err != nil {
    return packit.DetectResult{}, err
} else if shouldEnableReload {
    requirements = append(requirements, <require your preferred reload implementation such as watchexec>)
}
```

## Build

Use `TransformReloadableProcesses` to transform the non-reloadable process into a reloadable and non-reloadable process.

```go
if shouldEnableReload, err := reload.ShouldEnableLiveReload(); err != nil {
    return packit.BuildResult{}, err
} else if shouldEnableReload {
    nonReloadableProcess, reloadableProcess := reload.TransformReloadableProcesses(originalProcess, libreload.ReloadableProcessSpec{
        WatchPaths: []string{...},
        Shell: "none",
    })
    processes = append(processes, nonReloadableProcess, reloadableProcess)
} else {
    processes = append(processes, originalProcess)
}
```