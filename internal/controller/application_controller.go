/*
Copyright 2023 ahwhya.

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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/ahwhy/clusterops-operator/api/v1"
)

const (
	GenericRequeueDuraiton = 1 * time.Minute
)

var (
	CounterReconcileApplication int64
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.clusterops.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.clusterops.io,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.clusterops.io,resources=applications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Reconcile 调谐过程是并发执行的，这里等待 100毫秒，并对 reconcile 次数进行累计
	<-time.NewTicker(100 * time.Millisecond).C
	logger := log.FromContext(ctx)

	CounterReconcileApplication += 1
	logger.Info("Starting a reconile", "number", CounterReconcileApplication)

	// Get Application
	// 实例化一个 *v1.Application 类型的 app 对象，通过 r.Get() 方法查询触发当前调谐逻辑对应的 Application，将其写入 app
	app := &v1.Application{}
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		// 当 Application 不存在，结束本轮调谐
		if errors.IsNotFound(err) {
			logger.Info("Application not found.")
			return ctrl.Result{}, nil
		}
		// 其他错误情况，通过重试来处理
		logger.Error(err, "Failed to get the Application, will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
	}

	// reconcile sub-resource
	var result ctrl.Result
	var err error

	result, err = r.reconcileDeployment(ctx, app)
	if err != nil {
		logger.Error(err, "Fail to reconcile Deployment.")
		return result, err
	}

	result, err = r.reconcileService(ctx, app)
	if err != nil {
		logger.Error(err, "")
		return result, err
	}

	logger.Info("All resources have been reconciled.")
	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) reconcileDeployment(ctx context.Context, app *v1.Application) (ctrl.Result, error)

func (r *ApplicationReconciler) reconcileService(ctx context.Context, app *v1.Application) (ctrl.Result, error)

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Application{}).
		Complete(r)
}
