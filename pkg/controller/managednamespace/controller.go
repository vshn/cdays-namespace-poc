//
// Copyright (c) 2019, VSHN AG, info@vshn.ch
// Licensed under "BSD 3-Clause". See LICENSE file.
//
//

package managednamespace

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/vshn/cdays-namespace-poc/pkg/apis/control/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileManagedNamespace) handle(ctx context.Context, managedNamespace *v1alpha1.ManagedNamespace, request reconcile.Request, reqLogger logr.Logger) error {

	// Define a new Namespace object
	namespace := newNamespaceForCR(managedNamespace)

	// Set ManagedNamespace instance as the owner and controller
	if err := controllerutil.SetControllerReference(managedNamespace, namespace, r.scheme); err != nil {
		return err
	}

	res, err := controllerutil.CreateOrUpdate(ctx, r.client, namespace, controllerutil.MutateFn(func(existing apiruntime.Object) error {
		ns := existing.(*corev1.Namespace)
		ns.Annotations["openshift.io/display-name"] = managedNamespace.Spec.Description
		return nil
	}))

	if err != nil {
		return err
	}

	reqLogger.Info("Reconciled: "+string(res), "Namespace.Name", namespace.Name)

	managedNamespace.Status.CreatedNamespace = namespace.UID
	managedNamespace.Status.Phase = corev1.NamespaceActive
	if err := r.client.Status().Update(ctx, managedNamespace); err != nil {
		return err
	}

	return nil
}

// newNamespaceForCR returns a new Namespace as specified by the cr
func newNamespaceForCR(cr *v1alpha1.ManagedNamespace) *corev1.Namespace {
	labels := map[string]string{
		"app.kubernetes.io/managed-by": "appuio-namespace-operator",
	}
	annotations := map[string]string{
		"appuio.ch/customer":        cr.Namespace,
		"openshift.io/display-name": cr.Spec.Description,
	}
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Namespace + "-" + cr.Name,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.NamespaceSpec{},
	}
}
