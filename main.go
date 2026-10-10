package main

import (
	"context"
	"log"

	"github.com/WhoIsTheAxl/Prismnet/bedrock"
)

func main() {
	auth, err := bedrock.Auth(bedrock.AuthConfig{
		TokenPath: "token.json",
	})
	if err != nil {
		log.Fatal(err)
	}

	realm := bedrock.Realm(
		"3WAYyAhQYgAihKg",
		auth.TokenSource,
	)

	client, err := realm.Connect(context.Background())
	// client, err := bedrock.Connect(context.Background(), "zeqa.net:19132", auth)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	session := client.Session()

	// ─────────────────────────────
	// Session events
	// ─────────────────────────────

	session.OnPacket(func(packet any) {
		// Every packet received from the server.
	})

	session.OnDisconnect(func(err error) {
		log.Println("Disconnected:", err)
	})

	session.OnPlayer(func(player *bedrock.Player) {
		log.Println("Player joined:", player.Name)
	})

	// ─────────────────────────────
	// Commands
	// ─────────────────────────────

	cmd := session.Command()

	cmd.Run("/say Hello")
	cmd.Whisper("@a", "Hello")

	// ─────────────────────────────
	// Local player
	// ─────────────────────────────

	player := session.Player()

	log.Println("Username:", player.Name)
	log.Println("UUID:", player.UUID)
	log.Println("Position:", player.Position)

	player.Move(
		bedrock.Position{
			X: 10,
			Y: 64,
			Z: 20,
		},
	)

	player.LookAt(
		bedrock.Position{
			X: 15,
			Y: 64,
			Z: 25,
		},
	)

	player.Punch()
	player.StartBreak()
	player.StopBreak()

	// ─────────────────────────────
	// Other players
	// ─────────────────────────────

	for _, player := range session.Players() {
		log.Println(
			"Player:",
			player.Name,
			player.Position,
		)
	}

	select {}
}