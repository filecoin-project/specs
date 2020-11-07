---
title: Block Reward Minting
bookCollapseSection: true
weight: 2
dashboardWeight: 1
dashboardState: reliable
dashboardAudit: n/a
dashboardTests: 0
math-mode: true
---

# Block Reward Minting

In this section, we provide the mathematical specification for Simple Minting, Baseline Minting and Block Reward Issuance. We will provide the details and key mathematical properties for the above concepts.

## Economic parameters



- `$M_\infty$` is the total asymptotic number of tokens to be emitted as storage-mining block rewards. Per the [Token Allocation spec]{https://spec.filecoin.io/#systems__filecoin_token__token_allocation}, `$M_\infty := 55\% \cdot \texttt{FIL\_BASE} = 0.55 \cdot 2\times 10^9 FIL = 1.1 \times 10^9 FIL$`. The dimension of the `$M_\infty$` quantity is tokens.

- `$\lambda$` is the "simple exponential decay" minting rate corresponding to a 6-year half-life. The meaning of "simple exponential decay" is that the total minted supply at time `$t$` is `$M_\infty \cdot (1 - e^{-\lambda t})$`, so the specification of `$\lambda$` in symbols becomes the equation `$1 - e^{-\lambda \cdot 6yr} = \frac{1}{2}$`. Note that a "year" is poorly defined. The simplified definition of `$1yr := 365d$` was agreed upon for Filecoin. Of course, `$1d = 86400s$`, so `$1yr = 31536000s$`. We can solve this equation as

{{<katex>}}
\[\lambda = \frac{\ln 2}{6yr} = \frac{\ln 2}{189216000s} \approx 3.663258818 \times 10^{-9} Hz\] 
{{</katex>}}

The dimension of the `$\lambda$` quantity is `time$^{-1}$`.

- `$\gamma$` is the mixture between baseline and simple minting. A `$\gamma$` value of 1.0 corresponds to pure baseline minting, while a `$\gamma$` value of 0.0 corresponds to pure simple minting. Per the [Parameter Recommendation spec]{https://spec.filecoin.io/#algorithms__cryptoecon__initial-parameter-recommendation}, we currently use `$\gamma := 0.7$`. The `$\gamma$` quantity is dimensionless.

- `$b(t)$` is the baseline function, which was designed as an exponential 

{{<katex>}}
$$b(t) = b_0 \cdot e^{g t}$$
{{</katex>}}

where 


  - `$b_0$` is the "initial baseline". The dimension of the `$b_0$` quantity is information.
  - `$g$` is related to the baseline's "annual growth rate" (`$g_a$`) by the equation `$\exp(g \cdot 1yr) = 1 + g_a$`, which has the solution

{{<katex>}}
  $$g = \frac{\ln\left(1 + g_a\right)}{31536000s}.$$
{{</katex>}}

While `$g_a$` is dimensionless, the dimension of the `$g$` quantity is `time$^{-1}$`.

The dimension of the `$b(t)$` quantity is information.

## Simple Minting

- `$M_{\infty B}$` is the total number of tokens to be emitted via baseline minting: `$M_{\infty B} = M_\infty \cdot \gamma$`. Correspondingly, `$M_{\infty S}$` is the total asymptotic number of tokens to be emitted via simple minting: `$M_{\infty S} = M_\infty \cdot (1 - \gamma)$`. Of course, `$M_{\infty B} + M_{\infty S} = M_\infty$`.

- `$M_S(t)$` is the total number of tokens that should ideally have been emitted by simple minting up until time `$t$`. It is defined as `$M_S(t) = M_{\infty S} \cdot (1 - e^{-\lambda t})$`. It is easy to verify that `$\lim_{t\rightarrow\infty} M_S(t) = M_{\infty S}$`.


Note that `$M_S(t)$` is easy to calculate, and can be determined quite independently of the network's state. (This justifies the name "simple minting".)


## Baseline Minting

To define `$M_B(t)$` (which is the number of tokens that should be emitted up until time `$t$` by baseline minting), we must introduce a number of auxiliary variables, some of which depend on network state.


- `$R(t)$` is the instantaneous network raw-byte power (the total amount of bytes among all active sectors) at time `$t$`. This quantity is state-dependent---it depends on the activities of miners on the network (specifically: commitment, expiration, faulting, and termination of sectors). The dimension of the `$R(t)$` quantity is information. 

- `$\overline{R}(t)$` is the capped network raw-byte power, defined as `$\overline{R}(t):= \min\{b(t), R(t)\}$`. Its dimension is also information. 

- `$\overline{R}_\Sigma(t)$` is the cumulative capped raw-byte power, defined as `$\overline{R}_\Sigma(t) := \int_0^t \overline{R}(x)\, \mathrm{d}x$`. The dimension of `$\overline{R_\Sigma}(t)$` is `information$\cdot$time` (a dimension often referred to as "spacetime").

- `$\theta(t)$` is the "effective network time", and is defined as the solution to the equation

{{<katex>}}
    $$\int_0^{\theta(t)} b(x)\, \mathrm{d}x = \int_0^t \overline{R}(x)\, \mathrm{d}x = \overline{R}_\Sigma(t)$$
{{</katex>}}

By plugging in the definition of `$b(x)$` and evaluating the integral, we can solve for a closed form of `$\theta(t)$` as follows:

{{<katex>}}
    $$\int_0^{\theta(t)} b(x)\, \mathrm{d}x = \frac{b_0}{g} \left( e^{g\theta(t)} - 1 \right) = \overline{R}_\Sigma(t)$$
{{</katex>}}

{{<katex>}}
    $$\theta(t) = \frac{1}{g} \ln \left(\frac{g \overline{R}_\Sigma(t)}{b_0}+1\right)$$
{{</katex>}}

- `$M_B(t)$` is defined similarly to `$M_S(t)$`, just with `$\theta(t)$` in place of `$t$` and `$M_{\infty B}$` in place of `$M_{\infty S}$`: 

{{<katex>}}
$$M_B(t) = M_{\infty B} \cdot \left(1 - e^{-\lambda \theta(t)}\right)$$
{{</katex>}}

## Block Reward Issuance


- `$M(t)$`, the total number of tokens to be emitted as expected block rewards up until time `$t$`, is defined as the sum of simple and baseline minting:

{{<katex>}}
$$M(t) = M_S(t) + M_B(t)$$
{{</katex>}}


Now we have defined a continuous target trajectory for _cumulative_ minting. But minting actually occurs _incrementally_, and also in _discrete_ increments. Periodically, a "tipset" is formed consisting of multiple winners, each of which receives an equal, finite amount of reward. A single miner may win multiple times, but may only submit one block and may still receive rewards _as if_ they submitted multiple winning blocks. The mechanism by which multiple wins are rewarded is multiplication by a variable called `WinCount`, so we refer to the finite quantity minted and awarded for each win as "reward per `WinCount`" or "per win reward".

- `$\tau$` is the duration of an "epoch" or "round" (these are synonymous). Per the [spec](https://spec.filecoin.io/#glossary__epoch), `$\tau = 30s$`. The dimension of `$\tau$` is time.
- `$E$` is a parameter which determines the expected number of wins per round. While `$E$` could be considered dimensionless, it useful to give it a dimension of "wins". In Filecoin, the value of `$E$` is 5.
- `$W(n)$` is the total number of wins by all miners in the tipset during round `$n$`. This also has dimension "wins". For each `$n$`, `$W(n)$` is a random variable with the independent identical distribution `$\mathrm{Poisson}(E)$`.
- `$w(n)$` is the "reward per `WinCount`" or "per win reward" for round `$n$`. It is defined by:

{{<katex>}}
$$w(n) = \frac{\max\{M(n\tau+\tau) - M(n\tau),0\}}{E}$$
{{</katex>}}

The dimension of $W(n)$ is `tokens$\cdot$wins$^{-1}$`.

- While `$M(t)$` is a continuous target for minted supply, the discrete and random amount of tokens which have been minted as of time `$t$` is

{{<katex>}}
$$m(t) = \sum_{k=0}^{\left\lfloor t/\tau\right\rfloor-1} w(k) W(k)$$ 
{{</katex>}}

`$m(t)$` depends on past values of both `$W(n)$` and `$R(n\tau)$`.
