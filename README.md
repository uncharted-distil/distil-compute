# distil-compute

![CircleCI](https://circleci.com/gh/unchartedsoftware/distil-compute.svg?style=svg&circle-token=440a62840d79d910d1ad47db988efc0e83861ef3)

## Dependencies

- [Go](https://golang.org/) programming language binaries with the `GOPATH` environment variable specified and `$GOPATH/bin` in your `PATH`.

## Development

#### Clone the repository

```bash
mkdir -p $GOPATH/src/github.com/unchartedsoftware
cd $GOPATH/src/github.com/unchartedsoftware
git clone git@github.com:unchartedsoftware/distil-compute.git
cd distil-compute
git submodule init
git submodule update
```

#### [OPTIONAL] Install protocol buffer compiler

Linux

```bash
curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/protoc-3.3.0-linux-x86_64.zip
unzip protoc-3.3.0-linux-x86_64.zip -d protoc3
sudo mv protoc3/bin/protoc /usr/bin/protoc
```

OSX

```bash
curl -OL https://github.com/google/protobuf/releases/download/v3.3.0/protoc-3.3.0-osx-x86_64.zip
unzip protoc-3.3.0-osx-x86_64.zip -d protoc3
mv protoc3/bin/protoc /usr/bin/protoc
```

#### Install remaining dependencies

```bash
make install
```

#### [OPTIONAL] Generate code

To generate TA3TA2 interface protobuf files if the `pipeline/execute.proto` or `pipeline/ta3ta2-api/*.proto` have changed run:

```bash
make proto
```

To regenerate the PANDAS dataframe parser if the `primitive/compute/result/complex_field.peg` file is changed, run:

```bash
make peg
```
