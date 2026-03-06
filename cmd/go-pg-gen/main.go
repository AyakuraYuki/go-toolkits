package main

import "github.com/AyakuraYuki/go-toolkits/cmd/go-pg-gen/internal/gen"

func main() {
	if err := gen.Cmd.Execute(); err != nil {
		gen.Log.Fatal().Err(err).Msg("oops, something went wrong")
	}
}
