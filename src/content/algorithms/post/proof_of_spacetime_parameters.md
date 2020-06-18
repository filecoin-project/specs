+++
title = "PoSt Parameters"
author = ["Chhi'med Kunzang"]
draft = false
+++

DEPRECATED: Needs to be updated with WinningPoSt and with WindowedPoSt.

This section describes parameters for Rational-PoSt, the Proof-of-Spacetime used in Filecoin.

| Parameter             | Type   | Value | Description                                                                                                                                                                                    |
|-----------------------|--------|-------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| POST-CHALLENGE-BLOCKS | BLOCKS | 480   | The time offset before which the actual work of generating the PoSt cannot be started. This is some delta before the end of the Proving Period, and as such less than a single Proving Period. |
| POST-CHALLENGE-HOURS  | HOURS  | 2     | PoSt challenge time (see POST\_CHALLENGE\_BLOCKS).                                                                                                                                             |
| POST-PROVING-PERIOD   | BLOCKS | 5760  | The time interval in which a PoSt has to be submitted                                                                                                                                          |

{{% notice todo %}}
**TODO**: The above values are tentative and need both backing from research as well as detailed reasoning why we picked them.
{{% /notice %}}
