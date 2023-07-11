package model

type Default struct {
	Addr string `mapstructure:"ADDRESS"`
	Port string `mapstructure:"PORT"`
	Debug bool `mapstructure:"DEBUG"`
	AddressRB string `mapstructure:"ADDRESSRB"`
	UserRB string `mapstructure:"USERRB"`
	PassRB string `mapstructure:"PASSRB"`
	PortRB int `mapstructure:"PORTRB"`
	Queue string `mapstructure:"QUEUE"`
}