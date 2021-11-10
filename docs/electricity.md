# Electricity

The game will contain machines that both product and consume power.
These machines will be connected by wires for wired connections,
or laser emitter / storage crystal pairs for wireless connections.

## Requirements

Toggling a switch's power must be a fast operation.
Similarly, placing or destroying wire tiles must not cause lag.
Since there will usually be more switch-flipping than wire-destroying,
the switches power is the critical use case.

## Structure

A **circuit** is a continuously connected segment of wire
and the directly connected machines.

A **circuit group** is a set of circuits which are connected by a switch.

When two circuits are connected by a closed switch, 
they are considered part of the same circuit group,
and therefore power transmission is allowed.
When two circuits are connected by an open switch,
they are NOT considered part of the same circuit group 
and therefore power transmission is NOT allowed.

**Batteries** are a special case, 
as they draw/store a variable quantity of power. 
Battery power is calculated after all other power calculations. 
If a circuit consumes more power than it draws, 
then the extra power is drawn from the battery. 
Similarly, if a circuit produces more power than it consumes, 
then the extra power is stored in the battery.

A **laser emitter** is used to transmit power 
from one circuit to another wirelessely.
A **storage crystal** is used to recieve the wireless power transmision.
The storage crystal also has a battery charge internally.
The laser emitter periodically fires a power pulse beam, 
which delivers a packet of charge to the storage crystal.
If the recieving storage crystal is fully charged, 
then the laser emitter will not fire.

## Additional Considerations

Laser emitter should always display trace beam to paired storage crystal.
