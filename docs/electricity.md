# Electricity

The game will contain machines that both product and consume power.
These machines will be connected by wires.

## Requirements

Toggling a switch's power must be a fast operation.
Similarly, placing or destroying wire tiles must not cause lag.
Since there will usually be more switch-flipping than wire-destroying,
the switches power is the critical use case.

## Structure

A **circuit** is a continuously connected segment of wire
and the directly connected machines.

A **circuit group** is a set of circuits which are connected by a 
diode or switch.

When two circuits are connected by a closed switch, 
power transmission is bidirectional.
In contrast, when two circuits are connected by diode,
power transmission is unidirectional.
When two circuits are connected by an open switch,
they are NOT considered part of the same circuit group 
and therefore no power transmission is allowed.

**Batteries** are a special case, 
as they draw/store a variable quantity of power. 
Battery power is calculated after all other power calculations. 
If a circuit consumes more power than it draws, 
then the extra power is drawn from the battery. 
Similarly, if a circuit produces more power than it consumes, 
then the extra power is stored in the battery.

### Pseudocode

```golang
type CircuitID int
func (id CircuitID) Get() *Circuit // ...

type MachineID int
func (id MachineID) Get() *Machine // ...

type Machine struct {
	Power int
	IsBattery bool
}

type Circuit struct {
	MachineIDs []MachineID
	BasePower  int
}

func (c *Circuit) CalculateBasePower() {
	basePower := 0
	for _,id := range c.MachineIDs {
		machine := id.Get()
		if machine.IsStaticPower {
			basePower += machine.Power
		}
	}
	c.BasePower = basePower
}

func (c *Circuit) Batteries() []*Machine {
	batteries := []*Machine{}
	for _,id := range c.MachineIDs {
		machine := id.Get()
		if machine.IsBattery {
			batteries = append(batteries, machine)
		}
	}
	return batteries
}

var circuits = map[CircuitID]*Circuit{}

type CircuitGroup struct {
	CircuitIDs []CircuitID
}

type CircuitConnection struct {
	From,To CircuitID
	AllowForward, AllowBackward bool
	PowerDifference int
}

func (cc *CircuitConnection) CalculatePowerDifference() {
	cc.PowerDifference = cc.To.Get().BasePower - cc.From.Get().BasePower
}

func (cc *CircuitConnection) Allowed() bool {
	if cc.PowerDifference < 0 && cc.AllowForward { return true }
	if cc.PowerDifference > 0 && cc.AllowBackward { return true }
	return false
}

func TickCircuits() {
	for _,circuit := range circuits {
		circuit.CalculateBasePower()
	}
	for _,circuitGroup := range circuitGroups {
		for _,connection := circuitGroup.Connections {
			connection.CalculatePowerDifference()
			if connection.Allowed() {
				connection.TransmitPower()
			}
		}
	}
	for _,circuit := range circuits {
		for _,battery := circuit.Batteries() {
			
		}
	}
}

func (cg *CircuitGroup) Tick() {
	circuitBasePowers := make([]int, len(cg.CircuitIDs))
	for idx,id := range cg.CircuitIDs {
		circuitBasePowers[idx] = id.Get().BasePower()
	}
}
```
