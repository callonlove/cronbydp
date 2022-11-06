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

	logr "github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CronbydpReconciler reconciles a Cronbydp object
type DeploymentWatcher struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

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
func (r *DeploymentWatcher) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := r.Log.WithValues("Deploymentwatcher", req.NamespacedName)
	instance := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			log.Info("Cronbydp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		} // Error reading the object - requeue the request.
		log.Error(err, "Failed to get Cronbydp")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}
	// log.Info("命名空间:" + req.Namespace + "下的deployment:" + req.Name + "发生变更")
	for _, ns := range NowConfig.Spec.Namespaces {
		if ns == instance.Namespace {
			log.Info("判断命名空间")
			for annk, annv := range instance.Annotations {
				log.Info("判断注解名称")
				if annk == NowConfig.Spec.Annotationname {
					cronjoblist := &batchv1.CronJobList{}
					r.List(ctx, cronjoblist)
					for cronk, _ := range cronjoblist.Items {
						log.Info("查找相同注解的cronjob")
						for cronannk, cronannv := range cronjoblist.Items[cronk].Annotations {
							if cronannk == annk && cronannv == annv {
								log.Info("判断相同注解名称的cronjob的镜像")
								if cronjoblist.Items[cronk].Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image != instance.Spec.Template.Spec.Containers[0].Image {
									log.Info("镜像不一致，更换相同注解名称的cronjob的镜像")
									cronjob := cronjoblist.Items[cronk]
									cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image = instance.Spec.Template.Spec.Containers[0].Image
									r.Update(ctx, &cronjob)
								}
							}
						}

					}
				}

			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentWatcher) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
