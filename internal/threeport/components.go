package threeport

import (
	"fmt"

	"github.com/threeport/threeport/internal/kube"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// InstallThreeportControlPlaneComponents installs the threeport API and
// controllers.
func InstallThreeportControlPlaneComponents(
	kubeClient dynamic.Interface,
	mapper *meta.RESTMapper,
) error {

	var dbCreateConfig = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "db-create",
				"namespace": "threeport-control-plane",
			},
			"data": map[string]interface{}{
				"db.sql": `CREATE USER IF NOT EXISTS tp_rest_api
  LOGIN
;
CREATE DATABASE IF NOT EXISTS threeport_api
    encoding='utf-8'
;
GRANT ALL ON DATABASE threeport_api TO tp_rest_api;
`,
			},
		},
	}
	if _, err := kube.CreateResource(dbCreateConfig, kubeClient, *mapper); err != nil {
		return fmt.Errorf("failed to create DB create configmap: %w", err)
	}

	var apiSecret = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      "db-config",
				"namespace": ControlPlaneNamespace,
			},
			"stringData": map[string]interface{}{
				"env": `DB_HOST=crdb
DB_USER=tp_rest_api
DB_PASSWORD=tp-rest-api-pwd
DB_NAME=threeport_api
DB_PORT=26257
DB_SSL_MODE=disable
NATS_HOST=nats-js
NATS_PORT=4222
`,
			},
		},
	}
	if _, err := kube.CreateResource(apiSecret, kubeClient, *mapper); err != nil {
		return fmt.Errorf("failed to create API server secret: %w", err)
	}

	var apiDeployment = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "threeport-api-server",
				"namespace": ControlPlaneNamespace,
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app.kubernetes.io/name": "threeport-api-server",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app.kubernetes.io/name": "threeport-api-server",
						},
					},
					"spec": map[string]interface{}{
						"initContainers": []interface{}{
							map[string]interface{}{
								"name":            "db-init",
								"image":           "cockroachdb/cockroach:v22.2.2",
								"imagePullPolicy": "IfNotPresent",
								"command": []interface{}{
									"bash",
									"-c",
									//- "cockroach sql --insecure --host crdb --port 26257 -f /etc/threeport/db-create/db.sql && cockroach sql --insecure --host crdb --port 26257 --database threeport_api -f /etc/threeport/db-load/create_tables.sql && cockroach sql --insecure --host crdb --port 26257 --database threeport_api -f /etc/threeport/db-load/fill_tables.sql"
									"cockroach sql --insecure --host crdb --port 26257 -f /etc/threeport/db-create/db.sql",
								},
								"volumeMounts": []interface{}{
									//- name: db-load
									//  mountPath: "/etc/threeport/db-load"
									map[string]interface{}{
										"name":      "db-create",
										"mountPath": "/etc/threeport/db-create",
									},
								},
							},
						},
						"containers": []interface{}{
							map[string]interface{}{
								"name":            "api-server",
								"image":           "threeport-rest-api-dev:latest",
								"imagePullPolicy": "Never",
								"ports": []interface{}{
									map[string]interface{}{
										"containerPort": 1323,
										"name":          "http",
										"protocol":      "TCP",
									},
								},
								"volumeMounts": []interface{}{
									map[string]interface{}{
										"name":      "db-config",
										"mountPath": "/etc/threeport/",
									},
									map[string]interface{}{
										"name":      "code-path",
										"mountPath": "/threeport",
										//env:
										//- name: AUTHENTIK_HOST
										//  value: "authentik.threeport-api.svc.cluster.local"
										//- name: AUTHENTIK_SCHEME
										//  value: "https"
										//- name: AUTHENTIK_BOOTSTRAP_PASSWORD
										//  valueFrom:
										//    secretKeyRef:
										//      name: authentik-bootstrap
										//      key: authentik-bootstrap-password
										//- name: AUTHENTIK_BOOTSTRAP_TOKEN
										//  valueFrom:
										//    secretKeyRef:
										//      name: authentik-bootstrap
										//      key: authentik-bootstrap-token
									},
								},
							},
						},
						"volumes": []interface{}{
							map[string]interface{}{
								"name": "db-config",
								"secret": map[string]interface{}{
									"secretName": "db-config",
								},
							},
							map[string]interface{}{
								"name": "db-create",
								"configMap": map[string]interface{}{
									"name": "db-create",
								},
							},
							map[string]interface{}{
								"name": "db-load",
								"configMap": map[string]interface{}{
									"name": "db-load",
								},
							},
							map[string]interface{}{
								"name": "code-path",
								"hostPath": map[string]interface{}{
									"path": "/threeport",
									"type": "Directory",
								},
							},
						},
					},
				},
			},
		},
	}
	if _, err := kube.CreateResource(apiDeployment, kubeClient, *mapper); err != nil {
		return fmt.Errorf("failed to create API server deployment: %w", err)
	}

	var apiService = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata": map[string]interface{}{
				"name":      "threeport-api-server",
				"namespace": ControlPlaneNamespace,
			},
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{
					"app.kubernetes.io/name": "threeport-api-server",
				},
				"type": "LoadBalancer",
				"ports": []interface{}{
					map[string]interface{}{
						"name":       "http",
						"port":       80,
						"protocol":   "TCP",
						"targetPort": 1323,
					},
				},
			},
		},
	}
	if _, err := kube.CreateResource(apiService, kubeClient, *mapper); err != nil {
		return fmt.Errorf("failed to create API server service: %w", err)
	}

	var workloadControllerSecret = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      "workload-controller-config",
				"namespace": ControlPlaneNamespace,
			},
			"type": "Opaque",
			"stringData": map[string]interface{}{
				"API_SERVER":      "http://threeport-api-server",
				"MSG_BROKER_HOST": "nats-js",
				"MSG_BROKER_PORT": "4222",
			},
		},
	}
	if _, err := kube.CreateResource(workloadControllerSecret, kubeClient, *mapper); err != nil {
		return fmt.Errorf("failed to create workload controller secret: %w", err)
	}

	var workloadControllerDeployment = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "threeport-workload-controller",
				"namespace": ControlPlaneNamespace,
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app.kubernetes.io/name": "threeport-workload-controller",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app.kubernetes.io/name": "threeport-workload-controller",
						},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":            "workload-controller",
								"image":           "threeport-workload-controller-dev:latest",
								"imagePullPolicy": "IfNotPresent",
								"envFrom": []interface{}{
									map[string]interface{}{
										"secretRef": map[string]interface{}{
											"name": "workload-controller-config",
										},
									},
								},
								"volumeMounts": []interface{}{
									map[string]interface{}{
										"name":      "code-path",
										"mountPath": "/threeport",
									},
								},
							},
						},
						"volumes": []interface{}{
							map[string]interface{}{
								"name": "code-path",
								"hostPath": map[string]interface{}{
									"path": "/threeport",
									"type": "Directory",
								},
							},
						},
					},
				},
			},
		},
	}
	if _, err := kube.CreateResource(workloadControllerDeployment, kubeClient, *mapper); err != nil {
		return fmt.Errorf("failed to create workload controller deployment: %w", err)
	}

	return nil
}