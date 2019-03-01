/*
Copyright 2019 Pivotal.

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

package rabbitmqcluster

import (
	"context"

	rabbitmqv1beta1 "github.com/pivotal/rabbitmq-for-kubernetes/pkg/apis/rabbitmq/v1beta1"
	"github.com/pivotal/rabbitmq-for-kubernetes/pkg/helpers"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new RabbitmqCluster Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRabbitmqCluster{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("rabbitmqcluster-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to RabbitmqCluster
	err = c.Watch(&source.Kind{Type: &rabbitmqv1beta1.RabbitmqCluster{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by RabbitmqCluster - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &rabbitmqv1beta1.RabbitmqCluster{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRabbitmqCluster{}

// ReconcileRabbitmqCluster reconciles a RabbitmqCluster object
type ReconcileRabbitmqCluster struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a RabbitmqCluster object and makes changes based on the state read
// and what is in the RabbitmqCluster.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rabbitmq.pivotal.io,resources=rabbitmqclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rabbitmq.pivotal.io,resources=rabbitmqclusters/status,verbs=get;update;patch
func (r *ReconcileRabbitmqCluster) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the RabbitmqCluster instance
	instance := &rabbitmqv1beta1.RabbitmqCluster{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// TODO(user): Change this to be the object type created by your controller
	// Define the desired Deployment object
	// deploy := &appsv1.Deployment{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      instance.Name + "-deployment",
	// 		Namespace: instance.Namespace,
	// 	},
	// 	Spec: appsv1.DeploymentSpec{
	// 		Selector: &metav1.LabelSelector{
	// 			MatchLabels: map[string]string{"deployment": instance.Name + "-deployment"},
	// 		},
	// 		Template: corev1.PodTemplateSpec{
	// 			ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"deployment": instance.Name + "-deployment"}},
	// 			Spec: corev1.PodSpec{
	// 				Containers: []corev1.Container{
	// 					{
	// 						Name:  "nginx",
	// 						Image: "nginx",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// fmt.Println("CHECKING WORKING DIRECTORY")
	// dir, err := os.Getwd()
	// if err != nil {
	// 	return reconcile.Result{}, err
	// }
	// fmt.Println(dir)
	resourcesYml, err := helpers.Build("./templates", instance.Name, instance.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	resources, decodeError := helpers.Decode(resourcesYml)
	if decodeError != nil {
		return reconcile.Result{}, decodeError
	}
	for _, resource := range resources {
		if err := controllerutil.SetControllerReference(instance, resource.ResourceObject.(v1.Object), r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// TODO(user): Change this for the object type created by your controller
		// Check if the Deployment already exists
		found := resource.EmptyResource
		err = r.Get(context.TODO(), types.NamespacedName{Name: resource.Name, Namespace: resource.Namespace}, found)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Creating "+resource.ResourceObject.GetObjectKind().GroupVersionKind().Kind, "namespace", resource.Namespace, "name", resource.Name)
			err = r.Create(context.TODO(), resource.ResourceObject)
			if err != nil {
				return reconcile.Result{}, err
			}
		} else if err != nil {
			return reconcile.Result{}, err
		}
	}

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	// if !reflect.DeepEqual(deploy.Spec, found.Spec) {
	// 	found.Spec = deploy.Spec
	// 	log.Info("Updating Deployment", "namespace", deploy.Namespace, "name", deploy.Name)
	// 	err = r.Update(context.TODO(), found)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// }
	return reconcile.Result{}, nil
}