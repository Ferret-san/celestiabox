# Cli tool to upload and download files on Celestia

## Installation
You need to be running go 1.21, then clone the repo, and make sure to run some kind of node (light or brige) using [celestia-node](https://github.com/celestiaorg/celestia-node/tree/main).
Once your node starts running, make sure you keep your auth token and namespace handy.

## Usage
To submit a file, go to the root of this repo and run:

```
go run main.go -mode=submit -file=<path_to_file> -namespace=<namespace> -auth=<auth token>
```

this will return a height and the hex string of the blob commitment:
```
Height:  1873
Commitment string:  e254363061f350c092701ce336ab3e9b493d713883cd3dcf2c271760e9318eb1
```

With the height and commitment string, you can read from the data blob at the given height and write that data into a file by running:

```
go run main.go -mode=read -file=<file> -namespace=<namespace> -auth=<auth token> -commitment=<hex sting of PFB commitment> -height=<height>
```

This will create the file in the repo's directory with whatever file extension you give it. So if you uploaded a `.jpeg`, make sure you write the data into a file with that file extension, and that you tell others the file extension needed, or have them figure it out :)
