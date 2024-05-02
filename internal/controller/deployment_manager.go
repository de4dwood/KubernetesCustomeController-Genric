package controller

import (
	"context"
	stackv1alpha1 "kachi-bits/api/v1alpha1"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *AppStackReconciler) DeploymentManager(ctx context.Context, appstack *stackv1alpha1.AppStack, labels map[string]string) error {
	l := log.FromContext(ctx)
	deploymentNamespacedName := types.NamespacedName{
		Namespace: appstack.Namespace,
		Name:      appstack.Name,
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
		*deployment.Spec.Replicas = *appstack.Spec.Workload.Replicas
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
			log.Log.Error(err, "couldnt update appstack deployment status")
			return err
		}
	}

	return nil
}
