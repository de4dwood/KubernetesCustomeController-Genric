/*
Copyright 2024.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stackv1alpha1 "kachi-bits/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// AppStackReconciler reconciles a AppStack object
type AppStackReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stack.kachi-bits,resources=appstacks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stack.kachi-bits,resources=appstacks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stack.kachi-bits,resources=appstacks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AppStack object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *AppStackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	appstack := &stackv1alpha1.AppStack{}
	labels := map[string]string{
		"app":       appstack.Name,
		"manage-by": "app-stack",
	}
	err := r.Get(ctx, req.NamespacedName, appstack)
	if err != nil {
		l.Error(err, "couldnt get appstack object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// finalizer := "appstack.kachi-bits/finalizer"
	// if appstack.ObjectMeta.DeletionTimestamp.IsZero() {
	// 	if !controllerutil.ContainsFinalizer(appstack, finalizer) {
	// 		controllerutil.AddFinalizer(appstack, finalizer)
	// 		if err := r.Update(ctx, appstack); err != nil {
	// 			return ctrl.Result{}, err
	// 		}
	// 	}
	// } else {
	// 	if controllerutil.ContainsFinalizer(appstack, finalizer) {
	// 		if err := r.finalizerRun(appstack); err != nil {
	// 			return ctrl.Result{}, err
	// 		}

	// 		// remove our finalizer from the list and update it.
	// 		controllerutil.RemoveFinalizer(appstack, finalizer)
	// 		if err := r.Update(ctx, appstack); err != nil {
	// 			return ctrl.Result{}, err
	// 		}
	// 	}
	// 	// Stop reconciliation as the item is being deleted
	// 	return ctrl.Result{}, nil
	// }

	if err = r.DeploymentManager(ctx, appstack, labels); err != nil {
		l.Error(err, "deploymentManager failed")
		return ctrl.Result{}, err
	}
	if err = r.ServiceManager(ctx, appstack, labels); err != nil {
		l.Error(err, "serviceManager failed")
		return ctrl.Result{}, err
	}
	if err = r.IngressManager(ctx, appstack, labels); err != nil {
		l.Error(err, "ingressManager failed")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppStackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stackv1alpha1.AppStack{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		Complete(r)
}
