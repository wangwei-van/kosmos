package k8sadapter

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

func getSecrets(pod *corev1.Pod) []string {
	secretNames := []string{}
	for _, v := range pod.Spec.Volumes {
		switch {
		case v.Secret != nil:
			if strings.HasPrefix(v.Name, "default-token") {
				continue
			}
			klog.Infof("pod %s depends on secret %s", pod.Name, v.Secret.SecretName)
			secretNames = append(secretNames, v.Secret.SecretName)

		case v.CephFS != nil:
			klog.Infof("pod %s depends on secret %s", pod.Name, v.CephFS.SecretRef.Name)
			secretNames = append(secretNames, v.CephFS.SecretRef.Name)
		case v.Cinder != nil:
			klog.Infof("pod %s depends on secret %s", pod.Name, v.Cinder.SecretRef.Name)
			secretNames = append(secretNames, v.Cinder.SecretRef.Name)
		case v.RBD != nil:
			klog.Infof("pod %s depends on secret %s", pod.Name, v.RBD.SecretRef.Name)
			secretNames = append(secretNames, v.RBD.SecretRef.Name)
		default:
			klog.Warning("Skip other type volumes")
		}
	}
	if pod.Spec.ImagePullSecrets != nil {
		for _, s := range pod.Spec.ImagePullSecrets {
			secretNames = append(secretNames, s.Name)
		}
	}
	klog.Infof("pod %s depends on secrets %s", pod.Name, secretNames)
	return secretNames
}

func getConfigmaps(pod *corev1.Pod) []string {
	cmNames := []string{}
	for _, v := range pod.Spec.Volumes {
		if v.ConfigMap == nil {
			continue
		}
		cmNames = append(cmNames, v.ConfigMap.Name)
	}
	klog.Infof("pod %s depends on configMap %s", pod.Name, cmNames)
	return cmNames
}

func getPVCs(pod *corev1.Pod) []string {
	cmNames := []string{}
	for _, v := range pod.Spec.Volumes {
		if v.PersistentVolumeClaim == nil {
			continue
		}
		cmNames = append(cmNames, v.PersistentVolumeClaim.ClaimName)
	}
	klog.Infof("pod %s depends on pvc %v", pod.Name, cmNames)
	return cmNames
}

// nolint:unused
func checkNodeStatusReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type != corev1.NodeReady {
			continue
		}
		if condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// nolint:unused
func compareNodeStatusReady(old, new *corev1.Node) (bool, bool) {
	return checkNodeStatusReady(old), checkNodeStatusReady(new)
}

// nolint:unused
func podStopped(pod *corev1.Pod) bool {
	return (pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed) && pod.Spec.
		RestartPolicy == corev1.RestartPolicyNever
}

// nolint:unused
func nodeCustomLabel(node *corev1.Node, label string) {
	nodelabel := strings.Split(label, ":")
	if len(nodelabel) == 2 {
		node.ObjectMeta.Labels[strings.TrimSpace(nodelabel[0])] = strings.TrimSpace(nodelabel[1])
	}
}
