A read-only viewer CLI

```
$ web3-cli -c eth eth_blockNumber
"0x1103524"
$ web3-cli --rpc https://rpc.ankr.com/arbitrum eth_blockNumber | jq -r | xargs printf '%d\n'
117962934
$ export USDT_ADDRESS=0xdac17f958d2ee523a2206206994597c13d831ec7
$ web3-cli --chain eth --abi-name erc20 --call-to $USDT_ADDRESS symbol
[
    "USDT"
]
$ # fetch abi from explorer if you don't provide it
$ web3-cli --chain eth --call-to $USDT_ADDRESS getOwner
[
    "0xc6cde7c39eb2f0f0095f41570af89efc2c1ea828"
]
```
