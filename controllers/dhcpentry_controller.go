/*
Copyright 2021.

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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	// "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1beta1 "core.gardener.cloud/dhcpcontroller/api/v1beta1"
)

// DHCPEntryReconciler reconciles a DHCPEntry object
type DHCPEntryReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=core.core.gardener.cloud,resources=dhcpentries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.core.gardener.cloud,resources=dhcpentries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.core.gardener.cloud,resources=dhcpentries/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DHCPEntry object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *DHCPEntryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)

	// ctx := context.Background()
	log := r.Log.WithValues("DHCPEntry", req.NamespacedName)

	dhcpEntry := &corev1beta1.DHCPEntry{}
	if err := r.Get(ctx, req.NamespacedName, dhcpEntry); err != nil {

		log.Info("unable to fetch dhcpEntry", "Actual request", req, "Error", err)

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Fetched dhcpEntry")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DHCPEntryReconciler) SetupWithManager(mgr ctrl.Manager) error {

	predicateFunctions := predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// oldGeneration := e.MetaOld.GetGeneration()
			// newGeneration := e.MetaNew.GetGeneration()

			// vNetNew := e.ObjectNew.(*corev1beta1.DHCPEntry)
			// vNetOld := e.ObjectOld.(*corev1beta1.DHCPEntry)

			// if oldGeneration != newGeneration {
			// 	return true
			// }

			return false
		},

		CreateFunc: func(e event.CreateEvent) bool {
			// dhcpEntry := e.Object.(*corev1beta1.DHCPEntry)

			return true
		},

		DeleteFunc: func(e event.DeleteEvent) bool {
			// The reconciler adds a finalizer so we perform clean-up
			// when the delete timestamp is added
			// Suppress Delete events to avoid filtering them out in the Reconcile function
			// vNet := e.Object.(*corev1beta1.VirtualNet)
			// return vNet.Spec.NodeName != nil && *vNet.Spec.NodeName == r.HostName
			return true
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1beta1.DHCPEntry{}).
		WithEventFilter(predicateFunctions).
		Complete(r)
}
