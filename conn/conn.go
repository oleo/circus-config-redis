package conn
import(
	"log"
	"context"
	"github.com/go-redis/redis"
)
type Init struct {
	Host string
}
type conKey string


func rClient(host string) *redis.Client {
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

func (e *Init) Getstr(key string) string {
	out:=""
	k := conKey("jalla")
	ctx := context.WithValue(context.Background(),k, "Goredisssss")
	client := rClient(e.Host)
	err := ping(ctx,client)
	if err != nil {
		log.Println(err)
	}
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
