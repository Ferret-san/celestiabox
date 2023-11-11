# Cli tool to upload and download files on Celestia

## Installation

You need to be running go 1.21, then clone the repo, and make sure to run some kind of node (light or brige) using [celestia-node](https://github.com/celestiaorg/celestia-node/tree/main) with [state access (using the `--core.ip <rpc>` flag)](https://docs.celestia.org/developers/node-tutorial#connect-to-a-public-core-endpoint).

Once your node starts running, make sure you [keep your auth token](https://docs.celestia.org/developers/node-tutorial#auth-token) and [namespace](https://celestiaorg.github.io/celestia-app/specs/namespace.html) handy.

## Usage

To submit a file, go to the root of this repo and run:

```sh
go run main.go -mode=submit -file=<path_to_file> \
    -namespace=<namespace> -auth=<auth token>
```

This will return a height and the hex string of the blob commitment:

```sh
Successfully read 315987 bytes from <path_to_file>
Succesfully submitted blob to Celestia
Height:  1873
Commitment string:  e254363061f350c092701ce336ab3e9b493d713883cd3dcf2c271760e9318eb1
```

An example command looks like this:

```sh
go run main.go -mode=submit \
    -file=$HOME/celestiabox/<path_to_file> \
    -namespace=0000000000004a6f7368 -auth=$AUTH_TOKEN
```

With the height and commitment string, you can read from the data blob at the given height and write that data into a file by running:

```sh
go run main.go -mode=read -file=<file> \
    -namespace=<namespace> -auth=<auth token> \
    -commitment=<hex sting of PFB commitment> -height=<height>
```

This will create the file in the repo's directory with whatever file extension you give it. So if you uploaded a `.jpeg`, make sure you write the data into a file with that file extension, and that you tell others the file extension needed, or have them figure it out :)

Here's an example!

```sh
go run main.go -mode=read -file=diego.png \
    -namespace=0000000000004a6f7368 -auth=$AUTH_TOKEN \
    -commitment=15c4d44e62c098634fa2ccb57d1c2e690cd6e42b1b61d73c41150e57ec193658 \
    -height=82764
```

The output will look similar to below:

```sh
Requesting data from Celestia namespace 0000000000004a6f7368 commitment �NbcO̵}.i
      �+�AW�6X height 82764
Succesfully fetched data from Celestia namespace 0000000000004a6f7368 height 82764 commitment �NbcO̵}.i
                            �+�AW�6X
File written successfully!
```
