module github.com/grandper/go-screenplay-nats/v2

go 1.27.2

require (
	github.com/grandper/go-screenplay/v2 v2.0.0
	github.com/nats-io/nats-server/v2 v2.12.1
	github.com/nats-io/nats.go v1.47.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/antithesishq/antithesis-sdk-go v0.4.3-default-no-op // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-tpm v0.9.6 // indirect
	github.com/grandper/go-monad v1.0.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/nats-io/jwt/v2 v2.8.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// go-screenplay v2.0.0 is not tagged yet: build against the release branch checked out next to this repository.
replace github.com/grandper/go-screenplay/v2 => ../go-screenplay
