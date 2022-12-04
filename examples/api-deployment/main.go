package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// https://iximiuz.com/en/posts/kubernetes-api-go-types-and-common-machinery/
// https://github.com/iximiuz/client-go-examples/blob/main/serialize-typed-json/main.go
func main() {
	deployment := appsv1.Deployment{
		// 类型设置,没有设置时出现反序列化失败，需要注意
		// https://stackoverflow.com/questions/43462908/unable-to-declare-kind-type-for-kubernetes-api-type-declarations
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "web", Image: "nginx:1.21"},
					},
				},
			},
		},
	}

	fmt.Printf("%#v\n", &deployment)

	encoder := jsonserializer.NewSerializerWithOptions(
		nil, // jsonserializer.MetaFactory
		nil, // runtime.ObjectCreater
		nil, // runtime.ObjectTyper
		jsonserializer.SerializerOptions{
			Yaml:   false,
			Pretty: true,
			Strict: false,
		},
	)

	// Runtime.Encode() is just a helper function to invoke Encoder.Encode()
	// 对象序列化
	encoded, err := runtime.Encode(encoder, &deployment)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Deployment 资源对象序列号方法1 \n", string(encoded))

	encoded2, err := json.Marshal(deployment)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Deployment 资源对象序列号方法2 \n", string(encoded2))

	decoder := jsonserializer.NewSerializerWithOptions(
		jsonserializer.DefaultMetaFactory, // jsonserializer.MetaFactory
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectCreater
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectTyper
		jsonserializer.SerializerOptions{
			Yaml:   false,
			Pretty: false,
			Strict: false,
		},
	)

	// The actual decoding is much like stdlib encoding/json.Unmarshal but with some
	// minor tweaks - see https://github.com/kubernetes-sigs/json for more.
	// 反序列化
	decoded, err := runtime.Decode(decoder, encoded)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("反序列化 %#v\n", decoded)
}
