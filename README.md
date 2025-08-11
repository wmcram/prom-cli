# A CLI Helper for Prometheus endpoints

### Usage

`promcli` currently supports the following subcommands:
- `get`: for pretty-printing a metrics endpoint with filtering
- `watch`: for seeing a live view of selected metrics
- `graph`: like above, but visualized
- `mock`: for serving a mock metrics endpoint from a file

### Installing
```bash
git clone git@github.com:wmcram/prom-cli.git
cd prom-cli
make install
cd ..
rm -rf prom-cli
```