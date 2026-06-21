package manager

type Provider string

const (
	ProviderGithub           Provider = "github"
	ProviderGithubCRI        Provider = "github_cri"
	ProviderGithubWindowsCRI Provider = "github_windows_cri"
	// ProviderClaudeAgent is not a settings.json provider — it is the value the backend
	// returns in allocation_details.runner_application for a Claude managed-agent sandbox.
	ProviderClaudeAgent Provider = "claude_agent"
)
