FROM buildpack-deps:jessie-curl
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
ADD server /server
ENTRYPOINT ["/server"]
