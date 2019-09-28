+++
title = "Proof of Spacetime Parameters"
author = ["Chhi'med Kunzang"]
draft = false
+++

This section describes parameters for Rational-PoSt, the Proof-of-Spacetime used in Filecoin.

| Parameter             | Type   | Value | Description                                                                                                                                                                                              |
|-----------------------|--------|-------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| POST-CHALLENGE-BLOCKS | BLOCKS | 480   | The time offset at which the actual work of generating the PoSt can not be started earlier than. This is some delta before the end of the Proving Period, and as such less then a single Proving Period. |
| POST-CHALLENGE-HOURS  |        | 2     |                                                                                                                                                                                                          |
| POST-PROVING-PERIOD   | BLOCKS | 5760  | The time interval in which a PoSt has to be submitted                                                                                                                                                    |

{{% notice todo %}}
**TODO**: The above values are tentative and need both backing from research as well as detailed reasoning why we picked them.
{{% /notice %}}
