# WEB3-CLI

[![test](https://github.com/islishude/web3-cli/actions/workflows/test.yaml/badge.svg)](https://github.com/islishude/web3-cli/actions/workflows/test.yaml)

## Install

You can download the binary in the latest [release page](https://github.com/islishude/web3-cli/releases/latest).

or you have docker

```
docker pull ghcr.io/islishude/web3-cli
```

or if you have golang installed

```sh
# install with the laest tag
go install github.com/islishude/web3-cli@latest
# install with the latest commit
go install github.com/islishude/web3-cli@main
```

## Usage

**Do a simple jsonrpc call**

the default rpc endpoint is your local: `http://localhost:8545`

your command is like `web3-cli jsonrpc_method [jsonrpc_param...]`

```console
$ web3-cli web3_clientVersion
"Geth/v1.12.1-unstable-60070fe5-20230805/darwin-amd64/go1.20.7"
$ web3-cli eth_getBlockByNumber 1 false
{
    "baseFeePerGas": "0x342770c0",
    "difficulty": "0x0",
    "extraData": "0xd983010c01846765746888676f312e32302e378664617277696e",
    "gasLimit": "0xafa5bd",
    "gasUsed": "0x5208",
    "hash": "0xf5c52795b6fb4b69601e8613c436e07d964b33fb97b5b2b8670a106e0617bf6e",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "miner": "0x0000000000000000000000000000000000000000",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000000",
    "number": "0x1",
    "parentHash": "0xd7bcb282488afadddab1b33e1d5906d1cf367bfc6249b6eaf140961e986621c2",
    "receiptsRoot": "0xf78dfb743fbd92ade140711c8bbc542b5e307f0ab7984eff35d751969fe57efa",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "size": "0x2ad",
    "stateRoot": "0x05b0dc6527b6d706aa6400866014f8f0bc4996ac8f7b9c4d7b4801fccfad139f",
    "timestamp": "0x64fc7d7a",
    "totalDifficulty": "0x0",
    "transactions": [
        "0xfae530d45ec70be05c68aaa0628588f0055f4aae360a8c090c2591e8691a47ca"
    ],
    "transactionsRoot": "0x8b2e1c5c5426fb942d176fac9186a05911c41867709a9fa43b05e2051a64f712",
    "uncles": [],
    "withdrawals": [],
    "withdrawalsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
}
```

web3-cli supports built-in chain config as well.

```console
$ web3-cli --chain eth eth_blockNumber
"0x1142cb9"
$ web3-cli --chain bsc eth_chainId
"0x38"
```

You can run `web3-cli chains` to get all built-in chain configs.

Of course, you can also use your own rpc endpoint

```console
$ web3-cli --rpc https://rpc.ankr.com/arbitrum eth_chainId
"0xa4b1"
```

**Call a contract with human-readable params**

```console
$ export USDT_TOKEN=0xdac17f958d2ee523a2206206994597c13d831ec7
$ web3-cli --chain eth --abi-name erc20 --call-to $USDT_TOKEN symbol
[
    "USDT"
]
$ web3-cli --chain eth --call-to $USDT_TOKEN balanceOf 0x0000000000000000000000000000000000000000
[
    10123456
]
```

the `--abi-name` could be a built-in abi name, you can `web3-cli abis` to get the built-in abi list.

and it also could be a url and file path.

```
$ web3-cli --chain eth --abi-name https://http-server/abi/abi.json
$ web3-cli --chain eth --abi-name local/path/to/abi.json
```

one more, if you don't provide `--abi-name`, web3-cli can fetch the abi from explorer api automatically

```
$ web3-cli --chain eth --call-to $USDT_TOKEN getOwner
[
    "0xc6cde7c39eb2f0f0095f41570af89efc2c1ea828"
]
```

you may need to provide your explorer api endpoint as well if you provide a custom rpc

```console
$ # not required if use a built-in chain
$ web3-cli --chain eth --rpc https://my-own-rpc.com
$ web3-cli --rpc https://my-own-rpc.com --explorer-api https://my-custom-explorer.com/api
```

What about a complex abi type parameter, like array and tuple?

**bytes**

```solidity
    function logbyt(bytes memory _log) public returns (bytes memory log){
        counter++;
        return _log;
    }
```

you must use a hex string

```
0x776562332d636c6920697320736f20636f6f6c21
```

**array and slice**

```solidity
    function add(uint256[] calldata items) public pure returns (uint256 sum) {
        for (uint256 i = 0; i < items.length; i++){
            sum += items[i];
        }
    }
```

json array is valid

```json
["0x1", 100]
```

**tuple(struct)**

```solidity
    struct Payment {
        address payable to;
        uint256 value;
    }

    function transfer(Payment calldata item) external payable returns (bool success) {
       return item.to.send(msg.value);
    }
```

json array is valid (you can use it in Remix)

```json
["0x0000000000000000000000000000000000000000", "0x1"]
```

json object is valid as well!

the key is the tuple name, for most cases, you can just use the field name.

```json
{ "to": "0x0000000000000000000000000000000000000000", "value": "0x1" }
```

## Contribute

**How to add your chain to built-in list?**

Add it to `internal/chains/list.go` file.

for example:

```go
var eth = &Chain{
	Name:     "eth",
	Id:       1,
	Endpoint: "https://cloudflare-eth.com",
	Explorer: "https://api.etherscan.io/api",
	Alias:    []string{"eth-mainnet", "mainnet"},
}
```

**How to add an ABI to built-in list?**

Add a valid json file to `internal/abis/abi` directory.
