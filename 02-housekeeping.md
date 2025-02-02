# Chapter 2: Housekeeping!

In this chapter, we'll undertake some cleanup. We will delete all resources we set up during the workshop.

## Instructions

### Step 1 - Destroy your cluster with Pulumi

Depending on your progress in the workshop, you have to destroy more or less stacks. We will start with the last stack
we created, which is the `01-create-nginx-deployment` stack. Please change into the directory and destroy the stack:

```bash
cd 01-create-nginx-deployment
pulumi destroy -y -f
```

Now we can destroy the infrastructure stack in the `00-create-kubernetes-cluster` folder:

```bash
cd 00-create-kubernetes-cluster
pulumi destroy -y -f
```

### Step 2 - Now Celebrate, You're Done!

<p align="center">
  <img src="https://cdn.dribbble.com/users/234969/screenshots/5414177/burst_trophy_dribbble.gif">
</p>
