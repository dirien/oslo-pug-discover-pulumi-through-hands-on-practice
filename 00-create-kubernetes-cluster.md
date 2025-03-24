# Chapter 0 - Hello, DigitalOcean World!

<img src="docs/static/media/chap0.png">

## Introduction

In this chapter, we will develop a simple Pulumi program that creates a single virtual machine with a basic web server
running on it. Our goal is to become acquainted with the Pulumi CLI, understand the structure of a Pulumi program, and
learn how to create multiple stacks and override default values.

### Modern Infrastructure As Code with Pulumi

Pulumi is an open-source infrastructure-as-code tool for creating, deploying and managing cloud
infrastructure. Pulumi works with traditional infrastructures like VMs, networks, and databases and modern
architectures, including containers, Kubernetes clusters, and serverless functions. Pulumi supports dozens of public,
private, and hybrid cloud service providers.

Pulumi is a multi-language infrastructure as Code tool using imperative languages to create a declarative
infrastructure description.

You have a wide range of programming languages available, and you can use the one you and your team are the most
comfortable with. Currently, (2/2025) Pulumi supports the following languages:

* Node.js (JavaScript / TypeScript)

* Python

* Go

* Java

* .NET (C#, VB, F#)

* YAML


## Installing Pulumi CLI

To start with Pulumi, you need to install the CLI. Go to https://www.pulumi.com/docs/iac/download-install/ 
and follow the instructions that apply to your machine. 

Once you have installed the CLI, verify that it works via the following command:
```shell
> pulumi version
v3.147.0 # or similar
```
You are now ready to create your first Pulumi application. 

## Create a Pulumi account

When using infrastructure as code tools, we typically keep track of deployed resources in some sort of _state_ backend. 
Pulumi CLI saves the state of your program in our own Pulumi Cloud backend where you can manage your projects, collaborate with teammates and use shared secrets.

The Pulumi Cloud backend is optional and can be replaced with other state backends such as an S3 bucket on AWS, a blob storage on Azure or even locally on your machine
but in this workshop we will be using Pulumi Cloud. 

Start by running `pulumi login` and follow the instructions to connect the CLI with your Pulumi Cloud account after you've created it.

> **Important:**
> If you run Pulumi for the first time, you will be asked to log in. Follow the instructions on the screen to
> login. You may need to create an account first, don't worry it is free.

### Have a look at the Pulumi Cloud Console

After you have logged in, you can have a look at the Pulumi Cloud Console at https://app.pulumi.com/ to see your projects and stacks.

Go click around and see what you can find. You will see the projects you create in this workshop here.

Also, there is the tab item `Environmet` tab, which we will use later in the workshop to create different Pulumi ESC environments.

## Creating A New Project

Once the CLI is installed and your account connected, you can create a new project
```
pulumi new <language>
```
Where `<language>` is your language of choice for writing your infrastructure code.

You will need to run this command in an empty directory and `cd` into it so the command become:
```shell
mkdir digitalocean-pulumi-workshop
cd digitalocean-pulumi-workshop
pulumi new typescript
```
The snippet above creates an empty project in the TypeScript language, but you can choose another language such as 
`go`, `python`, `csharp`, `go` or `yaml`. Pulumi supports all of these. We will have snippets in all of these languages.

```yaml
name: 00-solution-do-yaml
description: A minimal DigitalOcean Pulumi YAML program
runtime: yaml

variables:
  clusterRegion: "fra1"
  nodePoolName: "default"
  nodeCount: 1
  version: 1.32.2-do.0

resources:
  do-cluster:
    type: digitalocean:KubernetesCluster
    properties:
      region: "${clusterRegion}"
      version: "${version}"
      destroyAllAssociatedResources: true
      nodePool:
        name: "${nodePoolName}"
        size: "s-2vcpu-2gb"
        nodeCount: "${nodeCount}"

outputs:
  name: "${do-cluster.name}"
  kubeconfig: "${do-cluster.kubeConfigs[0].rawConfig}"
```

<details>
  <summary>CSharp</summary>
    {% highlight csharp %}
    
    
    using System.Collections.Generic;
    using System.Linq;
    using Pulumi;
    using DigitalOcean = Pulumi.DigitalOcean;
    
    return await Deployment.RunAsync(() => 
    {
        var clusterRegion = "fra1";
        var nodePoolName = "default";
        var nodeCount = 1;
        var version = "1.32.2-do.0";
    
        var doCluster = new DigitalOcean.KubernetesCluster("do-cluster", new()
        {
            Region = clusterRegion,
            Version = version,
            DestroyAllAssociatedResources = true,
            NodePool = new DigitalOcean.Inputs.KubernetesClusterNodePoolArgs
            {
                Name = nodePoolName,
                Size = "s-2vcpu-2gb",
                NodeCount = nodeCount,
            },
        });
    
        return new Dictionary<string, object?>
        {
            ["name"] = doCluster.Name,
            ["kubeconfig"] = doCluster.KubeConfigs.Apply(kubeConfigs => kubeConfigs[0].RawConfig),
        };
    });
    {% endhighlight %}
</details>

<details>
  <summary>Go</summary>
    {% highlight go %}
    package main
    
    import (
        "github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
        "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    )
    
    func main() {
        pulumi.Run(func(ctx *pulumi.Context) error {
            clusterRegion := "fra1"
            nodePoolName := "default"
            nodeCount := 1
            version := "1.32.2-do.0"
            doCluster, err := digitalocean.NewKubernetesCluster(ctx, "do-cluster", &digitalocean.KubernetesClusterArgs{
                Region:                        pulumi.String(clusterRegion),
                Version:                       pulumi.String(version),
                DestroyAllAssociatedResources: pulumi.Bool(true),
                NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
                    Name:      pulumi.String(nodePoolName),
                    Size:      pulumi.String("s-2vcpu-2gb"),
                    NodeCount: pulumi.IntPtr(nodeCount),
                },
            })
            if err != nil {
                return err
            }
            ctx.Export("name", doCluster.Name)
            ctx.Export("kubeconfig", doCluster.KubeConfigs.ApplyT(func(kubeConfigs []digitalocean.KubernetesClusterKubeConfig) (*string, error) {
                return kubeConfigs[0].RawConfig, nil
            }).(pulumi.StringPtrOutput))
            return nil
        })
    }
    {% endhighlight %}
</details>

<details>
  <summary>Java</summary>
    {% highlight java %}
        package application;
        
        import com.pulumi.Context;
        import com.pulumi.Pulumi;
        import com.pulumi.digitalocean.KubernetesCluster;
        import com.pulumi.digitalocean.KubernetesClusterArgs;
        import com.pulumi.digitalocean.inputs.KubernetesClusterNodePoolArgs;
        
        public class App {
            public static void main(String[] args) {
                Pulumi.run(App::stack);
            }
        
            public static void stack(Context ctx) {
                final var clusterRegion = "fra1";
                final var nodePoolName = "default";
                final var nodeCount = 1;
                final var version = "1.32.2-do.0";
        
                var doCluster = new KubernetesCluster("do-cluster", KubernetesClusterArgs.builder()
                    .region(clusterRegion)
                    .version(version)
                    .destroyAllAssociatedResources(true)
                    .nodePool(KubernetesClusterNodePoolArgs.builder()
                        .name(nodePoolName)
                        .size("s-2vcpu-2gb")
                        .nodeCount(nodeCount)
                        .build())
                    .build());
        
                ctx.export("name", doCluster.name());
                ctx.export("kubeconfig", doCluster.kubeConfigs().applyValue(kubeConfigs -> kubeConfigs.get(0).rawConfig()));
            }
        }
    {% endhighlight %}
</details>


<details>
  <summary>Python</summary>
    {% highlight python %}
        import pulumi
        import pulumi_digitalocean as digitalocean
        
        cluster_region = "fra1"
        node_pool_name = "default"
        node_count = 1
        version = "1.32.2-do.0"
        do_cluster = digitalocean.KubernetesCluster("do-cluster",
            region=cluster_region,
            version=version,
            destroy_all_associated_resources=True,
            node_pool={
                "name": node_pool_name,
                "size": "s-2vcpu-2gb",
                "node_count": node_count,
            })
        pulumi.export("name", do_cluster.name)
        pulumi.export("kubeconfig", do_cluster.kube_configs[0].raw_config)
    {% endhighlight %}
</details>

<details>
  <summary>TypeScript</summary>
    {% highlight typescript %}
    import * as pulumi from "@pulumi/pulumi";
    import * as digitalocean from "@pulumi/digitalocean";
    
    const clusterRegion = "fra1";
    const nodePoolName = "default";
    const nodeCount = 1;
    const version = "1.32.2-do.0";
    const doCluster = new digitalocean.KubernetesCluster("do-cluster", {
        region: clusterRegion,
        version: version,
        destroyAllAssociatedResources: true,
        nodePool: {
            name: nodePoolName,
            size: "s-2vcpu-2gb",
            nodeCount: nodeCount,
        },
    });
    export const name = doCluster.name;
    export const kubeconfig = doCluster.kubeConfigs.apply(kubeConfigs => kubeConfigs[0].rawConfig);
    {% endhighlight %}
</details>

## Deploying the Cluster

To deploy the now created cluster, run the following command:

```shell
pulumi up
```

This command will show you a preview of all the resources and asks you if you want to deploy them. You can run dedicated
commands to see the preview or to deploy the resources.

```shell
pulumi preview
# or
pulumi up
```

> **Important:**
> You need to set your DigitalOcean API token as an environment variable `DIGITALOCEAN_TOKEN` before running the `pulumi up` command.
> The token will be provided to you in the workshop.

## Accessing the Cluster

After the deployment is complete, you can access the cluster by running the following command:

```shell
pulumi stack output kubeconfig --show-secrets > kubeconfig
export KUBECONFIG=$PWD/kubeconfig
kubectl get nodes
```

You should see the nodes in the cluster.

## Stretch Goal

- Create a new stack with a different region and node count. Use the `pulumi stack` command to create a new stack and
  deploy it with the new values.

## Learn more

- [Pulumi](https://www.pulumi.com/)
- [Pulumi Digitalocean Provider](https://www.pulumi.com/registry/packages/digitalocean/)
- [Pulumi ESC](https://www.pulumi.com/docs/pulumi-cloud/esc/)
