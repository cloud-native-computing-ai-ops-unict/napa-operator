package napa

import (
	"fmt"
	napav1alpha1 "github.com/cloud-native-computing-ai-ops-unict/napa-operator/api/napa/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

// imageForNapaAgent gets the Agent image which is managed by this controller
// from the NAPA_AGENT_IMAGE environment variable defined in the config/manager/manager.yaml
func imageForNapaAgent() (string, error) {
	var imageEnvVar = "NAPA_AGENT_IMAGE"
	image, found := os.LookupEnv(imageEnvVar)
	if !found {
		return "", fmt.Errorf("Unable to find %s environment variable with the image", imageEnvVar)
	}
	return image, nil
}

func labelsForNapaAgent(name string) map[string]string {
	var imageTag string
	image, err := imageForNapaAgent()
	if err == nil {
		imageTag = strings.Split(image, ":")[1]
	}
	return map[string]string{
		"app.kubernetes.io/name":       "NapaAgent",
		"app.kubernetes.io/instance":   name,
		"app.kubernetes.io/version":    imageTag,
		"app.kubernetes.io/part-of":    "napa-operator",
		"app.kubernetes.io/created-by": "controller-manager",
	}
}

// deploymentForNapaAgent returns a napa Agent Deployment object
func (r *AgentReconciler) deploymentForNapaAgent(
	m *napav1alpha1.Agent) (*appsv1.Deployment, error) {
	ls := labelsForNapaAgent(m.Name)
	replicas := int32(1)
	image, err := imageForNapaAgent()
	if err != nil {
		return nil, err
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           image,
						Name:            "napa-agent",
						ImagePullPolicy: corev1.PullIfNotPresent,
						SecurityContext: &corev1.SecurityContext{
							RunAsNonRoot:             &[]bool{true}[0],
							RunAsUser:                &[]int64{1001}[0],
							AllowPrivilegeEscalation: &[]bool{false}[0],
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "napa-agent",
						}},
					}},
				},
			},
		},
	}
	// Set the ownerRef for the Deployment
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
	if err := ctrl.SetControllerReference(m, dep, r.Scheme); err != nil {
		return nil, err
	}
	return dep, nil
}
