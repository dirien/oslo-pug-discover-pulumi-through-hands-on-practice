[![Join us on Slack!](docs/static/media/slack.svg)](https://slack.pulumi.com/)
<br/>
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://github.com/dirien/cfgmgmtcamp-2025-discover-pulumi-through-hands-on-practice.git)

<br/>

# Discover Pulumi Through Hands-On Practice

## Welcome, Pulumi Friends! 👋

The goal of this workshop to get started with Pulumi to deploy and manage infrastructure in the cloud.
Although Pulumi works with any cloud provider such as AWS, Azure and Google Cloud, for the purposes of this workshop 
we will be sticking with [DigitalOcean](https://www.digitalocean.com/) where we will learn how to deploy clustors and add
managed services to them, all via the Pulumi CLI. 

The workshop is split in small chapters to get you started quickly from installing the CLI all the way to deploying 
and managing nginx services on DigitalOcean.

### Repository

You can find the repository for this
workshop [here](https://github.com/dirien/cfgmgmtcamp-2025-discover-pulumi-through-hands-on-practice.git). Please feel
free to look for the examples in the different chapters if you get stuck.

### Content

- [Chapter 0 - Hello, DigitalOcean World!](./00-create-kubernetes-cluster.md)
- [Chapter 1 - Deploy the Application to Kubernetes](./01-create-nginx-deployment.md)
- [Chapter 2 - Housekeeping!](./02-housekeeping.md)

### Prerequisites

- [A GitHub Account](https://github.com/signup)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) - optional

There is also a [devcontainer.json](.devcontainer/devcontainer.json) file in this repository which you can use to spin
up a `devcontainer` with all the tools installed. Highly recommended if you are
using [VSCode](https://code.visualstudio.com/docs/devcontainers/containers), [GitHub Codespaces](https://docs.github.com/en/codespaces/overview)

### Troubleshooting Tips

If you encounter any challenges during the workshops, consider the following steps in order:

1. Don't hesitate to reach out to me! I'm always here to assist and get you back on track.
1. Review the example code available [here](https://github.com/dirien/cfgmgmtcamp-2025-discover-pulumi-through-hands-on-practice.git).
1. Search for the error on Google. Honestly, this method often provides the most insightful solutions.
1. Engage with the Pulumi Community on Slack. If you haven't joined yet, you can do
   so [here](https://slack.pulumi.com/).

### Want to know more?

If you enjoyed this workshop, please some of Pulumi's other [learning materials](https://www.pulumi.com/learn/)
