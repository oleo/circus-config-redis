package redisconfig
import(
	"encoding/json"
	"log"
	"context"
	"github.com/go-redis/redis"
	rc "github.com/oleo/circus-redis/conn"
)
type StringConfig struct {
	Master struct {
		Host string
		Port int
	}
	MessageBus struct {
		Host string
		Port int
		Topic string
	}
	Version string
	Flavour string
	Members  []Member
}
type Member struct {
	Id string
	Type string
	Host string
}
type Init struct {
	Host string
}

func (e *Init) GetConfig(key string) StringConfig {
	rc.Init{e.Host}
	raw_config := rc.Getstr(key)
	Config := StringConfig{}
	err = json.Unmarshal([]byte(raw_config), &Config)
	if err != nil {
		log.Fatal("Can't deode config JSON: ",err)
	}

	return Config
}

