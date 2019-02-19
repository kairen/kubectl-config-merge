package cmd

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/kairen/kubectl-config-merge/pkg/constants"
	"github.com/kairen/kubectl-config-merge/pkg/utils"
	"github.com/kairen/kubectl-config-merge/pkg/version"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	clientcmdapilatest "k8s.io/client-go/tools/clientcmd/api/latest"
)

var (
	mergeExample = `
	# Merge your kubeconfig and save it as "merge-config" in the current directory
	%[1]s config-merge kubeconfig1 kubeconfig2 ... 

	# Merge your kubeconfig and save it as "config" in "test" directory
	%[1]s config-merge kubeconfig1 kubeconfig2 ... --path test/config

	# Merge your kubeconfig with $HOME/.kube/config as "$HOME/.kube/config"
	%[1]s config-merge kubeconfig1 ... --home 

	# To view merged kubeconfig result
	%[1]s config-merge kubeconfig1 kubeconfig2 ... --view
`
)

const (
	configYAML = "yaml"
	configJSON = "json"
)

type MergeOptions struct {
	path      string
	output    string
	overwrite bool
	home      bool
	backup    bool
	view      bool
	version   bool

	kubeconfigs []string
	genericclioptions.IOStreams
}

func NewMergeOptions(streams genericclioptions.IOStreams) *MergeOptions {
	return &MergeOptions{IOStreams: streams}
}

func NewCmdMerge(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewMergeOptions(streams)
	cmd := &cobra.Command{
		Use:          "config-merge [kubeconfig] [flags]",
		Short:        "Merge two or more kubeconfig files",
		Example:      fmt.Sprintf(mergeExample, "kubectl"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if o.version {
				fmt.Fprintf(o.IOStreams.Out, "%s\n", version.GetVersion())
				return nil
			}

			if err := o.Parse(c, args); err != nil {
				return err
			}

			if err := o.Backup(); err != nil {
				return err
			}

			if err := o.Merge(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&o.output, "output", "o", configYAML, "Merged kubeconfig output type(json or yaml)")
	cmd.Flags().StringVar(&o.path, "path", constants.DefaultConfigName, "Merged kubeconfig output name and path")
	cmd.Flags().BoolVar(&o.home, "home", false, "If true, merge with $HOME/.kube/config to $HOME/.kube/config")
	cmd.Flags().BoolVar(&o.overwrite, "overwrite", false, "If true, force merge kubeconfig")
	cmd.Flags().BoolVar(&o.backup, "backup", true, "If true, backup $HOME/.kube/config file to $HOME/.kube/config.bk")
	cmd.Flags().BoolVar(&o.view, "view", false, "View merged kubeconfig result")
	cmd.Flags().BoolVar(&o.version, "version", o.version, "Show config-merge version")
	return cmd
}

func (o *MergeOptions) Parse(cmd *cobra.Command, args []string) error {
	o.kubeconfigs = args

	if len(o.kubeconfigs) == 0 {
		return cmd.Usage()
	}

	if len(o.kubeconfigs) < 1 && o.home {
		return fmt.Errorf("either one or more arguments are allowed")
	}

	if len(o.kubeconfigs) < 2 && !o.home {
		return fmt.Errorf("either two or more arguments are allowed")
	}

	if o.home {
		o.kubeconfigs = append(o.kubeconfigs, constants.HomeKubeconfig)
		o.path = constants.HomeKubeconfig
	}
	return nil
}

func (o *MergeOptions) encodeConfig(config *clientcmdapi.Config) ([]byte, error) {
	var err error
	var output []byte

	encode, err := runtime.Encode(clientcmdapilatest.Codec, config)
	if err != nil {
		return nil, err
	}

	switch o.output {
	case configYAML:
		output, err = yaml.JSONToYAML(encode)
	case configJSON:
		output, err = yaml.YAMLToJSON(encode)
		output, err = util.PrettyJson(output)
	default:
		err = fmt.Errorf("unsupported output type only save as yaml or json")
	}
	return output, err
}

func (o *MergeOptions) Backup() error {
	if !o.backup || !o.home || o.view {
		return nil
	}

	config, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return err
	}

	output, err := o.encodeConfig(config)
	if err != nil {
		return err
	}

	path := constants.HomeBackupKubeconfig
	if err := util.WriteFile(path, output, constants.DefaultFilePermission); err != nil {
		return err
	}
	return nil
}

func (o *MergeOptions) Merge() error {
	rules := clientcmd.ClientConfigLoadingRules{
		Precedence: o.kubeconfigs,
	}

	mergedConfig, err := rules.Load()
	if err != nil {
		return err
	}

	output, err := o.encodeConfig(mergedConfig)
	if err != nil {
		return err
	}

	if o.view {
		fmt.Fprintf(o.IOStreams.Out, "%s", string(output))
		return nil
	}

	if !o.overwrite {
		validations := map[string]int{
			"Contexts":  len(mergedConfig.Contexts),
			"Clusters":  len(mergedConfig.Clusters),
			"AuthInfos": len(mergedConfig.AuthInfos),
		}
		for k, v := range validations {
			if len(o.kubeconfigs) != v {
				return fmt.Errorf("merged config has conflict in %s. Use --overwrite to force merge", k)
			}
		}
	}

	if err := util.WriteFile(o.path, output, constants.DefaultFilePermission); err != nil {
		return err
	}
	return nil
}
