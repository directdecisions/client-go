// Copyright (c) 2022, Direct Decisions Go client AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package directdecisions_test

import (
	"context"
	"errors"
	"log"

	"directdecisions.com/directdecisions"
)

func ExampleNewClient() {
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

func Example_errorHandling() {
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
