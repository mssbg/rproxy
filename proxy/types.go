package proxy

const (
	STRATEGY_ROUND_ROBIN = "round-robin"
	STRATEGY_RANDOM      = "random"
	STRATEGY_VAR_NAME    = "RPROXY_STRATEGY"
)

type Host struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Service struct {
	Name     string     `yaml:"name"`
	Domain   string     `yaml:"domain"`
	Hosts    []Host     `yaml:"hosts"`
	NextHost <-chan int `yaml:"-"`
}

type Proxy struct {
	Strategy   string              `yaml:"strategy"`
	Listen     Host                `yaml:"listen"`
	Services   []Service           `yaml:"services"`
	ServiceMap map[string]*Service `yaml:"-"`
}

type Config struct {
	Proxy Proxy `yaml:"proxy"`
}
