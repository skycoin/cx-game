package agents

import (
	"fmt"
)

type AgentList struct {
	Agents []Agent
}

func (al *AgentList) Draw() {

}

func (al *AgentList) Tick() {
	for a := range al.Agents {
		fmt.Println(a)
	}
}
