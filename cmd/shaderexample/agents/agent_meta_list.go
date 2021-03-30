package agents

type AgentMetaList struct {
	AgentMetas []AgentMeta
}

func NewAgentMetaList() AgentMetaList {
	return AgentMetaList{
		AgentMetas: make([]AgentMeta, 128),
	}
}
