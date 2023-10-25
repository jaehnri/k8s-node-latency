package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	podName string

	kubeClient *kubernetes.Clientset
}

func NewClient() *Client {
	serviceName := "your-service-name"
	namespace := "your-namespace"

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panic("couldn't fetch the in-cluster config", err)
	}
	// creates the clientset
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic("couldn't create kubernetes config")
	}

	service, err := kubeClient.CoreV1().Services(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Panic("couldn't get server ClusterIP")
	}

	clusterIP := service.Spec.ClusterIP
	fmt.Printf("ClusterIP of service %s in namespace %s is %s\n", serviceName, namespace, clusterIP)

	return &Client{
		kubeClient: kubeClient,
	}
}

func (c *Client) Run() {
	sigint := make(chan os.Signal, 1)
	// interrupt signal sent from terminal
	signal.Notify(sigint, os.Interrupt)
	// sigterm signal sent from kubernetes
	signal.Notify(sigint, syscall.SIGTERM)

	<-sigint
	log.Fatal("SIGINT received, shutting down client")
}
