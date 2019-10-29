package main

import (
	"fmt"
	nicktmingv1 "github.com/nicktming/k8s-crd-controller/pkg/apis/example.com/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	nicktmingclientset "github.com/nicktming/k8s-crd-controller/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	nicktminginformers "github.com/nicktming/k8s-crd-controller/pkg/client/informers/externalversions"
	"k8s.io/client-go/tools/cache"
)

func main()  {
	config := &rest.Config{
		Host: "http://172.21.0.16:8080",
	}
	// 生成client
	databaseClient, _ := nicktmingclientset.NewForConfig(config)

	// 从api-server中获取
	myDatabase, _ := databaseClient.NicktmingV1().Databases("default").Get("my-database", metav1.GetOptions{})
	fmt.Printf("===>Database Name:%v(%v,%v,%v)\n", myDatabase.Name, myDatabase.Spec.User, myDatabase.Spec.Password, myDatabase.Spec.Encoding)


	factory := nicktminginformers.NewSharedInformerFactory(databaseClient, 10)

	// 添加event handler
	databaseInformer := factory.Nicktming().V1().Databases()
	databaseInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) {fmt.Printf("add: %v\n", obj.(*nicktmingv1.Database).Name)},
		UpdateFunc: func(oldObj, newObj interface{}) {fmt.Printf("update: %v\n", newObj.(*nicktmingv1.Database).Name)},
		DeleteFunc: func(obj interface{}){fmt.Printf("delete: %v\n", obj.(*nicktmingv1.Database).Name)},
	})

	// 启动
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)


	// 从本地缓存中获取元素
	databaseLister := databaseInformer.Lister()

	allDatabases, _ := databaseLister.List(labels.Everything())
	for _, p := range allDatabases {
		fmt.Printf("list database: %v\n", p.Name)
	}
	<- stopCh

}
