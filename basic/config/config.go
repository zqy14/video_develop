package config

type Config struct {
	System System `mapstructure:"system" json:"system"`
}

type System struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     uint64 `mapstructure:"port" json:"port"`
	Database string `mapstructure:"database" json:"database"`
}

type NaCos struct {
	NamespaceId string `mapstructure:"namespace_id" json:"namespace_id"`
	IpAddr      string `mapstructure:"ip_addr" json:"ip_addr"`
	Port        uint64 `mapstructure:"port" json:"port"`
	DataId      string `mapstructure:"data_id" json:"data_id"`
	Group       string `mapstructure:"group" json:"group"`
}
