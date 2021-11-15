# Game Design Document
# _CX Game_

> **Sci-fi sandbox game with survival elements**, inspired by games like Terraria,
> Starbound, Oxygen Not Included, Factorio, and Dyson Sphere Program.

## 🌎 **2D World Map (side view)** 🌎

### _Environment Layer_
It's placed **behind the background layer**. Its elements are just to define the environment of the planet/asteroid, so they are **indestructible**.

This layer has a small **parallax** effect on the **horizontal axis**, so its elements move less than the foreground layer.

### _Background Layer_
At first, this layer will only contain walls, that can be built by the player like any other object but **removed by a different tool** (On Terraria is the hammer).
- **Common walls:** These are the ones used to build the base and fill caves.
- **Mining areas:** These walls have lots of ores within. The player can put machines on the middle layer to mine these ores. These walls are probably going to use an auto tiling system to reduce their size as the resources are going to an end.

### _Middle Layer_
This layer is mainly composed of elements fixed on the walls, **lights** and **Pipesim** ([Wires](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#wire-circuits), [Liquid Pipes](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#liquid-pipes), [Gas Pipes](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#gas-pipes) and [Automation](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#-automation-system)).

The **windows** are going to be put on this layer and the walls behind them are cropped by a cutout shader. A window can only be placeable if there are walls behind each of its tiles.

### _Foreground Layer_
Most of the game is here. This layer is going to have machines, plants, enemies, spaceships, furniture, etc. All of these elements are still behind the player.

### _Layer over the player_
This layer contains the [liquids](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#liquids) (water, lava, biofluid, etc), [gases](https://github.com/skycoin/cx-game/blob/main/docs/game-design.md#gases), and some materials that go over the player, like the top part of a glass transportation tube.

### _Pipesim Overlay_
While editing a Pipesim element, the respective element (e.g. wires) is going to be highlighted. Also, some icons will appear on the objects to indicate where the player should connect the Pipesim.

![Overlay example](https://preview.redd.it/xhoa1gvod3v11.jpg?width=1013&format=pjpg&auto=webp&s=56e891501f255c8ae0cded9eed65c04a7dc24ff3)
> _Example of overlay on Oxygen Not Included_


## 🧱 **Base elements** 🧱

### ⚡ _Power System_ 

##### Wire circuits

##### Power Generators

##### Batteries

---
### 💨 _Gases System_

##### Gases

##### Gas pipes

##### Gas Pump

##### Gas Vent

---
### 🚿 _Plumbing System_

##### Liquids

##### Liquid pipes

##### Liquid Pump

##### Liquid Vent

---
### 🍖 _Food System_

##### Food

##### Cultivating plants

##### Creatures husbandry

---
### 👨‍💻 _Automation System_

##### Automation wires

##### Logic gates

##### Conveyor system

##### Droids

## 💊 **Survival elements** 💊
- **HP:** 
- **OXYGEN:** 
- **TEMPERATURE:** 
- **FOOD:** 
- **DISEASES:** 

## 🚀 **Space exploration** 🚀
### _Spaceships_ 

##### Building ships

### _Overworld_ 

##### Controls

##### Asteroids

##### Enemies

### _Planets_ 

##### Moon

##### Volcanic moon

##### Ocean Planet

##### Frozen Planet

##### Circuit Planet
