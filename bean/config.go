package bean

// ServerConfig ServerConfig Struct
type ServerConfig struct {
	Port        []string
	Mode        string
	GormLogMode string
	CallBack    string
}

// EthereumConfig EthereumConfig
type EthereumConfig struct {
	ChainID    int64
	Host       string
	Address    string
	Passphrase string
}

// DBConfig DBConfig Struct
type DBConfig struct {
	Dialect      string
	Database     string
	User         string
	Password     string
	Host         string
	Port         int
	Charset      string
	URL          string
	MaxIdleConns int
	MaxOpenConns int
}

// NsqConfig NsqConfig
type NsqConfig struct {
	Host  string
	Topic string
}
