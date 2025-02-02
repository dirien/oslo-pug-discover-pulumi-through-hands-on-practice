package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config := config.New(ctx, "")
		title := config.Get("nginxTitle")
		body := config.Get("nginxBody")
		kubeconfig, err := pulumi.NewStackReference(ctx, "kubeconfig", &pulumi.StackReferenceArgs{
			Name: pulumi.String("dirien/00-solution/dev"),
		})
		if err != nil {
			return err
		}
		doK8SProvider, err := kubernetes.NewProvider(ctx, "do_k8s_provider", &kubernetes.ProviderArgs{
			EnableServerSideApply: pulumi.Bool(true),
			Kubeconfig:            kubeconfig.GetStringOutput(pulumi.String("kubeconfig")),
		})
		if err != nil {
			return err
		}
		_, err = corev1.NewConfigMap(ctx, "nginxConfigMap", &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("nginx-config"),
			},
			Data: pulumi.StringMap{
				"index.html": pulumi.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%v</title>
</head>
<body>
    <h1>%v</h1>
</body>
</html>
`, pulumi.String(title), pulumi.String(body)),
			},
		}, pulumi.Provider(doK8SProvider))
		if err != nil {
			return err
		}
		_, err = appsv1.NewDeployment(ctx, "nginxDeployment", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("nginx-deployment"),
				Annotations: pulumi.StringMap{
					"reloader.stakater.com/auto": pulumi.String("true"),
				},
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(1),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("nginx"),
								Image: pulumi.String("nginx:latest"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
									},
								},
								VolumeMounts: corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										Name:      pulumi.String("nginx-html"),
										MountPath: pulumi.String("/usr/share/nginx/html/index.html"),
										SubPath:   pulumi.String("index.html"),
									},
								},
							},
						},
						Volumes: corev1.VolumeArray{
							&corev1.VolumeArgs{
								Name: pulumi.String("nginx-html"),
								ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
									Name: pulumi.String("nginx-config"),
									Items: corev1.KeyToPathArray{
										&corev1.KeyToPathArgs{
											Key:  pulumi.String("index.html"),
											Path: pulumi.String("index.html"),
										},
									},
								},
							},
						},
					},
				},
			},
		}, pulumi.Provider(doK8SProvider))
		if err != nil {
			return err
		}
		nginxService, err := corev1.NewService(ctx, "nginxService", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("nginx-service"),
			},
			Spec: &corev1.ServiceSpecArgs{
				Selector: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
				Type: pulumi.String(corev1.ServiceSpecTypeLoadBalancer),
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(8080),
						TargetPort: pulumi.Any(80),
					},
				},
			},
		}, pulumi.Provider(doK8SProvider))
		if err != nil {
			return err
		}
		_, err = helmv3.NewRelease(ctx, "reloader", &helmv3.ReleaseArgs{
			Chart:           pulumi.String("reloader"),
			Namespace:       pulumi.String("reloader"),
			CreateNamespace: pulumi.Bool(true),
			RepositoryOpts: &helmv3.RepositoryOptsArgs{
				Repo: pulumi.String("https://stakater.github.io/stakater-charts"),
			},
			Values: pulumi.Map{
				"reloader": pulumi.Any(map[string]interface{}{
					"reloadOnCreate": true,
				}),
			},
		}, pulumi.Provider(doK8SProvider))
		if err != nil {
			return err
		}
		ctx.Export("serviceName", nginxService.Metadata.Name())
		return nil
	})
}
