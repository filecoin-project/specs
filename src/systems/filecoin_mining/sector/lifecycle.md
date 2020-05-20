In Filecoin, miners contribute storage capacity to the network in units of Sectors. These sectors work similar to real-life shipping containers; they are used provide a unique ID for storage / retrieval processes as well as ensuring the datas 'dimensions' conform with all other Sectors in the network.

#### Sector creation

At creation, a sectors space (raw-byte power or sector size) and time (duration) are defined. The defined space will determine the total sum of the sectors pledge collateral in FIL. The new sector may contain deals (either Deals or VerifiedDeals), Committed Capacity or a mixture of both. 

The sector is then assigned a `SectorQuality`. which determines its Quality-Adjusted Power in the network, or consensus power (weight).

`SectorQuality` is determined through a weighted average of multipliers, based on their contents:

* Sectors full of `VerifiedDeals` will have a `SectorQuality` of `VerifiedDealWeightMultiplier/BaseMultiplier`.
* Sectors full of regular `Deals` will have a `SectorQuality` of `DealWeightMultiplier/BaseMultiplier`.
* Sectors with neither will have a `SectorQuality` of `BaseMultiplier/BaseMultiplier`.

Once `SectorQuality` has been assigned, the Sector is now ready for use. 

#### Committed Capacity (CC) upgrades

A sector entirely composed of Committed Capacity can later be upgraded to a Deals sector. This is currently done by resealing, though there are plans to make CC upgrades more efficient and cost-effective after the launch of mainnet.

#### Sector termination 

All sectors are expected to remain live until the end of their defined timespan and early dropping of sectors will result in slashing. This is done to provide clients a certain level  of guarantee on the reliability of their hosted data.

Miners may also choose to terminate a sector voluntarily and accept this slashing penalty.

#### Sector extensions

Miners can extend the lifetime of a sector at any time, though the Sector will be expected to remain live until its new end date. This can be done by submitting a `ExtendedSectorExpiration` message to the chain.


