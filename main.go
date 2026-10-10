package main

import (
 "context"
 "log"
 "github.com/WhoIsTheAxl/PrismNet/bedrock"
 "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Example: authenticate and connect to a Bedrock server.
// Set an address you control before running this example.
func main() {
 source, err := bedrock.Auth(bedrock.AuthConfig{TokenPath:"token.json"})
 if err != nil { log.Fatal(err) }
 client, err := bedrock.Connect(context.Background(),"localhost:19132",source.TokenSource)
 if err != nil { log.Fatal(err) }
 defer client.Close()
 session := client.Session()
 session.OnPacket(func(p packet.Packet) { log.Printf("received %T",p) })
 session.OnDisconnect(func(err error) { log.Printf("disconnected: %v",err) })
 if err:=session.Start();err!=nil {log.Fatal(err)}
 <-session.Done()
}
