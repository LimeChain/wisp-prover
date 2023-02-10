## CRC Prover Server

The Prover is a fork of [Iden3's prover-server](https://github.com/iden3/prover-server). It is a REST API Wrapper
for [go-rapidsnark](https://github.com/iden3/go-rapidsnark)

List of implemented features:

* Generate proof
* Verify proof

### Installation

1. Build prover server:
    ```
    go build ./cmd/prover/prover.go
    ```

2. Edit config file `configs/prover.yaml`

3. Download & unzip compiled circuits: `bash get-circuits.sh`

3. Run prover server:
     ```
    ./prover
    ```

## API

### Generate proof

```
POST /api/v1/proof/generate
Content-Type: application/json
{
  "inputs": {...}, // circuit specific inputs
  "circuit_name": "..." // name of a directory containing circuit_final.zkey, verification_key.json and circuit.wasm files
}
```

## Docker images

Build and run container:

```bash
docker build -t prover-server .
docker run -it -p 8002:8002 prover-server
```

## License

Iden'3 prover-server is part of the iden3 project copyright 2021 0KIMS association and published with GPL-3 license.
Please check the LICENSE file for more details.
