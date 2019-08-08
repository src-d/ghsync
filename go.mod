module github.com/src-d/ghsync

go 1.12

//Use github.com/mcuadros/go-kallax fork, ghsync branch
replace gopkg.in/src-d/go-kallax.v1 => github.com/mcuadros/go-kallax v1.3.6-0.20190516223806-dc0ad3de8cf0

// kallax does not work with modules, so this replacement is needed to be able
// to update the kallax models. Without this replacement google/go-github would
// point to a very old release
// replace github.com/google/go-github => github.com/google/go-github/v25
replace github.com/google/go-github => github.com/google/go-github/v25 v25.1.1

// src-d/go-queue requires master branch of go.uuid
// require github.com/satori/go.uuid master
require github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b // indirect

require (
	github.com/Masterminds/squirrel v1.1.0 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang-migrate/migrate/v4 v4.4.0
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-github/v25 v25.1.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/jpillora/backoff v0.0.0-20180909062703-3050d21c67d7 // indirect
	github.com/kami-zh/go-capturer v0.0.0-20171211120116-e492ea43421d // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/lib/pq v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/src-d/envconfig v1.0.0 // indirect
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/x-cray/logrus-prefixed-formatter v0.5.2 // indirect
	golang.org/x/crypto v0.0.0-20190530122614-20be4c3c3ed5 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/src-d/go-cli.v0 v0.0.0-20190422143124-3a646154da79
	gopkg.in/src-d/go-errors.v0 v0.1.0 // indirect
	gopkg.in/src-d/go-errors.v1 v1.0.0 // indirect
	gopkg.in/src-d/go-kallax.v1 v1.3.5
	gopkg.in/src-d/go-log.v1 v1.0.2
	gopkg.in/src-d/go-queue.v1 v1.0.6
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1 // indirect
)
