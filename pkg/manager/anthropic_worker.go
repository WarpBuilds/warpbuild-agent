package manager

import (
	"os"
	"path/filepath"
	"runtime"
)

// anthropicWorkerBinary is the Anthropic CLI that runs the self-hosted worker (`ant beta:worker run`).
// It is installed by cloud-init on claude_agent sandbox VMs (see backend-core's runner templates), NOT
// by the agent — the generic runner image ships Apache Ant at /usr/bin/ant, a different binary.
const anthropicWorkerBinary = "ant"

// anthropicWorkerPath returns the absolute path cloud-init installs the Anthropic `ant` CLI to. The
// worker is invoked by this absolute path (never via PATH/LookPath) so we never pick up Apache Ant at
// /usr/bin/ant.
func anthropicWorkerPath() string {
	if runtime.GOOS == "windows" {
		pd := os.Getenv("ProgramData")
		if pd == "" {
			pd = `C:\ProgramData`
		}
		return filepath.Join(pd, "warpbuild", "bin", "ant.exe")
	}
	return "/usr/local/bin/ant"
}
