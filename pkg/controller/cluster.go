// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"context"
	"fmt"

	"github.com/atomix/kubernetes-controller/pkg/apis/cloud/v1beta2"
	"github.com/atomix/local-replica/pkg/apis/storage/v1beta1"
	"github.com/golang/protobuf/jsonpb"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	apiPort = 5678
)

func (r *Reconciler) addService(cluster *v1beta2.Cluster, storage *v1beta1.CacheStorage) error {
	log.Info("Creating service", "Name:", cluster.Name, "Namespace:", cluster.Namespace)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      cluster.Name,
			Labels:    cluster.Labels,
		},

		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "api",
					Port: apiPort,
				},
			},
			PublishNotReadyAddresses: true,
			ClusterIP:                "None",
		},
	}
	if err := controllerutil.SetControllerReference(cluster, service, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.TODO(), service)
}

func (r *Reconciler) addDeployment(cluster *v1beta2.Cluster, storage *v1beta1.CacheStorage) error {
	log.Info("Creating Deployment", "Name", cluster.Name, "Namespace", cluster.Namespace)
	var replicas int32
	replicas = 1
	var env []corev1.EnvVar
	env = append(env, corev1.EnvVar{
		Name: "NODE_ID",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	})
	args := []string{
		"$(NODE_ID)",
		fmt.Sprintf("%s/%s", configPath, clusterConfigFile),
		fmt.Sprintf("%s/%s", configPath, protocolConfigFile),
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      cluster.Name,
			Labels:    cluster.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: cluster.Labels,
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cluster.Name,
					Namespace: cluster.Namespace,
					Labels:    cluster.Labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            storage.Name,
							Image:           storage.Spec.Image,
							ImagePullPolicy: storage.Spec.ImagePullPolicy,
							Args:            args,
							Env:             env,
						},
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(cluster, dep, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.TODO(), dep)
}

func (r *Reconciler) addConfigMap(cluster *v1beta2.Cluster, storage *v1beta1.CacheStorage) error {
	log.Info("Creating ConfigMap", "Name", cluster.Name, "Namespace", cluster.Namespace)

	config, err := newClusterConfig(cluster)
	if err != nil {
		return err
	}

	marshaller := jsonpb.Marshaler{}
	data, err := marshaller.MarshalToString(config)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: cluster.Namespace,
			Name:      cluster.Name,
		},
		Data: map[string]string{
			clusterConfigFile: data,
		},
	}
	if err := controllerutil.SetControllerReference(cluster, cm, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.TODO(), cm)
}
