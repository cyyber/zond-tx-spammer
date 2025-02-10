<img align="left" src="./.github/resources/goomy.png" width="75">
<h1>Transaction Spammer</h1>

tx-spammer is a simple tool that can be used to generate various types of random transactions for zond testnets.

tx-spammer can be used for stress testing (flooding the network with thousands of transactions) or to have a continuous amount of transactions over long time for testing purposes.

## Build

You can use this tool via pre-build docker images: [theqrl/tx-spammer](https://hub.docker.com/r/theqrl/tx-spammer)

Or build it yourself:

```
git clone https://github.com/theQRL/tx-spammer.git
cd tx-spammer
make
./bin/tx-spammer
```

## Usage

### `tx-spammer`
`tx-spammer` is a tool for sending mass transactions.

```
Usage of tx-spammer:
Required:
  -p, --privkey string        The private key of the wallet to send funds from.
  
  -h, --rpchost string        The RPC host to send transactions to (multiple allowed).
      --rpchost-file string   File with a list of RPC hosts to send transactions to.
      
Optional:
  -s, --seed string           The child wallet seed.
  -v, --verbose               Run the tool with verbose output.
```

The tool provides multiple scenarios, that focus on different aspects of transactions. One of the scenarios must be selected to run the tool:

#### `tx-spammer eoatx`

The `eoatx` scenario sends normal dynamic fee transactions.

```
Usage of ./bin/tx-spammer eoatx:
      --amount uint            Transfer amount per transaction (in gwei) (default 20)
      --basefee uint           Max fee per gas to use in transfer transactions (in gwei) (default 20)
  -c, --count uint             Total number of transfer transactions to send
      --data string            Transaction call data to send
      --gaslimit uint          Gas limit to use in transactions (default 21000)
      --max-pending uint       Maximum number of pending transactions
      --max-wallets uint       Maximum number of child wallets to use
  -p, --privkey string         The private key of the wallet to send funds from.
      --random-amount          Use random amounts for transactions (with --amount as limit)
      --random-target          Use random to addresses for transactions
      --rebroadcast uint       Number of seconds to wait before re-broadcasting a transaction (default 120)
      --refill-amount uint     Amount of ETH to fund/refill each child wallet with. (default 5)
      --refill-balance uint    Min amount of ETH each child wallet should hold before refilling. (default 2)
      --refill-interval uint   Interval for child wallet rbalance check and refilling if needed (in sec). (default 300)
  -h, --rpchost stringArray    The RPC host to send transactions to.
      --rpchost-file string    File with a list of RPC hosts to send transactions to.
  -s, --seed string            The child wallet seed.
  -t, --throughput uint        Number of transfer transactions to send per slot
      --tipfee uint            Max tip per gas to use in transfer transactions (in gwei) (default 2)
      --trace                  Run the script with tracing output
  -v, --verbose                Run the script with verbose output
```

#### `tx-spammer erctx`

The `erctx` scenario deploys an ERC20 contract and performs token transfers.

```
Usage of ./bin/tx-spammer erctx:
      --amount uint            Transfer amount per transaction (in gwei) (default 20)
      --basefee uint           Max fee per gas to use in transfer transactions (in gwei) (default 20)
  -c, --count uint             Total number of transfer transactions to send
      --max-pending uint       Maximum number of pending transactions
      --max-wallets uint       Maximum number of child wallets to use
  -p, --privkey string         The private key of the wallet to send funds from.
      --random-amount          Use random amounts for transactions (with --amount as limit)
      --random-target          Use random to addresses for transactions
      --rebroadcast uint       Number of seconds to wait before re-broadcasting a transaction (default 120)
      --refill-amount uint     Amount of ETH to fund/refill each child wallet with. (default 5)
      --refill-balance uint    Min amount of ETH each child wallet should hold before refilling. (default 2)
      --refill-interval uint   Interval for child wallet rbalance check and refilling if needed (in sec). (default 300)
  -h, --rpchost stringArray    The RPC host to send transactions to.
      --rpchost-file string    File with a list of RPC hosts to send transactions to.
  -s, --seed string            The child wallet seed.
  -t, --throughput uint        Number of transfer transactions to send per slot
      --tipfee uint            Max tip per gas to use in transfer transactions (in gwei) (default 2)
      --trace                  Run the script with tracing output
  -v, --verbose                Run the script with verbose output
```

#### `tx-spammer deploytx`

The `deploytx` scenario sends contract deployment transactions.

```
Usage of ./bin/tx-spammer deploytx:
      --basefee uint            Max fee per gas to use in deployment transactions (in gwei) (default 20)
      --bytecodes string        Bytecodes to deploy (, separated list of hex bytecodes)
      --bytecodes-file string   File with bytecodes to deploy (list with hex bytecodes)
  -c, --count uint              Total number of deployment transactions to send
      --gaslimit uint           Gas limit to use in deployment transactions (in gwei) (default 1000000)
      --max-pending uint        Maximum number of pending transactions
      --max-wallets uint        Maximum number of child wallets to use
  -p, --privkey string          The private key of the wallet to send funds from.
      --rebroadcast uint        Number of seconds to wait before re-broadcasting a transaction (default 120)
      --refill-amount uint      Amount of ETH to fund/refill each child wallet with. (default 5)
      --refill-balance uint     Min amount of ETH each child wallet should hold before refilling. (default 2)
      --refill-interval uint    Interval for child wallet rbalance check and refilling if needed (in sec). (default 300)
  -h, --rpchost stringArray     The RPC host to send transactions to.
      --rpchost-file string     File with a list of RPC hosts to send transactions to.
  -s, --seed string             The child wallet seed.
  -t, --throughput uint         Number of deployment transactions to send per slot
      --tipfee uint             Max tip per gas to use in deployment transactions (in gwei) (default 2)
      --trace                   Run the script with tracing output
  -v, --verbose                 Run the script with verbose output
```

### `tx-spammer gasburnertx`

The `gasburnertx` scenario sends out transactions with a configurable amount of gas units. Note that the estimated gas units is not 100% accurate.

```
Usage of tx-spammer gasburnertx:
Required (at least one of):
  -c, --count uint            Total number of gasburner transactions to send
  -t, --throughput uint       Number of gasburner transactions to send per slot
  
Optional:
      --basefee uint             Max fee per gas to use in gasburner transactions (in gwei) (default 20)
      --gas-units-to-burn uint   The number of gas units for each tx to cost (default 2000000)
      --max-pending uint         Maximum number of pending transactions
      --max-wallets uint         Maximum number of child wallets to use
  -p, --privkey string           The private key of the wallet to send funds from.
      --rebroadcast uint         Number of seconds to wait before re-broadcasting a transaction (default 120)
  -h, --rpchost stringArray      The RPC host to send transactions to.
      --rpchost-file string      File with a list of RPC hosts to send transactions to.
  -s, --seed string              The child wallet seed.
      --tipfee uint              Max tip per gas to use in gasburner transactions (in gwei) (default 2)
      --trace                    Run the script with tracing output
  -v, --verbose                  Run the script with verbose output
```

