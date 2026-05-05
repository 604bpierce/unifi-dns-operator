/*
Copyright 2026.

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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	dnsv1alpha1 "github.com/bpierce/unifi-dns-operator/api/v1alpha1"
	"github.com/bpierce/unifi-dns-operator/internal/unifi"
)

const (
	finalizerName = "dns.bx.network/finalizer"
)

// UnifiDNSEntryReconciler reconciles a UnifiDNSEntry object
type UnifiDNSEntryReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	UnifiClient *unifi.Client
}

// +kubebuilder:rbac:groups=dns.bx.network,resources=unifidnsentries,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dns.bx.network,resources=unifidnsentries/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dns.bx.network,resources=unifidnsentries/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *UnifiDNSEntryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch the UnifiDNSEntry instance
	entry := &dnsv1alpha1.UnifiDNSEntry{}
	if err := r.Get(ctx, req.NamespacedName, entry); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get UnifiDNSEntry")
		return ctrl.Result{}, err
	}

	// Handle deletion
	if !entry.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, entry)
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(entry, finalizerName) {
		controllerutil.AddFinalizer(entry, finalizerName)
		if err := r.Update(ctx, entry); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Convert CRD to Unifi DNS Policy
	policy := r.convertToUnifiPolicy(entry)

	var policyID string
	var err error

	// Check if policy already exists in Unifi
	if entry.Status.PolicyID != "" {
		// Update existing policy
		err = r.UnifiClient.UpdateDNSPolicy(entry.Status.PolicyID, policy)
		if err != nil {
			log.Error(err, "Failed to update DNS policy in Unifi")
			r.updateCondition(entry, "Degraded", metav1.ConditionTrue, "UpdateFailed", err.Error())
			if statusErr := r.Status().Update(ctx, entry); statusErr != nil {
				log.Error(statusErr, "Failed to update status")
			}
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}
		policyID = entry.Status.PolicyID
	} else {
		// Create new policy
		policyID, err = r.UnifiClient.CreateDNSPolicy(policy)
		if err != nil {
			log.Error(err, "Failed to create DNS policy in Unifi")
			r.updateCondition(entry, "Degraded", metav1.ConditionTrue, "CreateFailed", err.Error())
			if statusErr := r.Status().Update(ctx, entry); statusErr != nil {
				log.Error(statusErr, "Failed to update status")
			}
			return ctrl.Result{RequeueAfter: 30 * time.Second}, err
		}
	}

	// Update status
	entry.Status.PolicyID = policyID
	entry.Status.Synced = true
	now := metav1.Now()
	entry.Status.LastSyncTime = &now
	r.updateCondition(entry, "Available", metav1.ConditionTrue, "Synced", "DNS policy synchronized with Unifi")
	r.updateCondition(entry, "Degraded", metav1.ConditionFalse, "Synced", "DNS policy synchronized with Unifi")

	if err := r.Status().Update(ctx, entry); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	log.Info("Successfully reconciled UnifiDNSEntry", "policyID", policyID)
	return ctrl.Result{}, nil
}

// handleDeletion handles the deletion of a UnifiDNSEntry
func (r *UnifiDNSEntryReconciler) handleDeletion(ctx context.Context, entry *dnsv1alpha1.UnifiDNSEntry) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if controllerutil.ContainsFinalizer(entry, finalizerName) {
		// Delete policy from Unifi if it exists
		if entry.Status.PolicyID != "" {
			if err := r.UnifiClient.DeleteDNSPolicy(entry.Status.PolicyID); err != nil {
				log.Error(err, "Failed to delete DNS policy from Unifi")
				return ctrl.Result{RequeueAfter: 30 * time.Second}, err
			}
			log.Info("Deleted DNS policy from Unifi", "policyID", entry.Status.PolicyID)
		}

		// Remove finalizer
		controllerutil.RemoveFinalizer(entry, finalizerName)
		if err := r.Update(ctx, entry); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// convertToUnifiPolicy converts a UnifiDNSEntry CRD to a Unifi DNS Policy
func (r *UnifiDNSEntryReconciler) convertToUnifiPolicy(entry *dnsv1alpha1.UnifiDNSEntry) unifi.DNSPolicy {
	policy := unifi.DNSPolicy{
		Type:       unifi.DNSPolicyType(entry.Spec.Type),
		Enabled:    entry.Spec.Enabled,
		Domain:     entry.Spec.Domain,
		TTLSeconds: entry.Spec.TTLSeconds,
	}

	// Set type-specific fields
	data := entry.Spec.RecordData
	policy.IPv4Address = data.IPv4Address
	policy.IPv6Address = data.IPv6Address
	policy.TargetDomain = data.TargetDomain
	policy.MailServerDomain = data.MailServerDomain
	policy.Priority = data.Priority
	policy.TextValue = data.TextValue
	policy.Target = data.Target
	policy.Port = data.Port
	policy.Weight = data.Weight

	return policy
}

// updateCondition updates or adds a condition to the UnifiDNSEntry status
func (r *UnifiDNSEntryReconciler) updateCondition(entry *dnsv1alpha1.UnifiDNSEntry, conditionType string, status metav1.ConditionStatus, reason, message string) {
	condition := metav1.Condition{
		Type:               conditionType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: metav1.Now(),
		ObservedGeneration: entry.Generation,
	}

	meta.SetStatusCondition(&entry.Status.Conditions, condition)
}

// SetupWithManager sets up the controller with the Manager.
func (r *UnifiDNSEntryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dnsv1alpha1.UnifiDNSEntry{}).
		Named("unifidnsentry").
		Complete(r)
}
