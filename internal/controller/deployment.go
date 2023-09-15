package controller

import (
	"context"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/ahwhy/clusterops-operator/api/v1"
)

func (r *ApplicationReconciler) reconcileDeployment(ctx context.Context, app *v1.Application) (
	ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get Deployment
	var dp = &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: app.Namespace,
		Name:      app.Name,
	}, dp)

	if err == nil {
		logger.Info("The Deplyment has already exist.")
		// 判断 dp.Status 和 app.Status.Workflow 是否相等
		if reflect.DeepEqual(dp.Status, app.Status.Workflow) {
			return ctrl.Result{}, nil
		}

		app.Status.Workflow = dp.Status
		// 若不相等，则触发更新
		if err := r.Status().Update(ctx, app); err != nil {
			logger.Error(err, "Failed to update Application status.")
			return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
		}

		logger.Info("Ths Applications status has been updated.")
		return ctrl.Result{}, nil
	}
	// 非 NotFound 的场景
	if !errors.IsNotFound(err) {
		logger.Error(err, "Failed to get deployment, will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
	}

	// 若 NotFound，则触发 Create
	newDp := &appsv1.Deployment{}
	newDp.SetName(app.Name)
	newDp.SetNamespace(app.Namespace)
	newDp.SetLabels(app.Labels)
	newDp.Spec = app.Spec.Deployment.DeploymentSpec
	newDp.Spec.Template.SetLabels(app.Labels)

	// 将当前创建的 newDp 设置为 Application 类型的 app 资源的子资源
	if err := ctrl.SetControllerReference(app, newDp, r.Scheme); err != nil {
		logger.Error(err, "Failed to Set ControllerReference, will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
	}
	if err := r.Create(ctx, newDp); err != nil {
		logger.Error(err, "Failed to Create Deployment, will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
	}

	logger.Info("The Deployment has been created")
	return ctrl.Result{}, nil
}
