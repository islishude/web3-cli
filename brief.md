A read-only viewer CLI

```
$ web3-cli -c eth eth_blockNumber
"0x1103524"
$ web3-cli --rpc https://rpc.ankr.com/arbitrum eth_blockNumber | jq -r | xargs printf '%d\n'
117962934

$ web3-cli --chain eth --abi-name erc20 --call-to 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 symbol
[
"USDC"
]
$ web3-cli --chain eth --abi-name erc20 --call-to 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 balanceOf 0xA9D1e08C7793af67e9d92fe308d5697FB81d3E43
[
203540969923675
]
$ # without --abi-name params, it can try to gets abi from explorer
web3-cli --chain eth --call-to 0xdac17f958d2ee523a2206206994597c13d831ec7 totalSupply
```
