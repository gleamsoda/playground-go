package config

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBHost               string        `mapstructure:"DB_HOST"`
	DBPort               int           `mapstructure:"DB_PORT"`
	DBUser               string        `mapstructure:"DB_USER"`
	DBPassword           string        `mapstructure:"DB_PASSWORD"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
}

// Get reads configuration from file or environment variables.
var Get = sync.OnceValue(func() (cfg Config) {
	bindEnv(cfg)
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	return
})

func (c Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c Config) DBName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/playground?parseTime=true", c.DBUser, c.DBPassword, c.DBHost, c.DBPort)
}

func bindEnv(iface interface{}, ps ...string) {
	iv := reflect.ValueOf(iface)
	it := reflect.TypeOf(iface)
	for i := 0; i < it.NumField(); i++ {
		v := iv.Field(i)
		t := it.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnv(v.Interface(), append(ps, tv)...)
		default:
			_ = viper.BindEnv(strings.Join(append(ps, tv), "."))
		}
	}
}
