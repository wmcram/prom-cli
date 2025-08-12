# A CLI Helper for Prometheus endpoints

`promcli` is an easy-to-use command-line tool for interacting with 
Prometheus endpoints. Its output is explicitly designed to be human-readable
and easier to filter than `grep`ing and `awk`ing your way through `curl` output.

### Installing

The following will install `promcli` to your `GOBIN`:

```bash
git clone git@github.com:wmcram/prom-cli.git
cd prom-cli
make install
cd ..
rm -rf prom-cli
```

As long as `GOBIN` is in your `PATH`, you should be able to just type
`promcli` to use it.

### Usage

`promcli` currently supports the following subcommands:
- `get`: for pretty-printing a metrics endpoint with filtering
- `watch`: for seeing a live view of selected metrics as text
- `graph`: like above, but visualized as a graph
- `mock`: for serving a mock metrics endpoint from a file

