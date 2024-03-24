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
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	stackv1alpha1 "kachi-bits/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	finalizer := "appstack.kachi-bits/finalizer"
	if appstack.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(appstack, finalizer) {
			controllerutil.AddFinalizer(appstack, finalizer)
			if err := r.Update(ctx, appstack); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(appstack, finalizer) {
			if err := r.finalizerRun(appstack); err != nil {
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(appstack, finalizer)
			if err := r.Update(ctx, appstack); err != nil {
				return ctrl.Result{}, err
			}
		}
		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	if err = r.DeploymentManager(ctx, appstack, labels); err != nil {
		l.Error(err, "deploymentmgr failed")
		return ctrl.Result{}, err
	}
	if err = r.ServiceManager(ctx, appstack, labels); err != nil {
		l.Error(err, "deploymentmgr failed")
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

func (r *AppStackReconciler) DeploymentManager(ctx context.Context, appstack *stackv1alpha1.AppStack, labels map[string]string) error {
	l := log.FromContext(ctx)
	deploymentNamespacedName := types.NamespacedName{
		Namespace: appstack.Namespace,
		Name:      appstack.Name + "-deployment",
	}
	deployment := &appsv1.Deployment{}
	deployment.ObjectMeta.Name = deploymentNamespacedName.Name
	deployment.ObjectMeta.Namespace = deploymentNamespacedName.Namespace
	container := corev1.Container{}
	deployExist := true

	err := r.Get(ctx, deploymentNamespacedName, deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			deployExist = false
		} else {
			l.Error(err, "couldnt get deployment object")
			return err
		}
	}

	//annotations
	if appstack.Spec.Workload.DeploymentAnnotations != nil {
		deployment.ObjectMeta.Annotations = appstack.Spec.Workload.DeploymentAnnotations
	}
	if appstack.Spec.Workload.PodAnnotations != nil {

		deployment.Spec.Template.Annotations = appstack.Spec.Workload.PodAnnotations
	}
	//labels
	deployment.ObjectMeta.Labels = labels
	deployment.Spec.Template.Labels = labels
	if deployment.Spec.Selector != nil {
		deployment.Spec.Selector.MatchLabels = labels
	} else {
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
	}

	if appstack.Spec.Workload.Replicas != nil {
		deployment.Spec.Replicas = appstack.Spec.Workload.Replicas
	}
	//envs
	if appstack.Spec.Workload.Env != nil {
		container.Env = appstack.Spec.Workload.Env
	}
	if appstack.Spec.Workload.EnvFrom != nil {
		container.EnvFrom = appstack.Spec.Workload.EnvFrom
	}
	//volumes
	if appstack.Spec.Workload.VolumeMounts != nil {
		container.VolumeMounts = append(container.VolumeMounts, appstack.Spec.Workload.VolumeMounts...)
	}
	if appstack.Spec.Workload.Volumes != nil {
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, appstack.Spec.Workload.Volumes...)
	}
	//image
	if appstack.Spec.Workload.Image != "" {
		container.Image = appstack.Spec.Workload.Image
	}
	container.Name = appstack.Name
	//container
	deployment.Spec.Template.Spec.Containers = []corev1.Container{container}

	if err := controllerutil.SetControllerReference(appstack, deployment, r.Scheme); err != nil {
		return err
	}
	switch deployExist {
	case true:
		if err = r.Update(context.TODO(), deployment); err != nil {
			log.Log.Error(err, "couldnt update deployment")
			return err
		}

	case false:
		if err = r.Create(context.TODO(), deployment); err != nil {
			log.Log.Error(err, "couldnt create deployment")
			return err
		}
	}

	if !reflect.DeepEqual(appstack.Status.DeploymentStatus, deployment.Status) {
		appstack.Status.DeploymentStatus = deployment.Status
		if err = r.Status().Update(context.TODO(), appstack); err != nil {
			log.Log.Error(err, "couldnt update status")
			return err
		}
	}

	return nil
}
func (r *AppStackReconciler) ServiceManager(ctx context.Context, appstack *stackv1alpha1.AppStack, labels map[string]string) error {
	l := log.FromContext(ctx)
	serviceNamespacedName := types.NamespacedName{
		Namespace: appstack.Namespace,
		Name:      appstack.Name + "-service",
	}
	service := &corev1.Service{}
	service.ObjectMeta.Name = serviceNamespacedName.Name
	service.ObjectMeta.Namespace = serviceNamespacedName.Namespace
	serviceExist := true

	err := r.Get(ctx, serviceNamespacedName, service)
	if err != nil {
		if errors.IsNotFound(err) {
			serviceExist = false
		} else {
			l.Error(err, "couldnt get service object")
			return err
		}
	}

	//annotations
	if appstack.Spec.Service.Annotations != nil {
		service.ObjectMeta.Annotations = appstack.Spec.Service.Annotations
	}
	//labels
	service.ObjectMeta.Labels = labels
	service.Spec.Selector = labels

	if appstack.Spec.Service.Ports != nil {
		service.Spec.Ports = appstack.Spec.Service.Ports
	}
	switch strings.ToLower(appstack.Spec.Service.Type) {
	case "loadbalancer":
		service.Spec.Type = corev1.ServiceTypeLoadBalancer
	case "nodeport":
		service.Spec.Type = corev1.ServiceTypeNodePort
	default:
		service.Spec.Type = corev1.ServiceTypeClusterIP

	}

	if err := controllerutil.SetControllerReference(appstack, service, r.Scheme); err != nil {
		return err
	}
	switch serviceExist {
	case true:
		if err = r.Update(context.TODO(), service); err != nil {
			log.Log.Error(err, "couldnt update service")
			return err
		}

	case false:
		if err = r.Create(context.TODO(), service); err != nil {
			log.Log.Error(err, "couldnt create service")
			return err
		}
	}

	if !reflect.DeepEqual(appstack.Status.ServiceStatus, service.Status) {
		appstack.Status.ServiceStatus = service.Status
		if err = r.Status().Update(context.TODO(), appstack); err != nil {
			log.Log.Error(err, "couldnt update status")
			return err
		}
	}

	return nil
}

func (r *AppStackReconciler) finalizerRun(appstack *stackv1alpha1.AppStack) error {
	//nothgin to do for now
	return nil
}
