package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net/http"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "github.com/hashicorp/consul/api"
)

func main() {
    // Load the kubeconfig file to connect to the Kubernetes cluster
    kubeconfig := filepath.Join(
        homeDir(), ".kube", "config",
    )
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatalf("Failed to load kubeconfig: %v", err)
    }

    // Create a new Kubernetes client
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create Kubernetes client: %v", err)
    }

    // Discover services in the default namespace
    services, err := clientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Failed to discover services: %v", err)
    }

    for _, service := range services.Items {
        fmt.Printf("Discovered service: %s in namespace: %s\n", service.Name, service.Namespace)
        for _, port := range service.Spec.Ports {
            fmt.Printf("  Port: %d, Protocol: %s\n", port.Port, port.Protocol)
        }
    }

    // Create a new Consul client
    client, err := api.NewClient(api.DefaultConfig())
    if err != nil {
        log.Fatalf("Failed to create Consul client: %v", err)
    }

    // Register the service with Consul
    serviceID := "example-service"
    serviceName := "example-service"
    servicePort := 8080

    registration := &api.AgentServiceRegistration{
        ID:      serviceID,
        Name:    serviceName,
        Port:    servicePort,
        Address: "localhost",
    }

    err = client.Agent().ServiceRegister(registration)
    if err != nil {
        log.Fatalf("Failed to register service with Consul: %v", err)
    }

    fmt.Printf("Service %s registered with Consul\n", serviceName)

    // Discover services
    consulServices, err := client.Agent().Services()
    if err != nil {
        log.Fatalf("Failed to discover services: %v", err)
    }

    for _, service := range consulServices {
        fmt.Printf("Discovered service: %s at %s:%d\n", service.Service, service.Address, service.Port)
    }

    // Start an HTTP server
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from %s!", serviceName)
    })

    log.Printf("Starting HTTP server on port %d\n", servicePort)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", servicePort), nil))
}

// homeDir returns the home directory for the current user
func homeDir() string {
    if h := flag.Lookup("home").Value.String(); h != "" {
        return h
    }
    return
}