package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	jsonpatch "github.com/evanphx/json-patch"
	jsonpatch1 "github.com/mattbaird/jsonpatch"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	GlobalLabel        = "global"
	SelectorKey        = "clusterSelector"
	SelectedNodeKey    = "volume.kubernetes.io/selected-node"
	HostNameKey        = "kubernetes.io/hostname"
	BetaHostNameKey    = "beta.kubernetes.io/hostname"
	LabelOSBeta        = "beta.kubernetes.io/os"
	KosmosPodLabel     = "kosmos-pod"
	KosmosKubeletLabel = "kosmos-kubelet"
	TrippedLabels      = "tripped-labels"
	ClusterID          = "clusterID"
	NodeType           = "type"
	BatchPodLabel      = "pod-group.scheduling.sigs.k8s.io"
	TaintNodeNotReady  = "node.kubernetes.io/not-ready"

	TaintNodeUnreachable = "node.kubernetes.io/unreachable"
	CreatedbyDescheduler = "create-by-descheduler"
	DescheduleCount      = "sigs.k8s.io/deschedule-count"
)

type ClustersNodeSelection struct {
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Affinity     *corev1.Affinity    `json:"affinity,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
}

func CreateMergePatch(original, new interface{}) ([]byte, error) {
	pvByte, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}
	cloneByte, err := json.Marshal(new)
	if err != nil {
		return nil, err
	}
	patch, err := jsonpatch.CreateMergePatch(pvByte, cloneByte)
	if err != nil {
		return nil, err
	}
	return patch, nil
}

func CreateJSONPatch(original, new interface{}) ([]byte, error) {
	pvByte, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}
	cloneByte, err := json.Marshal(new)
	if err != nil {
		return nil, err
	}
	patchs, err := jsonpatch1.CreatePatch(pvByte, cloneByte)
	if err != nil {
		return nil, err
	}
	patchBytes, err := json.Marshal(patchs)
	if err != nil {
		return nil, err
	}
	return patchBytes, nil
}

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
var onlyOneSignalHandler = make(chan struct{})
var shutdownHandler chan os.Signal

func SetupSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice
	shutdownHandler = make(chan os.Signal, 2)
	stop := make(chan struct{})
	signal.Notify(shutdownHandler, shutdownSignals...)
	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1) // second signal. Exit directly.
	}()
	return stop
}

type Opts func(*rest.Config)

func NewClient(configPath string, opts ...Opts) (kubernetes.Interface, error) {
	var (
		config *rest.Config
		err    error
	)
	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("could not read config file for cluster: %v", err)
		}
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for master cluster: %v", err)
	}
	return client, nil
}

func NewClientFromByte(kubeConfig []byte, opts ...Opts) (kubernetes.Interface, error) {
	var (
		config *rest.Config
		err    error
	)

	clientconfig, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}
	config, err = clientconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for master cluster: %v", err)
	}
	return client, nil
}

func NewMetricClient(configPath string, opts ...Opts) (versioned.Interface, error) {
	var (
		config *rest.Config
		err    error
	)
	config, err = clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("could not read config file for cluster: %v", err)
		}
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	metricClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for master cluster: %v", err)
	}
	return metricClient, nil
}

func NewMetricClientFromByte(kubeConfig []byte, opts ...Opts) (versioned.Interface, error) {
	var (
		config *rest.Config
		err    error
	)

	clientconfig, err := clientcmd.NewClientConfigFromBytes(kubeConfig)
	if err != nil {
		return nil, err
	}
	config, err = clientconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(config)
	}

	metricClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create client for master cluster: %v", err)
	}
	return metricClient, nil
}

func IsVirtualNode(node *corev1.Node) bool {
	if node == nil {
		return false
	}
	valStr, exist := node.ObjectMeta.Labels[NodeType]
	if !exist {
		return false
	}
	return valStr == KosmosKubeletLabel
}

func IsVirtualPod(pod *corev1.Pod) bool {
	if pod.Labels != nil && pod.Labels[KosmosPodLabel] == "true" {
		return true
	}
	return false
}

func GetClusterID(node *corev1.Node) string {
	if node == nil {
		return ""
	}
	clusterName, exist := node.ObjectMeta.Labels[ClusterID]
	if !exist {
		return ""
	}
	return clusterName
}

func UpdateConfigMap(old, new *corev1.ConfigMap) {
	old.Labels = new.Labels
	old.Data = new.Data
	old.BinaryData = new.BinaryData
}

func UpdateSecret(old, new *corev1.Secret) {
	old.Labels = new.Labels
	old.Data = new.Data
	old.StringData = new.StringData
	old.Type = new.Type
}
