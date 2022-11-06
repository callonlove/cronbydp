/*
Copyright 2022 wangwei.

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
	"time"

	wwbaopintuicnv1 "cronbydp/api/v1"

	logr "github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CronbydpReconciler reconciles a Cronbydp object
type CronbydpReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var NowConfig wwbaopintuicnv1.Cronbydp

//+kubebuilder:rbac:groups=ww.baopintui.cn,resources=cronbydps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ww.baopintui.cn,resources=cronbydps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ww.baopintui.cn,resources=cronbydps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cronbydp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *CronbydpReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)
	log := r.Log.WithValues("Cronbydpcontroller", req.NamespacedName)
	err := r.Get(ctx, req.NamespacedName, &NowConfig)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Info("Cronbydp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		} // Error reading the object - requeue the request.
		log.Error(err, "Failed to get Cronbydp")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	} else {
		namespaces := ""
		for _, ns := range NowConfig.Spec.Namespaces {
			if namespaces != "" {
				namespaces = namespaces + "," + ns
			} else {
				namespaces = namespaces + ns
			}

		}
		annotationname := NowConfig.Spec.Annotationname
		log.Info("设置监听的命名空间为:" + namespaces + ",监听注解名称为:" + annotationname)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CronbydpReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wwbaopintuicnv1.Cronbydp{}).
		Complete(r)
}
