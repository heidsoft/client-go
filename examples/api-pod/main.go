package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

// https://iximiuz.com/en/posts/kubernetes-api-go-types-and-common-machinery/
// https://github.com/iximiuz/client-go-examples/blob/main/serialize-typed-json/main.go
func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取pod列表失败，错误是%s \n", err.Error())
	}
	fmt.Println("pod 列表如下:")
	for _, pod := range pods.Items {
		fmt.Printf("%+v \n", pod)
	}

	deployments, err := clientset.AppsV1().Deployments("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("获取deployment列表失败，错误是%s \n", err.Error())
	}
	fmt.Println("deployment 列表如下:")
	for _, deploy := range deployments.Items {
		fmt.Printf("%s\n", deploy.Name)
	}
}
