package controller

import (
	"context"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/ahwhy/clusterops-operator/api/v1"
)

func (r *ApplicationReconciler) reconcileService(ctx context.Context, app *v1.Application) (
	ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get Service
	var svc = &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: app.Namespace,
		Name:      app.Name,
	}, svc)

	if err == nil {
		logger.Info("The Service has already exist.")
		// 判断 svc.Status 和 app.Status.Network 是否相等
		if reflect.DeepEqual(svc.Status, app.Status.Network) {
			return ctrl.Result{}, nil
		}

		app.Status.Network = svc.Status
		// 若不相等，则触发更新
		if err := r.Status().Update(ctx, app); err != nil {
			logger.Error(err, "Failed to update Application service status.")
			return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
		}

		logger.Info("Ths Applications status has been updated.")
		return ctrl.Result{}, nil
	}
	// 非 NotFound 的场景
	if !errors.IsNotFound(err) {
		logger.Error(err, "Failed to get service, will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
	}

	// 若 NotFound，则触发 Create
	newSvc := &corev1.Service{}
	newSvc.SetName(app.Name)
	newSvc.SetNamespace(app.Namespace)
	newSvc.SetLabels(app.Labels)
	newSvc.Spec = app.Spec.Service.ServiceSpec
	newSvc.Spec.Selector = app.Labels

	// 将当前创建的 newSvc 设置为 Application 类型的 app 资源的子资源
	if err := ctrl.SetControllerReference(app, newSvc, r.Scheme); err != nil {
		logger.Error(err, "Failed to Set ControllerReference, will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
	}
	if err := r.Create(ctx, newSvc); err != nil {
		logger.Error(err, "Failed to Create Service, will requeue after a short time.")
		return ctrl.Result{RequeueAfter: GenericRequeueDuraiton}, err
	}

	logger.Info("The Service has been created")
	return ctrl.Result{}, nil
}
