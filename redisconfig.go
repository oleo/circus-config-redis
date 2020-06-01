package redisconfig
import(
	"encoding/json"
	"flag"
	"log"
//	"strings"
	"context"
	"fmt"
	"os"
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
func get(key string) {
	k := conKey("jalla")
	ctx := context.WithValue(context.Background(),k, "Goredisssss")
	// get connection to redis and get json-config
	cli := rClient(Host)
	err = ping(ctx,cli)
	if err != nil {
		fmt.Println(err)
	}
	raw_config := getstr(ctx,cli,key)
//	"circus:test-string-1:config")
	fmt.Printf("Got config:\n%s\n-\n",raw_config)
	Config := StringConfig{}
	err = json.Unmarshal([]byte(raw_config), &Config)
	if err != nil {
		log.Fatal("Can't deode config JSON: ",err)
	}
	fmt.Println("Will be using:")
	fmt.Printf("Host:   %s:%d\n",Config.MessageBus.Host,Config.MessageBus.Port)
	fmt.Print("Topic:   ")
	fmt.Println(Config.MessageBus.Topic)


}

func rClient(host string) *redis.Client {
	fmt.Printf("Will conntact redis at %s\n",host)
	client := redis.NewClient(&redis.Options{
		Addr: host,
		Password: "",
		DB: 0,
	})

	return client
}

func ping(ctx context.Context, client *redis.Client) error {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println(pong,err)

	return nil
}

func getstr(ctx context.Context, client *redis.Client,key string) string {
	out:=""
	fmt.Printf("Will try to get %s\n",key)
	Val, err := client.Get(ctx,key).Result()
	if err == redis.Nil {
		fmt.Println("no value found")
	} else if err != nil {
		panic(err)
	} else {
		out=Val
		return out
	}
	return out

}
