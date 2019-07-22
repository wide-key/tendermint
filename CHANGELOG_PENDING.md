## v0.32.2

\*\*

Special thanks to external contributors on this release:

Friendly reminder, we have a [bug bounty
program](https://hackerone.com/tendermint).

### BREAKING CHANGES:

- CLI/RPC/Config
  - [rpc] `/block_results` response format updated (see RPC docs for details)
    ```
    {
      "jsonrpc": "2.0",
      "id": "",
      "result": {
        "height": "2109",
        "txs_results": null,
        "begin_block_events": null,
        "end_block_events": null,
        "validator_updates": null,
        "consensus_param_updates": null
      }
    }

- Apps

- Go API
  - [libs] \#3811 Remove `db` from libs in favor of `https://github.com/tendermint/tm-cmn`

### FEATURES:

### IMPROVEMENTS:

- [abci] \#3809 Recover from application panics in `server/socket_server.go` to allow socket cleanup (@ruseinov)
- [rpc] \#3818 Make `max_body_bytes` and `max_header_bytes` configurable
- [p2p] \#3664 p2p/conn: reuse buffer when write/read from secret connection

### BUG FIXES:

