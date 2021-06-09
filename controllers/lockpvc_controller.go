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
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	k8scorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	lockk8siov1beta1 "github.com/alicefr/lock-pvc/api/v1beta1"
)

const (
	finalizer        = "lock.io/finalizer"
	LockedAnnotation = "locked"
)

// LockPVCReconciler reconciles a LockPVC object
type LockPVCReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func containsString(s []string, search string) bool {
	for _, v := range s {
		if search == v {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

//+kubebuilder:rbac:groups=lock.io,resources=lockpvcs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;update;patch;
//+kubebuilder:rbac:groups=lock.io,resources=lockpvcs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lock.io,resources=lockpvcs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LockPVC object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *LockPVCReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	reqLogger := r.Log.WithValues("lockpvc", req.NamespacedName)

	lockPVC := &lockk8siov1beta1.LockPVC{}
	if err = r.Client.Get(context.TODO(), req.NamespacedName, lockPVC); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	reqLogger.Info("Reconcile LockPVC")
	// Fetch the PVC to lock
	pvc := &k8scorev1.PersistentVolumeClaim{}
	if err = r.Client.Get(context.TODO(), client.ObjectKey{
		Namespace: lockPVC.ObjectMeta.Namespace,
		Name:      lockPVC.Spec.PVC,
	}, pvc); err != nil {
		return ctrl.Result{}, err
	}

	// Add finalizer to the LockPVC if it doesn't have
	if lockPVC.ObjectMeta.DeletionTimestamp.IsZero() {
		if !containsString(lockPVC.ObjectMeta.Finalizers, finalizer) {
			lockPVC.ObjectMeta.Finalizers = append(lockPVC.ObjectMeta.Finalizers, finalizer)
		}
		if err = r.Client.Update(context.TODO(), lockPVC); err != nil {
			return ctrl.Result{}, err
		}
		reqLogger.Info("Add finalizer")
	}

	// Check if the LockPVC is under deletion
	if !lockPVC.ObjectMeta.DeletionTimestamp.IsZero() {
		if containsString(lockPVC.ObjectMeta.GetFinalizers(), finalizer) {
			// Remove the annotation on the PVC if the PVC is locked and it is owened by the lockPVC
			if v, ok := pvc.ObjectMeta.Annotations[LockedAnnotation]; ok && v == lockPVC.ObjectMeta.Name {
				reqLogger.WithValues("PVC", lockPVC.Spec.PVC).Info("Remove annotation on PVC")
				// Remove the annotation from the PVC
				delete(pvc.ObjectMeta.Annotations, LockedAnnotation)
				if err := r.Update(context.Background(), pvc); err != nil {
					return ctrl.Result{}, err
				}
			}
			// TODO Stop watching the resource

			// Remove the finalizer
			lockPVC.ObjectMeta.Finalizers = removeString(lockPVC.ObjectMeta.Finalizers, finalizer)
			if err := r.Update(context.Background(), lockPVC); err != nil {
				return ctrl.Result{}, err
			}
			reqLogger.Info("Remove finalizer")
		}
		return ctrl.Result{}, nil
	}

	var acquired bool
	var message string
	// Check if the PVC is already lock and set the annotation on the pvc
	if v, ok := pvc.ObjectMeta.Annotations[LockedAnnotation]; !ok {
		reqLogger.WithValues("PVC", lockPVC.Spec.PVC).Info("Add annotation on PVC")
		patch := []byte(fmt.Sprintf(`{"metadata":{"annotations":{"%s": "%s"}}}`, LockedAnnotation, req.NamespacedName.Name))
		if err = r.Client.Patch(context.Background(), &k8scorev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: lockPVC.ObjectMeta.Namespace,
				Name:      lockPVC.Spec.PVC,
			},
		}, client.RawPatch(types.StrategicMergePatchType, patch)); err != nil {
			return r.updateStatus(false, "Failed to acquire lock", lockPVC)
		}
		acquired = true
		message = "Lock succesfully acquired"
	} else {
		// If the pvc already contains the annotation, check if it is owned by the lockpvc
		if lockPVC.ObjectMeta.Name != v {
			lockPVC.Status.Acquired = false
			lockPVC.Status.Message = fmt.Sprintf("PVC %s already lock", lockPVC.Spec.PVC)
			reqLogger.WithValues("PVC", lockPVC.Spec.PVC).Error(fmt.Errorf("PVC already locked by %s", v), "Aquire lock")
		} else {
			acquired = true
			message = "Lock succesfully acquired"
			reqLogger.WithValues("PVC", lockPVC.Spec.PVC).Info("Lock already acquired")
		}
	}
	// TODO Check if watch for the resource has been already started, otherwise start
	return r.updateStatus(acquired, message, lockPVC)
}

// SetupWithManager sets up the controller with the Manager.
func (r *LockPVCReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lockk8siov1beta1.LockPVC{}).
		Complete(r)
}

func (r *LockPVCReconciler) updateStatus(acquired bool, message string, lockPVC *lockk8siov1beta1.LockPVC) (ctrl.Result, error) {
	lockPVC.Status.Acquired = acquired
	lockPVC.Status.Message = message
	err := r.Client.Status().Update(context.TODO(), lockPVC)
	return ctrl.Result{}, err
}
