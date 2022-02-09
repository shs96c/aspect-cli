/*
Copyright © 2022 Aspect Build Systems Inc

Not licensed for re-use.
*/

package system

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	goplugin "github.com/hashicorp/go-plugin"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	rootFlags "aspect.build/cli/pkg/aspect/root/flags"
	"aspect.build/cli/pkg/aspecterrors"
	"aspect.build/cli/pkg/ioutils"
	"aspect.build/cli/pkg/plugin/sdk/v1alpha2/plugin"
	plugin_mock "aspect.build/cli/pkg/plugin/sdk/v1alpha2/plugin/mock"
	system_mock "aspect.build/cli/pkg/plugin/system/mock"
)

func createInterceptorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "TestCommand",
	}

	// Required flags for interceptor hooks
	cmd.PersistentFlags().Bool(rootFlags.InteractiveFlagName, false, "")

	return cmd
}

func TestPluginSystemInterceptors(t *testing.T) {
	t.Run("executes hooks in reverse order of interceptors", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup
		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}
		ctx := context.Background()
		cmd := createInterceptorCommand()

		ps := NewPluginSystem().(*pluginSystem)
		plugin := plugin_mock.NewMockPlugin(ctrl)
		ps.addPlugin(&InternalPlugin{
			plugin: plugin,
			client: system_mock.NewMockClientProvider(ctrl),
		})

		// Expect the callbacks in reverse-order of execution
		gomock.InOrder(
			plugin.EXPECT().PostRunHook(gomock.Any(), gomock.Any()),
			plugin.EXPECT().PostTestHook(gomock.Any(), gomock.Any()),
			plugin.EXPECT().PostBuildHook(gomock.Any(), gomock.Any()),
		)

		// Hook interceptors
		buildInterceptor := ps.BuildHooksInterceptor(streams)
		testInterceptor := ps.TestHooksInterceptor(streams)
		runInterceptor := ps.RunHooksInterceptor(streams)

		err := buildInterceptor(ctx, cmd, []string{}, func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return testInterceptor(ctx, cmd, args, func(ctx context.Context, cmd *cobra.Command, args []string) error {
				return runInterceptor(ctx, cmd, args, func(ctx context.Context, cmd *cobra.Command, args []string) error {
					return nil
				})
			})
		})

		g.Expect(err).To(BeNil())
	})

	t.Run("executes plugin hooks in order plugins are added", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup
		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}
		ctx := context.Background()
		cmd := createInterceptorCommand()

		// Plugins to be invoked
		ps := NewPluginSystem().(*pluginSystem)
		plugin1 := plugin_mock.NewMockPlugin(ctrl)
		plugin2 := plugin_mock.NewMockPlugin(ctrl)
		ps.addPlugin(&InternalPlugin{
			plugin: plugin1,
			client: system_mock.NewMockClientProvider(ctrl),
		})
		ps.addPlugin(&InternalPlugin{
			plugin: plugin2,
			client: system_mock.NewMockClientProvider(ctrl),
		})

		// Expect the callbacks in reverse-order of execution, plugins in order added
		gomock.InOrder(
			plugin1.EXPECT().PostTestHook(gomock.Any(), gomock.Any()),
			plugin2.EXPECT().PostTestHook(gomock.Any(), gomock.Any()),
			plugin1.EXPECT().PostBuildHook(gomock.Any(), gomock.Any()),
			plugin2.EXPECT().PostBuildHook(gomock.Any(), gomock.Any()),
		)

		// Hook interceptors
		buildInterceptor := ps.BuildHooksInterceptor(streams)
		testInterceptor := ps.TestHooksInterceptor(streams)

		err := buildInterceptor(ctx, cmd, []string{}, func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return testInterceptor(ctx, cmd, args, func(ctx context.Context, cmd *cobra.Command, args []string) error {
				return nil
			})
		})

		g.Expect(err).To(BeNil())
	})

	t.Run("returns pass nested interceptor errors to parent", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup
		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}
		ctx := context.Background()
		cmd := createInterceptorCommand()

		// Plugin to be invoked
		ps := NewPluginSystem().(*pluginSystem)
		plugin := plugin_mock.NewMockPlugin(ctrl)
		ps.addPlugin(&InternalPlugin{
			plugin: plugin,
			client: system_mock.NewMockClientProvider(ctrl),
		})

		// Expect the callbacks in reverse-order of execution
		gomock.InOrder(
			plugin.EXPECT().PostRunHook(gomock.Any(), gomock.Any()),
			plugin.EXPECT().PostTestHook(gomock.Any(), gomock.Any()),
			plugin.EXPECT().PostBuildHook(gomock.Any(), gomock.Any()),
		)

		// Hook interceptors
		buildInterceptor := ps.BuildHooksInterceptor(streams)
		testInterceptor := ps.TestHooksInterceptor(streams)
		runInterceptor := ps.RunHooksInterceptor(streams)

		// Return error in nested interceptor
		err := buildInterceptor(ctx, cmd, []string{}, func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return testInterceptor(ctx, cmd, args, func(ctx context.Context, cmd *cobra.Command, args []string) error {
				return runInterceptor(ctx, cmd, args, func(ctx context.Context, cmd *cobra.Command, args []string) error {
					return fmt.Errorf("test error")
				})
			})
		})

		g.Expect(err).To(MatchError("test error"))
	})

	t.Run("parent interceptor errors override child errors", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup
		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}
		ctx := context.Background()
		cmd := createInterceptorCommand()

		// Plugin to be invoked
		ps := NewPluginSystem().(*pluginSystem)
		plugin := plugin_mock.NewMockPlugin(ctrl)
		ps.addPlugin(&InternalPlugin{
			plugin: plugin,
			client: system_mock.NewMockClientProvider(ctrl),
		})

		// Expect the callbacks in reverse-order of execution
		gomock.InOrder(
			plugin.EXPECT().PostRunHook(gomock.Any(), gomock.Any()),
			plugin.EXPECT().PostTestHook(gomock.Any(), gomock.Any()),
			plugin.EXPECT().PostBuildHook(gomock.Any(), gomock.Any()),
		)

		// Hook interceptors
		buildInterceptor := ps.BuildHooksInterceptor(streams)
		testInterceptor := ps.TestHooksInterceptor(streams)
		runInterceptor := ps.RunHooksInterceptor(streams)

		// Override error of nested interceptor
		err := buildInterceptor(ctx, cmd, []string{}, func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return testInterceptor(ctx, cmd, args, func(ctx context.Context, cmd *cobra.Command, args []string) error {
				runInterceptor(ctx, cmd, args, func(ctx context.Context, cmd *cobra.Command, args []string) error {
					return fmt.Errorf("error 1")
				})
				return fmt.Errorf("error 2")
			})
		})

		g.Expect(err).To(MatchError("error 2"))
	})

	t.Run("ExitCode is not modified on error from interceptor", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup
		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}
		ctx := context.Background()
		cmd := createInterceptorCommand()

		ps := NewPluginSystem().(*pluginSystem)

		// Hook interceptor returning an error
		runInterceptor := ps.RunHooksInterceptor(streams)
		err := runInterceptor(ctx, cmd, []string{}, func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return &aspecterrors.ExitError{
				Err:      fmt.Errorf("error 1"),
				ExitCode: 123,
			}
		})

		g.Expect(err).NotTo(BeNil())
		g.Expect(err.(*aspecterrors.ExitError).Err).To(MatchError("error 1"))
		g.Expect(err.(*aspecterrors.ExitError).ExitCode).To(Equal(123))
	})

	t.Run("ExitCode set to 1 on interceptor error of type aspecterrors.ExitError when a plugin returns an error", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Setup
		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}
		ctx := context.Background()
		cmd := createInterceptorCommand()

		ps := NewPluginSystem().(*pluginSystem)

		// Plugin returning an error
		plugin := plugin_mock.NewMockPlugin(ctrl)
		plugin.EXPECT().
			PostRunHook(gomock.Any(), gomock.Any()).
			DoAndReturn(func(
				isInteractiveMode bool,
				promptRunner ioutils.PromptRunner,
			) error {
				return fmt.Errorf("plugin error")
			})
		ps.addPlugin(&InternalPlugin{
			plugin: plugin,
			client: system_mock.NewMockClientProvider(ctrl),
		})

		// Hook interceptors
		runInterceptor := ps.RunHooksInterceptor(streams)
		err := runInterceptor(ctx, cmd, []string{}, func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return &aspecterrors.ExitError{
				Err:      fmt.Errorf("interceptor error"),
				ExitCode: 123,
			}
		})

		g.Expect(err).NotTo(BeNil())
		g.Expect(err.(*aspecterrors.ExitError).ExitCode).To(Equal(1))
	})
}

