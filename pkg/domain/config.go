package domain

// Config contains application configuration
type Config struct {
	// Reserved for setting up themes/tones
	UIColor  string `json:"uiColor"`
	LogPath  string `json:"logPath"`
	// Specify log verbosity
	LogLevel string `json:"logLevel"`
	// Enable debugging the UI
	DebugUI bool `json:"debugUI"`
	// Static state configurations by us
	DefaultStateConfigs string `json:"defaultStateConfigs"`
	// State configurations set by the user should be saved here
	CustomStateConfigs string `json:"customStateConfigs"`
}