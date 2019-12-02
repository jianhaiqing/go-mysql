package canal

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pingcap/errors"
	"github.com/siddontang/go-mysql/mysql"
)

type DumpConfig struct {
	// mysqldump execution path, like mysqldump or /usr/bin/mysqldump, etc...
	// If not set, ignore using mysqldump.
	ExecutionPath string `toml:"mysqldump"`

	// Will override Databases, tables is in database table_db
	Tables  []string `toml:"tables"`
	TableDB string   `toml:"table_db"`

	Databases []string `toml:"dbs"`

	// Ignore table format is db.table
	IgnoreTables []string `toml:"ignore_tables"`

	// Dump only selected records. Quotes are mandatory
	Where string `toml:"where"`

	// If true, discard error msg, else, output to stderr
	DiscardErr bool `toml:"discard_err"`

	// Set true to skip --master-data if we have no privilege to do
	// 'FLUSH TABLES WITH READ LOCK'
	SkipMasterData bool `toml:"skip_master_data"`

	// Only MySQL's GTID is supported, cause mysqldump's parameter is a little different between MySQL and MariaDB
	// mysqldump of MySQL uses --set-gtid-purged, MariaDB uses --gtid instead.
	// GtidPurged:auto when MySQL supports GTID,and GtidPurged:none for gtid is disabled,or gtid is not supported,
	// or you strongly wants file-position replication.
	GtidPurged string `toml:"set_gtid_purged"`

	// Set to change the default max_allowed_packet size
	MaxAllowedPacketMB int `toml:"max_allowed_packet_mb"`

	// Set to change the default protocol to connect with
	Protocol string `toml:"protocol"`

	// Set extra options
	ExtraOptions []string `toml:"extra_options"`
}

type Config struct {
	Addr     string `toml:"addr"`
	User     string `toml:"user"`
	Password string `toml:"password"`

	Charset         string        `toml:"charset"`
	ServerID        uint32        `toml:"server_id"`
	Flavor          string        `toml:"flavor"`
	HeartbeatPeriod time.Duration `toml:"heartbeat_period"`
	ReadTimeout     time.Duration `toml:"read_timeout"`

	// IncludeTableRegex or ExcludeTableRegex should contain database name
	// Only a table which matches IncludeTableRegex and dismatches ExcludeTableRegex will be processed
	// eg, IncludeTableRegex : [".*\\.canal"], ExcludeTableRegex : ["mysql\\..*"]
	//     this will include all database's 'canal' table, except database 'mysql'
	// Default IncludeTableRegex and ExcludeTableRegex are empty, this will include all tables
	IncludeTableRegex []string `toml:"include_table_regex"`
	ExcludeTableRegex []string `toml:"exclude_table_regex"`

	// discard row event without table meta
	DiscardNoMetaRowEvent bool `toml:"discard_no_meta_row_event"`

	Dump DumpConfig `toml:"dump"`

	UseDecimal bool `toml:"use_decimal"`
	ParseTime  bool `toml:"parse_time"`

	TimestampStringLocation *time.Location

	// SemiSyncEnabled enables semi-sync or not.
	SemiSyncEnabled bool `toml:"semi_sync_enabled"`

	// Set to change the maximum number of attempts to re-establish a broken
	// connection
	MaxReconnectAttempts int `toml:"max_reconnect_attempts"`
}

func NewConfigWithFile(name string) (*Config, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return NewConfig(string(data))
}

func NewConfig(data string) (*Config, error) {
	var c Config

	_, err := toml.Decode(data, &c)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &c, nil
}

func NewDefaultConfig() *Config {
	c := new(Config)

	c.Addr = "127.0.0.1:3306"
	c.User = "root"
	c.Password = ""

	c.Charset = mysql.DEFAULT_CHARSET
	c.ServerID = uint32(rand.New(rand.NewSource(time.Now().Unix())).Intn(1000)) + 1001

	c.Flavor = "mysql"

	c.Dump.ExecutionPath = "mysqldump"
	c.Dump.DiscardErr = true
	c.Dump.SkipMasterData = false
	// add default value to disable mysqldump --set-gtid-purged
	c.Dump.GtidPurged = "none"

	return c
}
