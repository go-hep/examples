# groot

This repository holds various examples related to `groot`, the pure-Go package that deals with ROOT aspects.

## ex-tree-00

`ex-tree-00` creates a ROOT file with a simple flat TTree `tree`:

```
$> go run ./ex-tree-00.go
groot: -- created tree "tree":
groot: branch[0]: name="evt_i", title="evt_i/L"
groot: branch[1]: name="evt_a_e", title="evt_a_e/D"
groot: branch[2]: name="evt_a_t", title="evt_a_t/D"
groot: branch[3]: name="evt_b_e", title="evt_b_e/D"
groot: branch[4]: name="evt_b_t", title="evt_b_t/D"
groot: processing event 0...
groot: evt.i=          0
groot: evt.a.e=   -1.234
groot: evt.a.t=    0.665
groot: evt.b.e=   -0.126
groot: evt.b.t=    1.519
[...]
groot: processing event 9000...
groot: evt.i=       9000
groot: evt.a.e=    0.810
groot: evt.a.t=    0.304
groot: evt.b.e=    1.310
groot: evt.b.t=    0.094
```

## ex-tree-01

`ex-tree-01` reads the ROOT file created with `ex-tree-00` and displays the tree's content:

```
$> go run ./ex-tree-01.go
groot: processing event 0...
groot: ievt: 0
groot: evt.a.e=   -1.234
groot: evt.a.t=    0.665
groot: evt.b.e=   -0.126
groot: evt.b.t=    1.519
[...]
groot: processing event 9000...
groot: ievt: 9000
groot: evt.a.e=    0.810
groot: evt.a.t=    0.304
groot: evt.b.e=    1.310
groot: evt.b.t=    0.094
```
