# Scheduler FrameWork demo

This repo is a demo for the feature scheduler extender based on kubernetes 1.16.

## Get Started 

- Deploy
```
kubectl apply -f ./deploy/

Note that your default-scheduler should be started with args:
- --policy-config-file=/etc/kubernetes/scheduler-policy-config.json

There's a scheduler-policy-config.json in config folder for your example.
```

- Build local
```
make local
```

- Build image
```
make image
```

## Deploy a pod using this scheduler

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ex-test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: busybox
      isp: "mobile"
  template:
    metadata:
      labels:
        app: busybox
    spec:
      terminationGracePeriodSeconds: 5
      containers:
      - image: busybox:latest
        imagePullPolicy: IfNotPresent
        name: busybox
        command: ["sleep", "3600"]
```

By default, the pod will be in pending status with the message "2 Nodes doesn't have any isp labels."

Label the node with `kubectl label node $your-node-ip isp="mobile"`, the pod will be in running status.