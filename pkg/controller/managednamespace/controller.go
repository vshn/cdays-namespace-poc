//
// Copyright (c) 2019, VSHN AG, info@vshn.ch
// Licensed under "BSD 3-Clause". See LICENSE file.
//
//

package managednamespace

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/go-logr/logr"
	"github.com/vshn/cdays-namespace-poc/pkg/apis/control/v1alpha1"
	syncv1alpha1 "github.com/vshn/espejo/pkg/apis/sync/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var commonLabels = map[string]string{
	"app.kubernetes.io/managed-by": "appuio-namespace-operator",
}

func (r *ReconcileManagedNamespace) handle(ctx context.Context, managedNamespace *v1alpha1.ManagedNamespace, request reconcile.Request, reqLogger logr.Logger) error {

	namespace := newNamespaceForCR(managedNamespace)
	newNamespace := namespace.DeepCopy()

	res, err := controllerutil.CreateOrUpdate(ctx, r.client, namespace, controllerutil.MutateFn(func(existing apiruntime.Object) error {
		existingNamespace := existing.(*corev1.Namespace)
		existingNamespace.Annotations = newNamespace.Annotations
		existingNamespace.Labels = newNamespace.Labels
		existingNamespace.OwnerReferences = []metav1.OwnerReference{}
		if err := controllerutil.SetControllerReference(managedNamespace, existingNamespace, r.scheme); err != nil {
			return err
		}
		return nil
	}))

	if err != nil {
		return err
	}

	reqLogger.Info("Reconciled ManagedNamespace: "+string(res), "Namespace.Name", namespace.Name)

	managedNamespace.Status.CreatedNamespace = namespace.UID
	managedNamespace.Status.Phase = corev1.NamespaceActive

	if err := r.client.Status().Update(ctx, managedNamespace); err != nil {
		return err
	}

	syncConfig := newSyncConfigForCR(managedNamespace)
	newSyncConfig := syncConfig.DeepCopy()

	res, err = controllerutil.CreateOrUpdate(ctx, r.client, syncConfig, controllerutil.MutateFn(func(existing apiruntime.Object) error {
		existingConfig := existing.(*syncv1alpha1.SyncConfig)
		existingConfig.Labels = commonLabels
		existingConfig.Spec = newSyncConfig.Spec
		return nil
	}))

	if err != nil {
		return err
	}

	reqLogger.Info("Reconciled SyncConfig: "+string(res), "SyncConfig.Name", syncConfig.Name, "SyncConfig.Namespace", syncConfig.Namespace)

	return nil
}

// newNamespaceForCR returns a new Namespace as specified by the cr
func newNamespaceForCR(cr *v1alpha1.ManagedNamespace) *corev1.Namespace {
	annotations := map[string]string{
		"appuio.ch/customer":        cr.Namespace,
		"openshift.io/display-name": cr.Spec.Description,
	}
	labels := commonLabels
	labels["appuio.ch/customer"] = cr.Namespace
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Namespace + "-" + cr.Name,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.NamespaceSpec{},
	}
}

func newSyncConfigForCR(cr *v1alpha1.ManagedNamespace) *syncv1alpha1.SyncConfig {
	networkPolicy := &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: networkingv1.SchemeGroupVersion.String(),
			Kind:       "NetworkPolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "allow-from-same-namespace",
			Labels: commonLabels,
		},
		Spec: networkingv1.NetworkPolicySpec{
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{},
						},
					},
				},
			},
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
			},
		},
	}
	networkPolicyUnstructuredMap, _ := apiruntime.DefaultUnstructuredConverter.ToUnstructured(networkPolicy)
	networkPolicyUnstructured := unstructured.Unstructured{
		Object: networkPolicyUnstructuredMap,
	}
	syncItems := []unstructured.Unstructured{
		networkPolicyUnstructured,
	}
	syncConfig := &syncv1alpha1.SyncConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "base-config",
			Namespace: cr.Namespace,
			Labels:    commonLabels,
		},
		Spec: syncv1alpha1.SyncConfigSpec{
			NamespaceSelector: &syncv1alpha1.NamespaceSelector{
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"appuio.ch/customer": cr.Namespace,
					},
				},
			},
			SyncItems: syncItems,
		},
	}
	return syncConfig
}
