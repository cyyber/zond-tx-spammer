<h1>Transaction Spammer</h1>

tx-spammer is a simple tool that can be used to generate various types of random transactions for qrl testnets.

tx-spammer can be used for stress testing (flooding the network with thousands of transactions) or to have a continuous amount of transactions over long time for testing purposes.

## Build

You can use this tool via pre-build docker images: [theqrl/qrl-tx-spammer](https://hub.docker.com/r/theqrl/qrl-tx-spammer)

Or build it yourself:

```
git clone https://github.com/theQRL/qrl-tx-spammer.git
cd tx-spammer
make
./bin/tx-spammer
```