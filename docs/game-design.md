# Game Design Document
# _CX Game_

> **Sci-fi sandbox game with survival elements**, inspired by games like Terraria,
> Starbound, Oxygen Not Included, Factorio, and Dyson Sphere Program.

## 🌎 **2D World Map (side view)** 🌎

### 🌄 _Layers_

#### 1. Environment Layer
It's placed **behind the background layer**. Its elements are just to define the environment of the planet/asteroid, so they are **indestructible**.

This layer has a small **parallax** effect on the **horizontal axis**, so its elements move less than the foreground layer.

#### 2. Background Layer
At first, this layer will only contain walls, that can be built by the player like any other object but **removed by a different tool** (On Terraria is the hammer).
- **Common walls:** These are the ones used to build the base and fill caves.
- **Mining areas:** These walls have lots of ores within. The player can put machines on the middle layer to mine these ores. These walls are probably going to use an auto tiling system to reduce their size as the resources are going to an end.

#### 3. Middle Layer
This layer is mainly composed of elements fixed on the walls, **lights** and **Pipesim** ([Wires](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#wire-circuits), [Liquid Pipes](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#liquid-pipes), [Gas Pipes](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#gas-pipes) and [Automation](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#-automation-system)).

The **windows** are going to be put on this layer and the walls behind them are cropped by a cutout shader. A window can only be placeable if there are walls behind each of its tiles.

#### 4. Foreground Layer
Most of the game is here. This layer is going to have machines, plants, enemies, spaceships, furniture, etc. All of these elements are still behind the player.

#### 5. Layer over the player
This layer contains the [liquids](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#liquids) (water, lava, biofluid, etc), [gases](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#gases), and some materials that go over the player, like the top part of a glass transportation tube.

### 🔎 _Views_

#### Pipesim overlay
While editing a Pipesim element, the respective element (e.g. wires) is going to be highlighted. Also, some icons/effects will appear on the objects to indicate where the player should connect each Pipesim.

![Overlay example](https://preview.redd.it/xhoa1gvod3v11.jpg?width=1013&format=pjpg&auto=webp&s=56e891501f255c8ae0cded9eed65c04a7dc24ff3)
> _Example of overlay on Oxygen Not Included_


## 🧱 **Base elements** 🧱

### ⚡ _Power System_ 

#### Wire circuits
The wires autotiles like any other Pipesim element and each wire segment has a **circuit id**. If somehow the wires of a circuit get disconnected then it'll be split into two **circuit id**.

Each "circuit id" carries the information of:
- **Current state:** On/Off
- **Current power:** Current wattage into the circuit, based on the machines that are turned on.
- **Maximum power:** The amount of wattage that this circuit would carry if all the machines connected to it are turned on.

Each machine connected to a circuit should have a **circuit id**. So if we turned it on, the current power of the respective circuit is changed.

#### Power Generators

Machines that supply power to the circuit when turned on.

These machines usually require some kind of fuel (solar, gas, etc) to generate power and the amount of power generated is different on each kind of machine.

#### Power Consumers

Many machines of the game are going to consume energy to work. These consumption values will be subtracted from the amount of power available on the respective **circuit id**.

#### Batteries

The batteries can be both a power supplier or consumer, depending on the overall status of the **circuit id** it's connected to.

- If the **circuit is receiving more power than it needs**, the batteries will **consume power** to charge themselves.
- If the **circuit is demanding more power than the generators are providing**, the batteries will also **supply power** to the circuit and uncharge themselves.

#### Switches


---
### 💨 _Gases System_

#### Gases

Different gases will be spawned together with the world, among the caves. The gases have different densities and they slowly distribute based on this property. Less dense gases go up and more dense gases deposit on the bottom.

If the gases get in contact with a vacuum area, they'll be slowly drained to that area.

The player must be careful because if the gas gets drained on the surface of a planet with no atmosphere, the gas will continue to go up until it vanishes from the planet.

#### Geyser

Geysers are natural formations of the world, impossible to destroy. They are spawned over a black indestructible ore and releases gas from time to time in the world.

Each geyser has the following properties:
- **Gas released:** Which kind of gas does the geyser emit.
- **Emission temperature:** The temperature of the gases that came from that geyser.
- **Active time:** The amount of time it stays active releasing gases.
- **Dormant time:** The amount of time until it activates again.
- **Frequency:** The frequency it releases gases while active.
- **Quantity of gas released:** How much gas it emits in each eruption.

#### Gas pipes

#### Gas Pump

#### Gas Vent

---
### 🚿 _Plumbing System_

#### Liquids

#### Liquid pipes

#### Liquid Pump

#### Liquid Vent

---
### 🍖 _Food System_

#### Food

#### Cultivating plants

#### Creatures husbandry

---
### 👨‍💻 _Automation System_

#### Automation wires

#### Logic gates

#### Conveyor system

#### Droids

## 💊 **Survival elements** 💊
- **HP:** 
- **OXYGEN:** 
- **TEMPERATURE:** 
- **FOOD:** 
- **DISEASES:** 

## 🚀 **Space exploration** 🚀
### _Spaceships_ 

#### Building ships

### _Overworld_ 

#### Controls

#### Asteroids

#### Enemies

### _Planets_ 

#### Moon

#### Volcanic moon

#### Ocean Planet

#### Frozen Planet

#### Circuit Planet
