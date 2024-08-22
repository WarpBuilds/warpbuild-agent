# Proxy

Runs a local server on WarpBuild Runners. Serves the following use-cases:

- /\_apis/artifactcache: A local cache proxy server that is able to take requests from Buildkit's Remote Cache GHA backend and proxy them to work with WarpCache.
