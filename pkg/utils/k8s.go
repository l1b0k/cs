// Copyright 2021 l1b0k
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func init() {
	pflag.String("kubeconfig", "", "Paths to a kubeconfig. Only required if out-of-cluster.")
}

func ClientFromKubeConfig(conf []byte) (kubernetes.Interface, error) {
	rest, err := clientcmd.RESTConfigFromKubeConfig(conf)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(rest)
}

func ClientOrDie() kubernetes.Interface {
	rest, err := clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	if err != nil {
		fmt.Printf("failed build k8s client, %v", err)
		os.Exit(1)
	}
	return kubernetes.NewForConfigOrDie(rest)
}

func MergeApiConfig(newAPIConfig *clientcmdapi.Config, override bool) (*clientcmdapi.Config, error) {
	path := viper.GetString("kubeconfig")
	_, err := os.Stat(filepath.Dir(path))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error read %s ,%w", path, err)
		}
		_ = os.MkdirAll(filepath.Dir(path), 0750)
	}
	_, err = os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error read %s ,%w", path, err)
		}
		_ = os.WriteFile(path, nil, os.ModePerm)
	}

	oldApiConfig, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, err
	}
	newCtx, ok := newAPIConfig.Contexts[newAPIConfig.CurrentContext]
	if !ok {
		return nil, fmt.Errorf("error get ctx from config")
	}
	_, ok = oldApiConfig.Clusters[newCtx.Cluster]
	if ok && !override {
		return nil, fmt.Errorf("conflict clister %s", newCtx.Cluster)
	}

	oldApiConfig.Clusters[newCtx.Cluster] = newAPIConfig.Clusters[newCtx.Cluster]
	oldApiConfig.AuthInfos[newCtx.AuthInfo] = newAPIConfig.AuthInfos[newCtx.AuthInfo]
	oldApiConfig.Contexts[newAPIConfig.CurrentContext] = newCtx
	oldApiConfig.CurrentContext = newAPIConfig.CurrentContext
	err = clientcmd.WriteToFile(*oldApiConfig, path)
	if err != nil {
		return nil, err
	}

	return clientcmd.LoadFromFile(path)
}

func SetKubeConfigName(conf []byte, clusterID, name string) (*clientcmdapi.Config, error) {
	apiConfig, err := clientcmd.Load(conf)
	if err != nil {
		return nil, err
	}

	oldCtx, ok := apiConfig.Contexts[apiConfig.CurrentContext]
	if !ok {
		return nil, fmt.Errorf("error get context %s form kubeconfig", apiConfig.CurrentContext)
	}
	oldCluster, ok := apiConfig.Clusters[oldCtx.Cluster]
	if !ok {
		return nil, fmt.Errorf("error get cluster %s form kubeconfig", oldCtx.Cluster)
	}
	oldAuthInfo, ok := apiConfig.AuthInfos[oldCtx.AuthInfo]
	if !ok {
		return nil, fmt.Errorf("error get authinfo %s form kubeconfig", oldCtx.AuthInfo)
	}

	cfg := clientcmdapi.NewConfig()
	cfg.CurrentContext = name

	newCtx := clientcmdapi.NewContext()
	newCtx.Cluster = fmt.Sprintf("cluster-%s", clusterID)
	newCtx.AuthInfo = fmt.Sprintf("user-%s", clusterID)
	newCtx.Namespace = oldCtx.Namespace
	cfg.Contexts[name] = newCtx

	cfg.Clusters[fmt.Sprintf("cluster-%s", clusterID)] = oldCluster
	cfg.AuthInfos[fmt.Sprintf("user-%s", clusterID)] = oldAuthInfo

	return cfg, nil
}

func ParseInstanceID(node *corev1.Node) string {
	keys := strings.Split(node.Spec.ProviderID, ".")
	if len(keys) > 1 {
		return keys[1]
	}
	return ""
}

func ParseInternalIP(node *corev1.Node) net.IP {
	for _, a := range node.Status.Addresses {
		if a.Type != corev1.NodeInternalIP {
			continue
		}
		return net.ParseIP(a.Address)
	}
	return nil
}

func ParseRegion(node *corev1.Node) string {
	id, _ := node.Labels[corev1.LabelTopologyRegion]
	return id
}

func ParseZone(node *corev1.Node) string {
	id, _ := node.Labels[corev1.LabelTopologyZone]
	return id
}

func ParseInstanceType(node *corev1.Node) string {
	id, _ := node.Labels[corev1.LabelInstanceTypeStable]
	return id
}
