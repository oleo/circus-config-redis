package circusNode
import(
	"time"
	guuid "github.com/google/uuid"
)
type circusState struct {
		Enrolled bool
		Enrolling bool
}
type CircusNode struct {
  Uuid string `json:"uuid"`
	Name string `json:"name"`
	Timestamp int64 `json:"timestamp"`
	State circusState
}

func New(name string) CircusNode {
	initial_state := circusState{false,false}
	uuid := guuid.New()
	e := CircusNode{	uuid.String(),	name,time.Now().Unix(),initial_state}
	return e
}

func (e *CircusNode) isEnrolled() bool {
	return e.State.Enrolled
}
