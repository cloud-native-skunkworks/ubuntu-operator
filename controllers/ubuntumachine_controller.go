/*
Copyright 2022 Alex.

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
	"reflect"
	"strings"

	"github.com/cloud-native-skunkworks/ubuntu-operator/api/v1alpha1"
	ubuntumachineryiov1alpha1 "github.com/cloud-native-skunkworks/ubuntu-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// UbuntuMachineConfigurationReconciler reconciles a UbuntuMachine object
type UbuntuMachineConfigurationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *UbuntuMachineConfigurationReconciler) buildDaemonset(instance *v1alpha1.UbuntuMachineConfiguration) (*appsv1.DaemonSet, error) {
	// Create our payload configuration -
	var moduleList []string
	for _, mod := range instance.Spec.DesiredModules {

		joined := fmt.Sprintf("%s=%s", mod.Name, mod.Flags)
		moduleList = append(moduleList, joined)
	}
	var aptList []string
	for _, mod := range instance.Spec.DesiredPackages.Apt {
		aptList = append(aptList, mod.Name)
	}
	var snapList []string
	for _, mod := range instance.Spec.DesiredPackages.Snap {
		joined := fmt.Sprintf("%s=%s", mod.Name, mod.Confinement)
		snapList = append(snapList, joined)
	}

	hostPathType := v1.HostPathDirectoryOrCreate

	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-daemonset",
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"daemonset": instance.Name + "-daemonset"},
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"daemonset": instance.Name + "-daemonset"},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "socket-path",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/opt/ubuntu-operator/",
									Type: &hostPathType,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Env: []corev1.EnvVar{
								{
									Name:  "MODULE_LIST",
									Value: strings.Join(moduleList, ","),
								},
								{
									Name:  "APT_LIST",
									Value: strings.Join(aptList, ","),
								},
								{
									Name:  "SNAP_LIST",
									Value: strings.Join(snapList, ","),
								},
							},
							ImagePullPolicy: corev1.PullAlways,
							Name:            "controller",
							Image:           "tibbar/ubuntu-machine-controller:latest",
							SecurityContext: &corev1.SecurityContext{
								Privileged: &[]bool{true}[0],
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{"CAP_NET_BIND_SERVICE"},
								},
							},
							Args: []string{"--socketPath", "/opt/ubuntu-operator/uo2.sock"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "socket-path",
									MountPath: "/opt/ubuntu-operator/",
								},
							},
						},
					},
				},
			},
		},
	}
	return daemonset, nil
}

func (r *UbuntuMachineConfigurationReconciler) checkFinalizers(instance *v1alpha1.UbuntuMachineConfiguration,
	daemonset *appsv1.DaemonSet, ctx context.Context) (ctrl.Result, error) {
	finalizerName := "ubuntu.machinery.io/finalizer"
	// Check to see if the Cluster has a finalizer
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !ctrlutil.ContainsFinalizer(instance, finalizerName) {
			ctrlutil.AddFinalizer(instance, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// Check to see if the Cluster is under deletion
		if controllerutil.ContainsFinalizer(instance, finalizerName) {

			found := &appsv1.DaemonSet{}
			err := r.Get(context.TODO(), types.NamespacedName{Name: daemonset.Name, Namespace: daemonset.Namespace}, found)
			if err == nil {
				if err = r.Delete(ctx, found); err != nil {
					return ctrl.Result{}, err
				}
			}
			// remove our finalizer from the Cluster and update it.
			controllerutil.RemoveFinalizer(instance, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}
	}
	if err := controllerutil.SetControllerReference(instance, daemonset, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

//+kubebuilder:rbac:groups=ubuntu.machinery.io.canonical.com,resources=ubuntumachineconfiguration,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ubuntu.machinery.io.canonical.com,resources=ubuntumachineconfiguration/status,verbs=get;update;patch;create;delete
//+kubebuilder:rbac:groups=ubuntu.machinery.io.canonical.com,resources=ubuntumachineconfiguration/finalizers,verbs=patch;create;update;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the UbuntuMachine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile

func (r *UbuntuMachineConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	instance := &v1alpha1.UbuntuMachineConfiguration{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {

			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err

	}

	daemonset, err := r.buildDaemonset(instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Finalizers ----------------------------------------------------------------------
	if res, err := r.checkFinalizers(instance, daemonset, ctx); err != nil {
		return res, err
	}
	// Check Daemonset ------------------------------------------------------------------------
	found := &appsv1.DaemonSet{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: daemonset.Name, Namespace: daemonset.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		err = r.Create(context.TODO(), daemonset)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}
	if !reflect.DeepEqual(daemonset.Spec, found.Spec) {
		found.Spec = daemonset.Spec
		err = r.Update(context.TODO(), found)
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UbuntuMachineConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ubuntumachineryiov1alpha1.UbuntuMachineConfiguration{}).Watches(&source.Kind{Type: &appsv1.DaemonSet{}},
		&handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &v1alpha1.UbuntuMachineConfiguration{},
		}).
		Complete(r)
}
