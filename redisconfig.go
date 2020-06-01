package redisconfig
import(
	"encoding/json"
	"log"
	"context"
	//"fmt"
	"github.com/go-redis/redis"
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
}
type Init struct {
	Host string
}
type conKey string

func (e *Init) Get(key string) StringConfig {
	k := conKey("jalla")
	ctx := context.WithValue(context.Background(),k, "Goredisssss")
	// get connection to redis and get json-config
	cli := rClient(e.Host)
	err := ping(ctx,cli)
	if err != nil {
		log.Println(err)
	}
	raw_config := getstr(ctx,cli,key)
//	"circus:test-string-1:config")
	//fmt.Printf("Got config:\n%s\n-\n",raw_config)
	Config := StringConfig{}
	err = json.Unmarshal([]byte(raw_config), &Config)
	if err != nil {
		log.Fatal("Can't deode config JSON: ",err)
	}
	//fmt.Println("Will be using:")
	//fmt.Printf("Host:   %s:%d\n",Config.MessageBus.Host,Config.MessageBus.Port)
	//fmt.Print("Topic:   ")
	//fmt.Println(Config.MessageBus.Topic)

	return Config
}

func rClient(host string) *redis.Client {
	//fmt.Printf("Will conntact redis at %s\n",host)
	client := redis.NewClient(&redis.Options{
		Addr: host,
		Password: "",
		DB: 0,
	})

	return client
}

func ping(ctx context.Context, client *redis.Client) error {
	_ , err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func getstr(ctx context.Context, client *redis.Client,key string) string {
	out:=""
//	fmt.Printf("Will try to get %s\n",key)
	Val, err := client.Get(ctx,key).Result()
	if err == redis.Nil {
		log.Println("no value found")
	} else if err != nil {
		panic(err)
	} else {
		out=Val
		return out
	}
	return out

}
