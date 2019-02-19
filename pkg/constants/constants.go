package constants

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
)

const (
	DefaultConfigName     = "merge-config"
	DefaultFilePermission = 0664
)

var (
	HomeKubeconfig       = clientcmd.RecommendedHomeFile
	HomeBackupKubeconfig = fmt.Sprintf("%s.bk", HomeKubeconfig)
	HomeKubeconfigPath   = clientcmd.RecommendedConfigDir
)
