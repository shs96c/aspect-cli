package client

import (
	goplugin "github.com/hashicorp/go-plugin"

	sdk_config "aspect.build/cli/pkg/plugin/sdk/v1alpha2/config"
	"aspect.build/cli/pkg/plugin/sdk/v1alpha2/plugin"
)

// ClientFactory hides the call to goplugin.NewClient() and conversion to InternalPlugin.
type Factory interface {
	New(*goplugin.ClientConfig) (*InternalPlugin, error)
}

type clientFactory struct{}

func NewFactory() Factory {
	return &clientFactory{}
}

// New calls the goplugin.NewClient with the given config.
func (*clientFactory) New(clientConfig *goplugin.ClientConfig) (*InternalPlugin, error) {
	var internalPlugin *InternalPlugin

	client := goplugin.NewClient(clientConfig)

	rpcClient, err := client.Client()
	if err != nil {
		return internalPlugin, err
	}

	rawplugin, err := rpcClient.Dispense(sdk_config.DefaultPluginName)
	if err != nil {
		return internalPlugin, err
	}

	internalPlugin = &InternalPlugin{
		Plugin: rawplugin.(plugin.Plugin),
		Client: client,
	}

	return internalPlugin, nil
}

// ClientProvider is an interface for goplugin.Client returned by
// goplugin.NewClient.
type ClientProvider interface {
	Client() (goplugin.ClientProtocol, error)
	Kill()
}

type InternalPlugin struct {
	Plugin plugin.Plugin
	Client ClientProvider
}
