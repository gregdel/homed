package components

// Config represents a component configuration
type Config struct {
	Type         string `yaml:"type"`
	FriendlyName string `yaml:"friendly_name"`
	StateTopic   string `yaml:"state_topic"`
	CommandTopic string `yaml:"command_topic"`
	Internal     bool   `yaml:"internal"`
}