// Finder returning an error
type errorFinder struct{}

func (f *errorFinder) Find() (string, error) {
	return "", fmt.Errorf("Find Fail")
}

// PluginParser returning an error
type errorParser struct{}

func (f *errorParser) Parse(path string) ([]AspectPlugin, error) {
	return nil, fmt.Errorf("Parse Fail")
}

// PluginParser returning a preset array of plugins
type inlinePluginParser struct {
	plugins []AspectPlugin
}

// ClientFactory - pops plugins of the array when New() is invoked
type mockClientFactory struct {
	plugins []InternalPlugin
}

func (cf *mockClientFactory) New(config *goplugin.ClientConfig) (*InternalPlugin, error) {
	plugin, plugins := &cf.plugins[0], append(cf.plugins[:0], cf.plugins[0+1:]...)
	cf.plugins = plugins
	return plugin, nil
}

// Empty test plugin
type TestPlugin struct {
	plugin.Base
}

func (f *inlinePluginParser) Parse(path string) ([]AspectPlugin, error) {
	return f.plugins, nil
}

func TestConfigure(t *testing.T) {
	t.Run("fails when Finder fails to find plugin config file", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}

		ps := &pluginSystem{
			finder: &errorFinder{},
		}

		err := ps.Configure(streams)

		g.Expect(err).NotTo(BeNil())
	})

	t.Run("fails when Parser fails to parse plugin config file", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}

		finder := system_mock.NewMockFinder(ctrl)
		finder.EXPECT().Find().Return("/foo/bar")

		ps := &pluginSystem{
			finder: finder,
			parser: &errorParser{},
		}

		err := ps.Configure(streams)

		g.Expect(err).NotTo(BeNil())
	})

	t.Run("works when 0 plugins are found in config file", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}

		finder := system_mock.NewMockFinder(ctrl)
		finder.EXPECT().Find().Return("/foo/bar")

		ps := &pluginSystem{
			finder: finder,
			parser: &inlinePluginParser{plugins: []AspectPlugin{}},
		}

		err := ps.Configure(streams)

		g.Expect(err).To(BeNil())
	})

	t.Run("creates and persists each plugin", func(t *testing.T) {
		g := NewGomegaWithT(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var stdout strings.Builder
		streams := ioutils.Streams{Stdout: &stdout, Stderr: &stdout}

		testPlugin := AspectPlugin{
			Name:     "test plugin",
			From:     "...",
			LogLevel: "debug",
		}
		testPlugin2 := AspectPlugin{
			Name:     "test plugin2",
			From:     "...",
			LogLevel: "debug",
		}

		finder := system_mock.NewMockFinder(ctrl)
		finder.EXPECT().Find().Return("/foo/bar")

		p1 := &TestPlugin{}
		p2 := &TestPlugin{}

		ps := &pluginSystem{
			finder: finder,
			parser: &inlinePluginParser{plugins: []AspectPlugin{testPlugin, testPlugin2}},
			clientFactory: &mockClientFactory{
				plugins: []InternalPlugin{
					{
						plugin: p1,
						client: system_mock.NewMockClientProvider(ctrl),
					},
					{
						plugin: p2,
						client: system_mock.NewMockClientProvider(ctrl),
					},
				},
			},
			plugins: &PluginList{},
		}

		err := ps.Configure(streams)

		g.Expect(err).To(BeNil())
		g.Expect(ps.plugins.head.plugin.plugin).To(Equal(p1))
		g.Expect(ps.plugins.tail.plugin.plugin).To(Equal(p2))
	})
}
