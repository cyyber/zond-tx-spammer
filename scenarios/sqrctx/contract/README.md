## Commands

```
hypc --standard-json ./TestToken.input.json | jq '.contracts["contracts/TestToken.hyp"].TestToken.abi' > TestToken.abi
hypc --standard-json ./TestToken.input.json | jq -r '.contracts["contracts/TestToken.hyp"].TestToken.zvm.bytecode.object' > TestToken.bin
abigen --bin=./TestToken.bin --abi=./TestToken.abi --pkg=contract --out=TestToken.go
```