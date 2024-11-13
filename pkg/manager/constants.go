package manager

type Provider string

const (
	ProviderGithub           Provider = "github"
	ProviderGithubCRI        Provider = "github_cri"
	ProviderGithubWindowsCRI Provider = "github_windows_cri"
)
