package agents

type AgentMetaList struct {
	Agents []Agent
}

func NewAgentMetaList() AgentMetaList {
	return AgentMetaList{
		Agents: make([]Agent, 128),
	}
}
