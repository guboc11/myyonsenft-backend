# take-home-2023-guboc11

## Prerequisite
```
You have to install go
```

## Clone repository
```bash
git clone https://github.com/planetarium/take-home-2023-guboc11.git

cd take-home-2023-guboc11/backend
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
DELIGATOR_ADDRESS={your_wallet_address}
CONTRACT_ABI={contract_abi}
CONTRACT_ADDRESS={contract_address}
```

## Run api
```bash
// check if localhost:8080 is open and available
go run main.go
```

## Query api
### /mint
```bash
curl -X POST "http://localhost:8080/mint?address=0xF8c847Fc824B441f0b4D9641371e6eD3f56CF145"
```
### /balanceOf
```bash
curl "http://localhost:8080/balanceOf?address=0xF8c847Fc824B441f0b4D9641371e6eD3f56CF145"
```
### /history
```bash
curl "http://localhost:8080/history?address=0xF8c847Fc824B441f0b4D9641371e6eD3f56CF145"
```