# Intel® Trust Authority TDX CLI
An easy-to-use command line interface for attesting TDX TEE(TD) with Intel Trust Authority. The CLI will be installed and run inside the TD.

## Installation
### Before you begin
Make sure you have build essential and dependencies installed.
#### Ubuntu
```sh
sudo apt install build-essential
sudo snap install go --classic
```
#### SLES
```sh
sudo zypper install git make go
```

Note: make sure go is 1.19 or newer. 
```sh
go version
```
### Get the code
Checkout the code
```sh
git clone https://github.com/intel/trustauthority-client -b gcp-tdx-preview
```

### Build
Compile Intel Trust Authority TDX CLI. This will generate `trustauthority-cli` binary in current directory:

```sh
cd trustauthority-client/tdx-cli/
make cli
```

## Usage

### To get list of all the available commands

```sh
./trustauthority-cli --help
```
More info about a specific command can be found using
```sh
./trustauthority-cli <command> --help
```

### To get Intel Trust Authority signed token

`token` command requires Intel Trust Authority configuration to be passed in json format
```json
{
    "trustauthority_url": "https://portal.trustauthority.intel.com",
    "trustauthority_api_url": "https://api.trustauthority.intel.com",
    "trustauthority_api_key": "<trustauthority attestation api key>"
}
```
Save this data in config.json file and invoke `token` command
```sh
sudo ./trustauthority-cli token --config config.json --user-data <base64 encoded userdata>  --no-eventlog
```

### To verify Intel Trust Authority signed token
```sh
./trustauthority-cli verify --config config.json --token <attestation token in JWT format>
```

## License

This source is distributed under the BSD-style license found in the [LICENSE](../LICENSE)
file.
