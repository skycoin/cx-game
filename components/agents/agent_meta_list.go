package agents

// not used
type AgentMetaList struct {
	AgentMetas []*AgentMeta
}

func NewAgentMetaList() AgentMetaList {
	return AgentMetaList{
		AgentMetas: make([]*AgentMeta, 128),
	}
}
