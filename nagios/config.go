package nagios

import "time"

const (
	GitlabGitProvider = "gitlab"
	GithubGitProvider = "github"
)

type ProbeConfig struct {
	APIEndpoint             string        `mapstructure:"api-endpoint"`
	Debug                   bool          `mapstructure:"debug"`
	GitProvider             string        `mapstructure:"git-provider"`
	APIToken                string        `mapstructure:"api-token"`
	Project                 string        `mapstructure:"project"`
	Timeout                 time.Duration `mapstructure:"timeout"`
	TargetBranch            string        `mapstructure:"target-branch"`
	WarningLastUpdateDelay  time.Duration `mapstructure:"delay-warning-last-update"`
	CriticalLastUpdateDelay time.Duration `mapstructure:"delay-critical-last-update"`
}
