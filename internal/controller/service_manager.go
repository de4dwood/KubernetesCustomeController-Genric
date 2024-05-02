package controller

import (
	"context"
	stackv1alpha1 "kachi-bits/api/v1alpha1"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *AppStackReconciler) ServiceManager(ctx context.Context, appstack *stackv1alpha1.AppStack, labels map[string]string) error {
	l := log.FromContext(ctx)
	serviceNamespacedName := types.NamespacedName{
		Namespace: appstack.Namespace,
		Name:      appstack.Name,
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
			log.Log.Error(err, "couldnt update appstack service status")
			return err
		}
	}

	return nil
}
