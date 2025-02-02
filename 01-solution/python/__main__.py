import pulumi
import pulumi_kubernetes as kubernetes

config = pulumi.Config()
title = config.get("nginxTitle") 
body = config.get("nginxBody")
kubeconfig = pulumi.StackReference("kubeconfig", stack_name="dirien/00-solution/dev")
do_k8_s_provider = kubernetes.Provider("do_k8s_provider",
    enable_server_side_apply=True,
    kubeconfig=kubeconfig.outputs["kubeconfig"])
nginx_config_map = kubernetes.core.v1.ConfigMap("nginxConfigMap",
    metadata={
        "name": "nginx-config",
    },
    data={
        "index.html": f"""<!DOCTYPE html>
<html>
<head>
    <title>{title}</title>
</head>
<body>
    <h1>{body}</h1>
</body>
</html>
""",
    },
    opts = pulumi.ResourceOptions(provider=do_k8_s_provider))
nginx_deployment = kubernetes.apps.v1.Deployment("nginxDeployment",
    metadata={
        "name": "nginx-deployment",
        "annotations": {
            "reloader.stakater.com/auto": "true",
        },
    },
    spec={
        "replicas": 1,
        "selector": {
            "match_labels": {
                "app": "nginx",
            },
        },
        "template": {
            "metadata": {
                "labels": {
                    "app": "nginx",
                },
            },
            "spec": {
                "containers": [{
                    "name": "nginx",
                    "image": "nginx:latest",
                    "ports": [{
                        "container_port": 80,
                    }],
                    "volume_mounts": [{
                        "name": "nginx-html",
                        "mount_path": "/usr/share/nginx/html/index.html",
                        "sub_path": "index.html",
                    }],
                }],
                "volumes": [{
                    "name": "nginx-html",
                    "config_map": {
                        "name": "nginx-config",
                        "items": [{
                            "key": "index.html",
                            "path": "index.html",
                        }],
                    },
                }],
            },
        },
    },
    opts = pulumi.ResourceOptions(provider=do_k8_s_provider))
nginx_service = kubernetes.core.v1.Service("nginxService",
    metadata={
        "name": "nginx-service",
    },
    spec={
        "selector": {
            "app": "nginx",
        },
        "type": kubernetes.core.v1.ServiceSpecType.LOAD_BALANCER,
        "ports": [{
            "port": 8080,
            "target_port": 80,
        }],
    },
    opts = pulumi.ResourceOptions(provider=do_k8_s_provider))
reloader = kubernetes.helm.v3.Release("reloader",
    chart="reloader",
    namespace="reloader",
    create_namespace=True,
    repository_opts={
        "repo": "https://stakater.github.io/stakater-charts",
    },
    values={
        "reloader": {
            "reloadOnCreate": True,
        },
    },
    opts = pulumi.ResourceOptions(provider=do_k8_s_provider))
pulumi.export("serviceName", nginx_service.metadata.name)
