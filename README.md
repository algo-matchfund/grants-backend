# Grants backend

## Prerequisite
To be able to fund, you must have [AlgoSigner](chrome://extensions/?id=kpfnocihfdjkfacalpieicmnmkclbhji) installed and have an existing account. Use TestNet for development purposes. Currently, the extension only available in Chrome, which means funding can also only be done in Chrome.

To get more Algos/USDC in TestNet: https://dispenser.testnet.aws.algodev.network/
## Run locally

```
git submodule init && git submodule update
mkdir gen
./generate.sh
go build -a -o ./build/grants-backend ./cmd/grants-backend/main.go
./build/grants-backend -config ci/config.yml
```

## Run in Docker

To run locally with docker-compose initialize all git submodules, and then run docker-compose up.

```bash
git submodule init && git submodule update
mkdir gen
./generate.sh
```
Before starting the containers:
- provide `ALGOD_TOKEN` for your own Algorand account in the `grants-frontend` section of `docker-compose.yml`.
- In `/etc/hosts`, add in this line to allow these URLs to be redirected to localhost:

`127.0.0.1 keycloak grants-backend grants-frontend`


Finally, run `docker-compose up`

This will start five docker containers:
* Postgres Database
* Grants Backend
* Grants Frontend
* KeyCloak
* API Documentation

Keycloak uses a public/private key pair to sign JSON Web Tokens. The default development configuration has the key pair seeded in both the keycloak realm and the backend configuration for ease of use.

The frontend is reachable at https://grants-frontend
For documentation on the backend API visit http://localhost:8091
The API itself is reachable at https://grants-backend:444

You may have to accept certs for all three of the URLs above.

## Keycloak Key Configuration
Keycloak tokens are verified by checking if they are unexpired and if they were signed by the keycloak server. The public key of the keycloak server must be added to the configuration file in PEM format.

The public key for the keycloak server can be obtained by visiting `http://localhost:8080/auth/realms/grants`.

It must then be placed in `ci/config.yml` in the following format: `"-----BEGIN PUBLIC KEY-----\n<public key>\n-----END PUBLIC KEY-----"`.

## Updating smart contract
`approval.teal` and `clear.teal` in the root directory are the smart contract programs. To update them:

1. `cd grants-smart-contract`
2. Set up venv (one time): `python3 -m venv venv`
3. Active venv: `. venv/bin/activate`
4. Install dependencies (one time): `pip install -r requirements.txt`
5. Compile and move: `python3 project/contracts.py && mv approval.teal clear.teal ..`