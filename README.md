# CRC Prover Server

[![Integration Test](https://github.com/LimeChain/crc-prover/actions/workflows/ci.yaml/badge.svg?branch=master)](https://github.com/LimeChain/crc-prover/actions/workflows/ci.yaml)

The Prover is a fork of [Iden3's prover-server](https://github.com/iden3/prover-server). It is a REST API Wrapper
for [go-rapidsnark](https://github.com/iden3/go-rapidsnark)

List of implemented features:

* Generate proof
* Verify proof

### Installation

1. (Optional) Create / edit your config file. Defaults to `configs/dev.yaml`.
2. Prepare compiled circuits, zkey and verification key.
    1. Option 1 (Development): Use `multiplier`:
          ```
          mkdir ./circuits && cp -R .github/workflows/test-circuit/multiplier ./circuits/multiplier
          ```
    2. Option 2: CRC circuits (`ssz2Poseidon` and `blsHeaderVerify`):
          ```
          bash get-circuits.sh
          ```
       Hardware requirements are `256GB RAM`, `32-core CPU` and `1 TB SSD`
3. Build the image
    ```
    docker build -t prover-server .
    ```
4. Run Prover
   ```
   docker run -it -p 8000:8000 prover-server
   ```
   If you want to use config, different from the default `dev` one you must pass it as an environmental
   variable `CONFIG={config}`

## API

### Generate proof

```
POST /api/v1/proof/generate
Content-Type: application/json
{
    "circuit": "multiplier", // name of the requested circuit as specified in the config
    "inputs": {...} // circuit specific inputs
}
```

## License

Iden'3 prover-server is part of the iden3 project copyright 2021 0KIMS association and published with GPL-3 license.
Please check the LICENSE file for more details.
