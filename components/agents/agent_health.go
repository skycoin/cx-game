package agents

type HealthComponent struct {
	Current int
	Max     int
	Died    bool
}

func NewHealthComponent(max int) HealthComponent {
	return HealthComponent{Current: max, Max: max, Died: false}
}
