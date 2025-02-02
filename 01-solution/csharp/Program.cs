using System.Collections.Generic;
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var title = config.Get("nginxTitle");
    var body = config.Get("nginxBody");

    var kubeconfig = new StackReference("kubeconfig", new()
    {
        Name = "dirien/00-solution/dev",
    });

    var doK8SProvider = new Kubernetes.Provider("do_k8s_provider", new()
    {
        EnableServerSideApply = true,
        KubeConfig = kubeconfig.RequireOutput("kubeconfig").Apply(kubeconfig => 
        {
            return kubeconfig.ToString() ?? "";   
        }),
    });

    var nginxConfigMap = new Kubernetes.Core.V1.ConfigMap("nginxConfigMap", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Name = "nginx-config",
        },
        Data = 
        {
            { "index.html", @$"<!DOCTYPE html>
<html>
<head>
    <title>{title}</title>
</head>
<body>
    <h1>{body}</h1>
</body>
</html>
" },
        },
    }, new CustomResourceOptions
    {
        Provider = doK8SProvider,
    });

    var nginxDeployment = new Kubernetes.Apps.V1.Deployment("nginxDeployment", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Name = "nginx-deployment",
            Annotations = 
            {
                { "reloader.stakater.com/auto", "true" },
            },
        },
        Spec = new Kubernetes.Types.Inputs.Apps.V1.DeploymentSpecArgs
        {
            Replicas = 1,
            Selector = new Kubernetes.Types.Inputs.Meta.V1.LabelSelectorArgs
            {
                MatchLabels = 
                {
                    { "app", "nginx" },
                },
            },
            Template = new Kubernetes.Types.Inputs.Core.V1.PodTemplateSpecArgs
            {
                Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
                {
                    Labels = 
                    {
                        { "app", "nginx" },
                    },
                },
                Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
                {
                    Containers = new[]
                    {
                        new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                        {
                            Name = "nginx",
                            Image = "nginx:latest",
                            Ports = new[]
                            {
                                new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                                {
                                    ContainerPortValue = 80,
                                },
                            },
                            VolumeMounts = new[]
                            {
                                new Kubernetes.Types.Inputs.Core.V1.VolumeMountArgs
                                {
                                    Name = "nginx-html",
                                    MountPath = "/usr/share/nginx/html/index.html",
                                    SubPath = "index.html",
                                },
                            },
                        },
                    },
                    Volumes = new[]
                    {
                        new Kubernetes.Types.Inputs.Core.V1.VolumeArgs
                        {
                            Name = "nginx-html",
                            ConfigMap = new Kubernetes.Types.Inputs.Core.V1.ConfigMapVolumeSourceArgs
                            {
                                Name = "nginx-config",
                                Items = new[]
                                {
                                    new Kubernetes.Types.Inputs.Core.V1.KeyToPathArgs
                                    {
                                        Key = "index.html",
                                        Path = "index.html",
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }, new CustomResourceOptions
    {
        Provider = doK8SProvider,
    });

    var nginxService = new Kubernetes.Core.V1.Service("nginxService", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Name = "nginx-service",
        },
        Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
        {
            Selector = 
            {
                { "app", "nginx" },
            },
            Type = Kubernetes.Core.V1.ServiceSpecType.LoadBalancer,
            Ports = new[]
            {
                new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                {
                    Port = 8080,
                    TargetPort = 80,
                },
            },
        },
    }, new CustomResourceOptions
    {
        Provider = doK8SProvider,
    });

    var reloader = new Kubernetes.Helm.V3.Release("reloader", new()
    {
        Chart = "reloader",
        Namespace = "reloader",
        CreateNamespace = true,
        RepositoryOpts = new Kubernetes.Types.Inputs.Helm.V3.RepositoryOptsArgs
        {
            Repo = "https://stakater.github.io/stakater-charts",
        },
        Values = 
        {
            { "reloader", new Dictionary<string, object?>
            {
                ["reloadOnCreate"] = true,
            } },
        },
    }, new CustomResourceOptions
    {
        Provider = doK8SProvider,
    });

    return new Dictionary<string, object?>
    {
        ["serviceName"] = nginxService.Metadata.Apply(metadata => metadata.Name),
    };
});

