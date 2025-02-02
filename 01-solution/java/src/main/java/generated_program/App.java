package generated_program;

import com.pulumi.Config;
import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.resources.StackReference;
import com.pulumi.resources.StackReferenceArgs;
import com.pulumi.kubernetes.Provider;
import com.pulumi.kubernetes.ProviderArgs;
import com.pulumi.kubernetes.core_v1.ConfigMap;
import com.pulumi.kubernetes.core_v1.ConfigMapArgs;
import com.pulumi.kubernetes.meta_v1.inputs.ObjectMetaArgs;
import com.pulumi.kubernetes.apps_v1.Deployment;
import com.pulumi.kubernetes.apps_v1.DeploymentArgs;
import com.pulumi.kubernetes.apps_v1.inputs.DeploymentSpecArgs;
import com.pulumi.kubernetes.meta_v1.inputs.LabelSelectorArgs;
import com.pulumi.kubernetes.core_v1.inputs.PodTemplateSpecArgs;
import com.pulumi.kubernetes.core_v1.inputs.PodSpecArgs;
import com.pulumi.kubernetes.core_v1.Service;
import com.pulumi.kubernetes.core_v1.ServiceArgs;
import com.pulumi.kubernetes.core_v1.inputs.ServiceSpecArgs;
import com.pulumi.kubernetes.helm.sh_v3.Release;
import com.pulumi.kubernetes.helm.sh_v3.ReleaseArgs;
import com.pulumi.kubernetes.helm.sh_v3.inputs.RepositoryOptsArgs;
import com.pulumi.resources.CustomResourceOptions;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var config = ctx.config();
        var title = config.require("nginxTitle");
        var body = config.require("nginxBody");

        var kubeconfig = new StackReference("kubeconfig", StackReferenceArgs.builder()
            .name("dirien/00-solution/dev")
            .build());

        var doK8SProvider = new Provider("doK8SProvider", ProviderArgs.builder()
            .enableServerSideApply(true)
            .kubeconfig(kubeconfig.requireOutput("kubeconfig"))
            .build());

        var index = String.format("""
            <!DOCTYPE html>
            <html>
            <head>
                <title>%s</title>
            </head>
            <body>
                <h1>%s</h1>
            </body>
            </html>
            """, title, body);
        
        var nginxConfigMap = new ConfigMap("nginxConfigMap", ConfigMapArgs.builder()
            .metadata(ObjectMetaArgs.builder()
                .name("nginx-config")
                .build())
            .data(Map.of("index.html", index))
            .build(), CustomResourceOptions.builder()
                .provider(doK8SProvider)
                .build());

        var nginxDeployment = new Deployment("nginxDeployment", DeploymentArgs.builder()
            .metadata(ObjectMetaArgs.builder()
                .name("nginx-deployment")
                .annotations(Map.of("reloader.stakater.com/auto", "true"))
                .build())
            .spec(DeploymentSpecArgs.builder()
                .replicas(1)
                .selector(LabelSelectorArgs.builder()
                    .matchLabels(Map.of("app", "nginx"))
                    .build())
                .template(PodTemplateSpecArgs.builder()
                    .metadata(ObjectMetaArgs.builder()
                        .labels(Map.of("app", "nginx"))
                        .build())
                    .spec(PodSpecArgs.builder()
                        .containers(ContainerArgs.builder()
                            .name("nginx")
                            .image("nginx:latest")
                            .ports(ContainerPortArgs.builder()
                                .containerPort(80)
                                .build())
                            .volumeMounts(VolumeMountArgs.builder()
                                .name("nginx-html")
                                .mountPath("/usr/share/nginx/html/index.html")
                                .subPath("index.html")
                                .build())
                            .build())
                        .volumes(VolumeArgs.builder()
                            .name("nginx-html")
                            .configMap(ConfigMapVolumeSourceArgs.builder()
                                .name("nginx-config")
                                .items(KeyToPathArgs.builder()
                                    .key("index.html")
                                    .path("index.html")
                                    .build())
                                .build())
                            .build())
                        .build())
                    .build())
                .build())
            .build(), CustomResourceOptions.builder()
                .provider(doK8SProvider)
                .build());

        var nginxService = new Service("nginxService", ServiceArgs.builder()
            .metadata(ObjectMetaArgs.builder()
                .name("nginx-service")
                .build())
            .spec(ServiceSpecArgs.builder()
                .selector(Map.of("app", "nginx"))
                .type("LoadBalancer")
                .ports(ServicePortArgs.builder()
                    .port(8080)
                    .targetPort(80)
                    .build())
                .build())
            .build(), CustomResourceOptions.builder()
                .provider(doK8SProvider)
                .build());

        var reloader = new Release("reloader", ReleaseArgs.builder()
            .chart("reloader")
            .namespace("reloader")
            .createNamespace(true)
            .repositoryOpts(RepositoryOptsArgs.builder()
                .repo("https://stakater.github.io/stakater-charts")
                .build())
            .values(Map.of("reloader", Map.of("reloadOnCreate", true)))
            .build(), CustomResourceOptions.builder()
                .provider(doK8SProvider)
                .build());

        ctx.export("serviceName", nginxService.metadata().applyValue(metadata -> metadata.name()));
    }
}
