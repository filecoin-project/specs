# Pledge Collateral

Filecoin includes a concept of "Pledge Collateral", which is FIL collateral that storage miners must lock up when participating as miners.

Pledge collateral serves several functions in Filecoin. It:
- makes it possible to slash misbehaving or slow miners
- ensures that miners have skin in the game (for the Filecoin network as a whole)
- increases the cost of launching a 51% attack


## Computing Pledge Collateral

The total pledge collateral across all miners is a fixed proportion of available FIL.
Available FIL is computed as the total amount of FIL that has been mined, plus the total amount of FIL that's been vested, minued the amount of FIL which has been burned.

```go
availableFil := minedFil + vestedFil - burnedFil
```

Pledge collateral is subdivided into two kinds: power collateral and per-capita collateral.
Power collateral is split across miners according to their share of the total network power, and per-capita collateral is split across miners evenly.
Two parameters, `POWER_COLLATERAL_PROPORTION` and `PER_CAPITA_COLLATERAL_PROPORTION`, relate the total amount of collateral to the `availableFil`.

```go
totalPowerCollateral := availableFil * POWER_COLLATERAL_PROPORTION
totalPerCapitaCollateral := availableFil * PER_CAPITA_COLLATERAL_PROPORTION
totalPledgeCollateral := totalPowerCollateral + totalPerCapitaCollateral
```

Power-based collateral ensures that miners' collateral is proportional to their economic size and to their expected rewards.
The presence of per-capital collateral acts as a deterrent against Sibyl attacks.
We intend for the `POWER_COLLATERAL_PROPORTION` to be several times larger than the `PER_CAPITA_COLLATERAL_PROPORTION`.

To calculate any particular miner's collateral requirements, we need to know the miner's power, the total network power, and the total number of miners in the network.

```go
minerPowerCollateral := totalPowerCollateral * minerPower / totalNetworkPower
minerPerCapitaCollateral := totalPerCapitaCollateral / numMiners
```

Putting all these variables together, we have each miner's individual collateral requirement:
```go
minerPledgeCollateral := availableFil * ( POWER_COLLATERAL_PROPORTION * minerPower / totalNetworkPower PER_CAPITA_COLLATERAL_PROPORTION / numMiners)
```

## Dealing with Undercollateralization

In the course of normal events, miners may become undercollateralized.

They cannot directly undercollateralized themselves by adding more power, as commitSector will fail if they do not have sufficient collateral to cover their power requirements.
However, their collateral requirement could increase due to growth in availableFil, a reduction in the total network power, or a reduction in the total number of miners.
In such cases, the miner may continue to submit PoSts and mine blocks. When they win blocks, their block rewards will be garnished while they remain undercollateralized.

## Parameter Choices

We provisionally propose the following two parameters choices:

```go
POWER_COLLATERAL_PROPORTION := 0.2
PER_CAPITA_COLLATERAL_PROPORTION := 0.05
```

These are subject to change before launch.
