# ghsync

## Kallax Models

In order to update the kallax models, place this project in `$GOPATH/src/github.com/src-d/ghsync`.

Then follow these steps:

```shell
# Make sure you are not using modules
unset GO111MODULE

# Get kallax, replace it with mcuadros fork, branch ghsync
go get -u gopkg.in/src-d/go-kallax.v1/...
cd $GOPATH/src/gopkg.in/src-d/go-kallax.v1
git remote add mcuadros git@github.com:mcuadros/go-kallax.git
git fetch --all
git checkout -b ghsync mcuadros/ghsync

# Build kallax
rm $GOPATH/bin/kallax
go get -u github.com/golang-migrate/migrate
cd generator/cli/kallax
go install

# Make sure the $GOPATH/bin is in your path, if not run
export PATH=$GOPATH/bin:$PATH

# Back to ghsync, create the dependencies vendor folder
cd $GOPATH/src/github.com/src-d/ghsync
GO111MODULE=on go mod vendor

# Run kallax generation
go generate ./...
```
