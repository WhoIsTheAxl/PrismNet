package bedrock

import (
	"context"
	"fmt"
	"net/http"

	"github.com/df-mc/go-playfab/v2"
	"github.com/df-mc/go-xsapi/v2"
	"github.com/lactyy2/pia/signaling/messaging"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sandertv/gophertunnel/minecraft/realms"
	"github.com/sandertv/gophertunnel/minecraft/service"
	"golang.org/x/oauth2"
)



func createServiceTokenSource(
	ctx context.Context,
	tokenSource oauth2.TokenSource,
) (service.TokenSource, error) {

	ctx = auth.WithContextClient(ctx, http.DefaultClient)

	xboxClient, err := xsapi.ClientConfig{
		HTTPClient: http.DefaultClient,
		RTAMode:    xsapi.RTADisabled,
	}.New(
		ctx,
		auth.ContextSession(ctx, tokenSource),
	)
	if err != nil {
		return nil, fmt.Errorf("login to Xbox Live: %w", err)
	}

	discovery, err := service.Default(ctx)
	if err != nil {
		xboxClient.Close()
		return nil, fmt.Errorf("discover Minecraft services: %w", err)
	}

	env := new(service.AuthorizationEnvironment)

	if err := discovery.Environment(env); err != nil {
		xboxClient.Close()
		return nil, fmt.Errorf("get authorization environment: %w", err)
	}

	playFabClient, err := playfab.LoginWithXbox(
		ctx,
		env.PlayFabTitleID,
		xboxClient,
		playfab.ClientConfig{
			HTTPClient:    http.DefaultClient,
			CreateAccount: true,
		},
	)
	if err != nil {
		xboxClient.Close()
		return nil, fmt.Errorf("login to PlayFab: %w", err)
	}

	return env.TokenSource(
		playFabClient,
		service.TokenConfig{
			Device: service.DeviceConfig{
				TreatmentOverrides: []string{
					"mc-signaling-usewebsockets",
				},
			},
		},
	), nil
}

func realmConnect(
	ctx context.Context,
	tokenSource oauth2.TokenSource,
	inviteCode string,
) (*minecraft.Conn, error) {

	fmt.Println("[1] Resolving Realm...")

	realmClient := realms.NewClient(
		tokenSource,
		http.DefaultClient,
	)

	realm, err := realmClient.Realm(ctx, inviteCode)
	if err != nil {
		return nil, fmt.Errorf("resolve Realm: %w", err)
	}

	fmt.Println("[2] Getting Realm network address...")

	address, err := realm.Address(ctx)
	if err != nil {
		return nil, fmt.Errorf("get Realm address: %w", err)
	}

	fmt.Printf("    Network protocol: %s\n", address.NetworkProtocol)
	fmt.Printf("    Network ID:       %s\n", address.Address)

	if address.NetworkProtocol != realms.NetworkProtocolNetherNetJSONRPC {
		return nil, fmt.Errorf(
			"unsupported Realm network protocol: %s",
			address.NetworkProtocol,
		)
	}

	fmt.Println("[3] Creating Minecraft service token source...")

	serviceSource, err := createServiceTokenSource(
		ctx,
		tokenSource,
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("[4] Connecting to Player Messaging...")

	signaling, err := messaging.Dialer{}.DialContext(
		ctx,
		serviceSource,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"connect Player Messaging: %w",
			err,
		)
	}

	fmt.Printf(
		"    Local Player Messaging ID: %s\n",
		signaling.PlayerMessagingID(),
	)

	fmt.Println("[5] Connecting through NetherNet...")

	network := minecraft.NetherNet{
		Signaling: signaling,
	}

	dialer := minecraft.Dialer{
    TokenSource: tokenSource,
}

	conn, err := dialer.DialContextNetwork(
	    ctx,
	    network,
	    address.Address,
	)
	if err != nil {
	    return nil, err
	}
	
	if err := conn.DoSpawn(); err != nil {
        conn.Close()
        return nil, err
    }
	
	return conn, nil
}