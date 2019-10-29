package main

import (
	"k8s.io/client-go/rest"
	clientset "github.com/resouer/k8s-crd-controller/pkg/client/clientset/versioned"
)

func main()  {
	config := &rest.Config{
		Host: "http://172.21.0.16:8080",
	}

	databaseClient := clientset
}
