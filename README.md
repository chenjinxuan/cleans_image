# cleans_image

运行参数

参数|默认值|说明
-|-|-
once | false | 是否只执行一次
cron|"0 0 3 * MON *"|cron表达式，定时执行
force|true|是否强制删除镜像（不被使用，不被依赖的镜像都会被删除）
tryRun| false | 尝试运行，会列出需要删除的镜像ID,但不执行删除操作

### 本地运行
golang 版本1.13.3
```
git clone https://github.com/chenjinxuan/cleans_image.git
cd cleans_image
go run main.go --once=true
```

### docker运行
```
docker run -it -d --name=clean-image -v /var/run/docker.sock:/var/run/docker.sock  jinxuan/cleans_image:0.0.1 --once=true
```

### k8s下运行
```
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: cleans-image
  labels:
    app: cleans-image
spec:
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: cleans-image
    spec:
      tolerations:
        - key: node-role.kubernetes.io/master
          effect: NoSchedule
      containers:
        - name: cleans-image
          image: jinxuan/cleans_image:0.0.1
          args:
            - "--cron=0 0 0 * * *"
          volumeMounts:
            - name: sock
              mountPath: /var/run/docker.sock
      terminationGracePeriodSeconds: 30
      volumes:
        - name: sock
          hostPath:
            path: /var/run/docker.sock
```