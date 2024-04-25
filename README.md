# myyonseinft-backend

## Prerequisite
```
You have to install go
```

## Clone repository
```bash
git clone https://github.com/guboc11/myyonsenft-backend.git

cd myyonseinft-backend/backend
```

## Download & install go packages
```bash
go get
```

## Add .env
```bash
vi .env
```
#### .env
```
PRIVATE_KEY={your_ethereum_holeskey_network_wallet_private_key}
SENDER_ADDRESS={your_wallet_address}
CONTRACT_ABI={contract_abi}
CONTRACT_ADDRESS={contract_address}
```

## Run API
```bash
// check if localhost:8080 is open and available
go run main.go
```

## Query API
### /mint
```bash
curl -X POST "http://localhost:8080/mint?address={CONTRACT_ADDRESS}"
```
### /balanceOf
```bash
curl "http://localhost:8080/balanceOf?address={CONTRACT_ADDRESS}"
```
### /history
```bash
curl "http://localhost:8080/history?address={CONTRACT_ADDRESS}"
```