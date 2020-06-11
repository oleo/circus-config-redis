package messagebus
import(
	"log"
	"context"
	//"errors"
	//"regexp"
	"strconv"
	"time"
	"encoding/json"
	"github.com/go-redis/redis"
  circusNode "github.com/oleo/circus-redis/node"


)
type Init struct {
	Host string
	Port int
	Topic string
}
type status struct {
	Enrolled bool
	Enrolling bool
}

type ID struct {
  Uuid string `json:"uuid"`
	Name string `json:"name"`
	Timestamp int64 `json:"timestamp"`
}
type Heartbeats struct {
	Timestamp int64 `json:"timestamp"`
	Beats []ID
}

type conKey string
var main_client  *redis.Client

var State = status{}

func init() {
	State.Enrolled = false
	State.Enrolling = false
}
func (e *Init) GetState() *status {
	return &State
}
func (e *Init) GetIdStruct() ID {
	_id := ID{}
	return _id
}
func (e *Init) Ctx() context.Context {
	k := conKey("jalla")
	ctx := context.WithValue(context.Background(),k, "Goredisssss")
	return ctx
}
func rClient(host string) *redis.Client {
	k := conKey("jalla")
	ctx := context.WithValue(context.Background(),k, "Goredisssss")

	client := redis.NewClient(&redis.Options{
		Addr: host,
		Password: "",
		DB: 0,
	})
	// Test connection
	err := ping(ctx,client)
	if err != nil {
		log.Println(err)
	}
	return client
}

func ping(ctx context.Context, client *redis.Client) error {
	_ , err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}
func (e *Init) GetClient() *redis.Client {

	k := conKey("jalla")
	ctx := context.WithValue(context.Background(),k, "Goredisssss")

	if(main_client == nil) {
		main_client = rClient(e.Host+":"+strconv.Itoa(e.Port))
	}
	err := ping(ctx,main_client)
	if err != nil {
		// Not able to use main_client. close and open again
		client := rClient(e.Host+":"+strconv.Itoa(e.Port))
		log.Println(err)
		return client
	}

	// All ok
	return main_client

}
func (e *Init) Getstr(key string) string {
	out:=""
	client := e.GetClient()

	Val, err := client.Get(e.Ctx(),key).Result()
	if err == redis.Nil {
		log.Println("no value found")
	} else if err != nil {
		panic(err)
	} else {
		out=Val
	}
	return out

}
func (e *Init) Getbool(key string) (bool,error) {
	out:=false

	Val := e.Getstr(key)
	log.Printf("Found %s for %s\n",Val,key)
	b, err := strconv.ParseBool(Val)
		if err != nil {
			return out,err
		}
		return b,nil

}

func (e *Init) SetTopic(key string,value string,ttl int) {
	client := e.GetClient()

	err := client.Set(e.Ctx(),key,value,time.Second * time.Duration(ttl)).Err()
	if err != nil {
		panic(err)
	}
	return
}
func (e *Init) PublishTopic(topic string,msg string) {
	client := e.GetClient()

	err := client.Publish(e.Ctx(),topic,msg).Err()
	if err != nil {
		panic(err)
	}
	return
}
func (e *Init) Set(key string,value string,ttl int) {
	e.SetTopic(e.Topic+key,value,ttl)
}
func (e *Init) Publish(msg string) {
	e.PublishTopic(e.Topic,msg)
}
func (e *Init) PublishService(service string,msg string) {
	e.PublishTopic(e.Topic+service,msg)
}

func (e *Init) Action(action string,msg string) {
	e.PublishService("action/"+action,msg)
}
// Actions
func (e *Init) Enroll(node circusNode.CircusNode) {
	node.State.Enrolling=true
	id_ := ID{	Uuid: node.Uuid,	Name: node.Name, Timestamp: node.Timestamp	}
	jsondata_id, _ := json.Marshal(id_)
	e.Action("enroll",string(jsondata_id))
}

// Sets
func (e *Init) SendHeartbeat(node circusNode.CircusNode,secs int) {
	node.State.Enrolling=true
	id_ := ID{	Uuid: node.Uuid,	Name: node.Name, Timestamp: time.Now().Unix()	}
	jsondata_id, _ := json.Marshal(id_)
	e.Set("heartbeat/"+node.Uuid,string(jsondata_id),secs)
}

// Gets
func (e *Init) Enrolled(id string) bool {
	valuestr,err := e.Getbool(e.Topic+"enrolled/"+id)
	if err != nil {
		log.Printf("Error retrieving boolean value from redis.")
		return false
	}
	return valuestr
}

func (e *Init) GetHeartbeats() Heartbeats {
	var hb_struct Heartbeats
	hb_struct.Timestamp = time.Now().Unix()
	
	client := e.GetClient()

	// get list of keys
	keys,err := client.Keys(e.Ctx(),e.Topic+"heartbeat/*").Result()
	if err != nil {
		panic(err)
	}
	//re := regexp.MustCompile(`heartbeat/(.+)$`)
	for _, key := range keys {
		// for each key:
		// retrieve value (should be json) 
		value,err := client.Get(e.Ctx(),key).Result()
		if err != nil {
			panic(err)
		}
		// unmarshall to ID-struct
		var tmp = e.GetIdStruct()
		err = json.Unmarshal([]byte(value), &tmp)

		// append to hb_struct
		hb_struct.Beats = append(hb_struct.Beats, tmp)
		//log.Printf("Got %s in %s\n",re.FindStringSubmatch(key)[1],value)

	}
	// Return hb_struct
	return hb_struct
}

