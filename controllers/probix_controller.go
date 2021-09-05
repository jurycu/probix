/*
Copyright 2021 jurycu.

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
	"reflect"
	"time"

	v1prometheus "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ferulaxv1alpha1 "github.com/jurycu/probix/api/v1alpha1"
)

// ProbixReconciler reconciles a Probix object
type ProbixReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Finalizer
}

const (
	delay = 3 * time.Second
)

//+kubebuilder:rbac:groups=acmp.aliyun-inc.com,resources=probixes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=acmp.aliyun-inc.com,resources=probixes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=acmp.aliyun-inc.com,resources=probixes/finalizers,verbs=update

//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=podmonitors,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Probix object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ProbixReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	probixLog := log.FromContext(ctx)
	instance := ferulaxv1alpha1.Probix{}
	if err := r.Client.Get(context.TODO(), req.NamespacedName, &instance); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		probixLog.Error(err, "Get HTTPTrigger failed")
		return ctrl.Result{Requeue: true, RequeueAfter: delay}, err
	}

	//删除逻辑
	if instance.DeletionTimestamp != nil && instance.DeletionGracePeriodSeconds != nil {
		//delete logic
		probixLog.Info("ProbixReconciler start...|delete")
		if err := r.delete(ctx, instance); err != nil {
			return reconcile.Result{Requeue: true, RequeueAfter: delay}, err
		}
		return reconcile.Result{}, nil
	}

	pm := v1prometheus.PodMonitor{}
	if err := r.Client.Get(ctx, req.NamespacedName, &pm); err != nil && errors.IsNotFound(err) {
		//创建逻辑
		probixLog.Info("ProbixReconciler start...|create")
		if err := r.create(ctx, instance); err != nil {
			if err = r.updateStatus(req, ferulaxv1alpha1.FAILED, "create failed,err:"+err.Error()); err != nil {
				return reconcile.Result{Requeue: true, RequeueAfter: delay}, err
			}
			return reconcile.Result{Requeue: true, RequeueAfter: delay}, err
		}
		if err = r.updateStatus(req, ferulaxv1alpha1.SUCCESS, "create success!"); err != nil {
			return reconcile.Result{Requeue: true, RequeueAfter: delay}, err
		}
		return reconcile.Result{}, nil
	}

	//更新逻辑
	probixLog.Info("ProbixReconciler start...|update")
	if err := r.update(ctx, req, instance); err != nil {
		if err = r.updateStatus(req, ferulaxv1alpha1.SUCCESS, "update failed,err:"+err.Error()); err != nil {
			return reconcile.Result{Requeue: true, RequeueAfter: delay}, err
		}
		return reconcile.Result{Requeue: true, RequeueAfter: delay}, err
	}
	if err := r.updateStatus(req, ferulaxv1alpha1.SUCCESS, "update success!"); err != nil {
		return reconcile.Result{Requeue: true, RequeueAfter: delay}, err
	}
	return ctrl.Result{}, nil
}

func (r *ProbixReconciler) delete(ctx context.Context, instance ferulaxv1alpha1.Probix) error {
	probixLog := log.FromContext(ctx)
	deletePodMonitor := getDeleteInstance(instance)
	if err := r.Delete(ctx, deletePodMonitor); err != nil {
		probixLog.Error(err, "failed to delete PodMonitor")
		return err
	}
	probixLog.Info("success to delete PodMonitor")
	if err := r.Finalizer.RemoveFinalizer(r.Client, &instance); err != nil {
		return err
	}
	return nil
}

func (r *ProbixReconciler) create(ctx context.Context, instance ferulaxv1alpha1.Probix) error {
	probixLog := log.FromContext(ctx)
	createPodMonitor := getCreatePodMonitor(instance)
	if err := r.Create(ctx, createPodMonitor); err != nil {
		probixLog.Error(err, "failed to create PodMonitor")
		return err
	}
	probixLog.Info("success to create PodMonitor,name:" + createPodMonitor.Name)
	if err := r.Finalizer.AddFinalizer(r.Client, &instance); err != nil {
		return err
	}
	return nil
}

func (r *ProbixReconciler) update(ctx context.Context, req ctrl.Request, instance ferulaxv1alpha1.Probix) error {
	probixLog := log.FromContext(ctx)
	oldPM := v1prometheus.PodMonitor{}
	if err := r.Client.Get(ctx, req.NamespacedName, &oldPM); err != nil {
		return err
	}
	updatePodMonitor := getUpdatePodMonitor(instance, oldPM)
	if err := r.Update(ctx, updatePodMonitor); err != nil {
		probixLog.Error(err, "failed to update PodMonitor")
		return err
	}
	return nil
}

func (r *ProbixReconciler) updateStatus(req ctrl.Request, status ferulaxv1alpha1.Result, msg string) error {
	instanceNew := ferulaxv1alpha1.Probix{}
	if err := r.Client.Get(context.TODO(), req.NamespacedName, &instanceNew); err != nil {
		return err
	}
	instanceNew.Status.Status = status
	instanceNew.Status.Message = msg
	if err := r.Client.Status().Update(context.TODO(), &instanceNew); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProbixReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ferulaxv1alpha1.Probix{}, builder.WithPredicates(predicate.Funcs{
			CreateFunc: func(event event.CreateEvent) bool {
				return true
			},
			DeleteFunc: func(deleteEvent event.DeleteEvent) bool {
				return !deleteEvent.DeleteStateUnknown
			},
			UpdateFunc: func(updateEvent event.UpdateEvent) bool {
				if !reflect.DeepEqual(updateEvent.ObjectOld.GetLabels(), updateEvent.ObjectNew.GetLabels()) {
					return true
				}
				return updateEvent.ObjectOld.GetGeneration() != updateEvent.ObjectNew.GetGeneration()
			},
			GenericFunc: func(genericEvent event.GenericEvent) bool {
				return true
			},
		})).
		Complete(r)
}
