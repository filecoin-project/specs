# Ubercalc visualizations

## Install

- Install [orient](https://github.com/filecoin-project/orient)
- Install [obs](https://github.com/nicola/obs)

## Usage

```
# Run these commands from the root folder of the spec repo on two different terminals
$ ORIENT_CACHE_DIR=/absolute/path/to/cache ./orient/bin/orient web --system=src/orient/rational-calc.orient --port=8000
$ obs serve 8081 src/orient/viz
```

