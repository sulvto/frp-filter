FROM centos:7

# 清理原有的仓库配置文件
RUN mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup

# 创建新的仓库配置文件
RUN tee /etc/yum.repos.d/Centos-Base.repo <<EOF
[base]
name=CentOS-\$releasever - Base
baseurl=http://mirrors.aliyun.com/centos/\$releasever/os/\$basearch/
gpgcheck=1
enabled=1
gpgkey=https://mirrors.aliyun.com/repo/Centos-\$releasever-keyring.gpg
EOF

ENV http_proxy=
ENV https_proxy=
ENV HTTP_PROXY=
ENV HTTPS_PROXY=

RUN yum --nogpgcheck install -y gcc-c++ make wget curl && \
    wget https://dl.google.com/go/go1.18.linux-amd64.tar.gz && \
    tar -xvf go1.18.linux-amd64.tar.gz -C /usr/local && \
    rm go1.18.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:$PATH"

WORKDIR /app

COPY /app /
