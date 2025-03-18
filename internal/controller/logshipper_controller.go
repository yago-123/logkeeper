/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	loggingv1alpha1 "github.com/yago-123/logkeeper/api/v1alpha1"
)

// LogShipperReconciler reconciles a LogShipper object
type LogShipperReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=logging.yago.ninja,resources=logshippers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.yago.ninja,resources=logshippers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=logging.yago.ninja,resources=logshippers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LogShipper object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.2/pkg/reconcile
func (r *LogShipperReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Retrieve the custom resource from the API server
	logShipper := &loggingv1alpha1.LogShipper{}
	if err := r.Get(ctx, req.NamespacedName, logShipper); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Obtain list of nodes and filter them based on nodeSelector
	var nodes corev1.NodeList
	if err := r.List(ctx, &nodes); err != nil {
		return ctrl.Result{}, err
	}
	desiredNodes := []corev1.Node{}
	for _, node := range nodes.Items {
		if matchesSelector(node, logShipper.Spec.NodeSelector) {
			desiredNodes = append(desiredNodes, node)
		}
	}

	// Ensure a pod exists for each desired node
	for _, node := range desiredNodes {
		// Generate pod name based on log shipper name and node name
		podName := fmt.Sprintf("%s-%s", logShipper.Name, node.Name)

		existingPod := &corev1.Pod{}
		err := r.Get(ctx, client.ObjectKey{Name: podName}, existingPod)
		if err != nil && errors.IsNotFound(err) {
			// Pod doesn't exist, create it
			newPod := createLogShipperPod(podName, logShipper.Spec.Image, node.Name, logShipper.Namespace)
			if err := r.Create(ctx, newPod); err != nil {
				return ctrl.Result{}, err
			}
		} else if err != nil {
			// Error occurred while checking for existing pod
			return ctrl.Result{}, err
		} else {
			// Pod exists, ensure it's running and matches desired state
			// todo(): ensure it's running and matches desired state
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LogShipperReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1alpha1.LogShipper{}).
		Named("logshipper").
		Complete(r)
}

// matchesSelector returns true if the node matches the given nodeSelector keys and values.
func matchesSelector(node corev1.Node, selector map[string]string) bool {
	if len(selector) == 0 {
		return true
	}

	for key, value := range selector {
		if nodeValue, exists := node.Labels[key]; !exists || nodeValue != value {
			return false
		}
	}
	return true
}

func createLogShipperPod(podName, podImage, nodeName, namespace string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    map[string]string{"app": "log-shipper", "node": nodeName},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "log-shipper",
					Image: podImage,
				},
			},
			NodeName: nodeName,
		},
	}
}
