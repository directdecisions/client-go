# Direct Decisions API v1 Go client

[![Go](https://github.com/directdecisions/client-go/workflows/Go/badge.svg)](https://github.com/directdecisions/client-go/actions)
[![PkgGoDev](https://pkg.go.dev/badge/directdecisions.com/directdecisions)](https://pkg.go.dev/directdecisions.com/directdecisions)
[![NewReleases](https://newreleases.io/badge.svg)](https://newreleases.io/github/directdecisions/client-go)

Package directdecisions is a Go client library for accessing the [Direct Decisions](https://directdecisions.com) v1 API.

You can view the client API docs here: [https://pkg.go.dev/directdecisions.com/directdecisions](https://pkg.go.dev/directdecisions.com/directdecisions)

You can view Direct Decisions API v1 docs here: [https://api.directdecisions.com/v1](https://api.directdecisions.com/v1)

## Installation

This package requires Go 1.20 version or later.

Run `go get directdecisions.com/directdecisions` from command line.

## Usage

```go
import "directdecisions.com/directdecisions"
```

Create a new Client, then use the exposed services to access different parts of the API.

## Features

This client implements all Direct API features.

- Create votings
- Retrieve voting information
- Set voting choices
- Delete votings
- Vote with a ballot
- Unvote
- Get submitted ballot
- Calculate results

## Examples

To run a voting:

```go
package main

import (
    "context"
    "log"

    "directdecisions.com/directdecisions"
)

func main() {
	client := directdecisions.NewClient("my-api-key", nil)

	ctx := context.Background()

	v, err := client.Votings.Create(ctx, []string{"Margarita", "Pepperoni", "Capricciosa"})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created voting with ID %s", v.ID)

	if _, err := client.Votings.Vote(ctx, v.ID, "Leonardo", map[string]int{
		"Pepperoni": 1,
		"Margarita": 2,
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("Leonardo Voted for Pepperoni in voting with ID %s", v.ID)

	if _, err := client.Votings.Vote(ctx, v.ID, "Michelangelo", map[string]int{
		"Capricciosa": 1,
		"Margarita":   2,
		"Pepperoni":   2,
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("Michelangelo Voted for Capricciosa in voting with ID %s", v.ID)

	results, tie, err := client.Votings.Results(ctx, v.ID)
	if err != nil {
		log.Fatal(err)
	}

	if tie {
		log.Printf("Voting with ID %s is tied", v.ID)
	} else {
		log.Printf("Voting with ID %s is not tied", v.ID)
	}

	log.Printf("Results for voting with ID %s: %v", v.ID, results)
}
```

Error handling with Go 1.20 multiple errors wrapped:

```go
package main

import (
    "context"
    "log"

    "directdecisions.com/directdecisions"
)

func main() {
	client := directdecisions.NewClient("my-api-key", nil)

	ctx := context.Background()

	_, err := client.Votings.Create(ctx, []string{"Margarita", "Pepperoni", "Capricciosa"})
	if err != nil {
		if errors.Is(err, directdecisions.ErrHTTPStatusUnauthorized) {
			log.Fatal("Invalid API key")
		}
		if errors.Is(err, directdecisions.ErrChoiceTooLong) {
			log.Fatal("Some of the choices are too long")
		}
		// ...
		log.Fatal(err)
	}
}
```

## Versioning

Each version of the client is tagged and the version is updated accordingly.

This package uses Go modules.

To see the list of past versions, run `git tag`.

## Contributing

We love pull requests! Please see the [contribution guidelines](CONTRIBUTING.md).

## License

This library is distributed under the BSD-style license found in the [LICENSE](LICENSE) file.
