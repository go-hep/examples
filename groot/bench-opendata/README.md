# bench-opendata

`bench-opendata` is a set of naive analyses intended to reproduce the ones at [IRIS-HEP](https://github.com/iris-hep/adl-benchmarks-index/) and show how one would use [groot](https://go-hep.org/x/hep/groot).

Currently, all 8 analyses have been implemented using the "basic" `groot` interface.

More analyses using different styles (`rsql`, `basic+struct` or `rarrow+dframe`) might appear if time permits (PR accepted!)

## Example

```
$> go get github.com/go-hep/examples/groot/bench-opendata
$> bench-opendata -help
Usage of ./bench-opendata:
  -bench string
    	comma-separated list of opendata benchmark examples to run (01-basic,01-rsql,02-basic,03-basic,04-basic,05-basic,06-basic,07-basic,08-basic)
  -f string
    	input file to analyze (default "root://eospublic.cern.ch//eos/root-eos/benchmark/Run2012B_SingleMu.root")
  -list
    	list all available benchmarks and exits
  -profile
    	enable/disable CPU profiling

$> bench-opendata -list
bench-opendata: available OpenData benchmark examples: ["01-basic" "01-rsql"
 "02-basic" "03-basic" "04-basic" "05-basic" "06-basic" "07-basic" "08-basic"]

$> bench-opendata  -f ./testdata/Run2012B_SingleMu.root
bench-opendata: running benchs: ["01-basic" "01-rsql" "02-basic" "03-basic" "04-basic" "05-basic" "06-basic" "07-basic" "08-basic"]
bench-opendata: running "01-basic"...
tree: 53446198 entries
hmet: 5.3446198e+07
bench-opendata: running "01-basic"... [err=<nil>] delta=12.551175919s
bench-opendata: running "01-rsql"...
tree: 53446198 entries
hmet: 5.3446198e+07
bench-opendata: running "01-rsql"... [err=<nil>] delta=1m9.692351135s
bench-opendata: running "02-basic"...
tree: 53446198 entries
hJetPt: 1.70952895e+08
bench-opendata: running "02-basic"... [err=<nil>] delta=34.778735989s
bench-opendata: running "03-basic"...
tree: 53446198 entries
hJetPt: 6.3845511e+07
bench-opendata: running "03-basic"... [err=<nil>] delta=47.967790653s
bench-opendata: running "04-basic"...
tree: 53446198 entries
hmet: 2.349018e+06
bench-opendata: running "04-basic"... [err=<nil>] delta=50.719104916s
bench-opendata: running "05-basic"...
tree: 53446198 entries
hmet: 2.823648e+06
bench-opendata: running "05-basic"... [err=<nil>] delta=1m19.357797894s
bench-opendata: running "06-basic"...
tree: 53446198 entries
h1: 109698
h2: 36566
bench-opendata: running "06-basic"... [err=<nil>] delta=1.922018984s
bench-opendata: running "07-basic"...
tree: 53446198 entries
h1: 37446
bench-opendata: running "07-basic"... [err=<nil>] delta=515.737113ms
bench-opendata: running "08-basic"...
tree: 53446198 entries
hmet: 8644
hlep: 8644
bench-opendata: running "08-basic"... [err=<nil>] delta=307.818562ms
```

![basic-08](https://github.com/go-hep/examples/raw/master/groot/bench-opendata/imgs/08-basic.png)
