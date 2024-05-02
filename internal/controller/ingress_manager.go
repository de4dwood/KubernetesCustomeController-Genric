package controller

import (
	"context"
	stackv1alpha1 "kachi-bits/api/v1alpha1"
	"reflect"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *AppStackReconciler) IngressManager(ctx context.Context, appstack *stackv1alpha1.AppStack, labels map[string]string) error {
	l := log.FromContext(ctx)
	IngressNamespacedName := types.NamespacedName{
		Namespace: appstack.Namespace,
		Name:      appstack.Name,
	}
	ingress := &networkingv1.Ingress{}
	ingress.ObjectMeta.Name = IngressNamespacedName.Name
	ingress.ObjectMeta.Namespace = IngressNamespacedName.Namespace
	IngressExist := true

	err := r.Get(ctx, IngressNamespacedName, ingress)
	if err != nil {
		if errors.IsNotFound(err) {
			IngressExist = false
		} else {
			l.Error(err, "couldnt get ingress object")
			return err
		}
	}
	//annotations
	if appstack.Spec.Ingress.Annotations != nil {
		ingress.ObjectMeta.Annotations = appstack.Spec.Ingress.Annotations
	}
	//rules
	ingress.ObjectMeta.Labels = labels
	rules := []networkingv1.IngressRule{}
	for _, i := range appstack.Spec.Ingress.Rules {

		rule := networkingv1.IngressRule{
			Host: i.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{},
			},
		}

		for _, j := range i.Paths {

			pathtype := networkingv1.PathTypeImplementationSpecific
			if j.PathType != nil {
				pathtype = *j.PathType
			}

			path := networkingv1.HTTPIngressPath{
				Path:     j.Path,
				PathType: &pathtype,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: appstack.Name,
						Port: networkingv1.ServiceBackendPort{
							Number: j.Port,
						},
					},
				},
			}
			rule.HTTP.Paths = append(rule.HTTP.Paths, path)

		}

		rules = append(rules, rule)
	}

	ingress.Spec.Rules = rules

	//TLS
	if len(appstack.Spec.Ingress.TLS) != 0 {
		ingress.Spec.TLS = appstack.Spec.Ingress.TLS
	}

	//ingressClass
	if appstack.Spec.Ingress.IngressCLass != nil {
		ingress.Spec.IngressClassName = appstack.Spec.Ingress.IngressCLass
	}

	if err := controllerutil.SetControllerReference(appstack, ingress, r.Scheme); err != nil {
		return err
	}
	switch IngressExist {
	case true:
		if err = r.Update(context.TODO(), ingress); err != nil {
			log.Log.Error(err, "couldnt update ingress")
			return err
		}

	case false:
		if err = r.Create(context.TODO(), ingress); err != nil {
			log.Log.Error(err, "couldnt create ingress")
			return err
		}
	}

	if !reflect.DeepEqual(appstack.Status.IngressStatus, ingress.Status) {
		appstack.Status.IngressStatus = ingress.Status
		if err = r.Status().Update(context.TODO(), appstack); err != nil {
			log.Log.Error(err, "couldnt update appstack ingress status")
			return err
		}
	}

	return nil
}
func (r *AppStackReconciler) finalizerRun(appstack *stackv1alpha1.AppStack) error {
	//nothgin to do for now
	return nil
}
