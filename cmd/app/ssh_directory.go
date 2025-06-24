package main

import (
	"path/filepath"

	"github.com/charmbracelet/keygen"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/charmbracelet/wishlist"
)

func main() {
	k, err := keygen.New(
		filepath.Join(".wishlist", "server"),
		keygen.WithKeyType(keygen.Ed25519),
	)
	if err != nil {
		log.Fatal("Server keypair", "err", err)
	}
	if !k.KeyPairExists() {
		if err := k.WriteKeys(); err != nil {
			log.Fatal("Server keypair", "err", err)
		}
	}

	// wishlist config
	cfg := &wishlist.Config{
		Port: 2233,
		Factory: func(e wishlist.Endpoint) (*ssh.Server, error) {
			return wish.NewServer(
				wish.WithAddress(e.Address),
				wish.WithHostKeyPEM(k.RawPrivateKey()),
				wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
					return true
				}),
				wish.WithMiddleware(
					append(
						e.Middlewares,
						lm.Middleware(),
						activeterm.Middleware(),
					)...,
				),
			)
		},
		Endpoints: []*wishlist.Endpoint{
			{
				Name:       "Dragon Lair",
				Desc:       "Find out more about my experience and skills",
				Address:    "localhost:5173",
				Link:       wishlist.Link{Name: "", URL: "wwww"},
				RequestTTY: true,
			},
			{
				Name:       "Code Dragon's Elixir",
				Desc:       "My Personal Git Server",
				Address:    "localhost:23231",
				RequestTTY: true,
			},
		},
	}

	// start all the servers
	if err := wishlist.Serve(cfg); err != nil {
		log.Fatal("Serve", "err", err)
	}
}
