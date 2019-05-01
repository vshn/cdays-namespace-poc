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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

	// Check if this Namespace already exists
	found := &corev1.Namespace{}
	err := r.client.Get(ctx, types.NamespacedName{Name: namespace.Name}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Namespace", "Namespace.Name", namespace.Name)
		err = r.client.Create(ctx, namespace)
		if err != nil {
			return err
		}

		// Namespace created successfully - don't requeue
		return nil
	} else if err != nil {
		return err
	}

	// Namespace already exists - don't requeue
	reqLogger.Info("Skip reconcile: Namespace already exists", "Namespace.Name", found.Name)
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
