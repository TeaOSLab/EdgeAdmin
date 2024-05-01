FROM --platform=linux/amd64 alpine:latest
LABEL maintainer="goedge.cdn@gmail.com"
ENV TZ "Asia/Shanghai"
ENV VERSION 1.3.8
ENV ROOT_DIR /usr/local/goedge
ENV TAR_FILE edge-admin-linux-amd64-plus-v${VERSION}.zip

# remote official repository
ENV TAR_URL "https://dl.goedge.cn/edge/v${VERSION}/edge-admin-linux-amd64-plus-v${VERSION}.zip"

# your local repository
#ENV TAR_URL "http://192.168.2.61:8080/edge-admin-linux-amd64-plus-v${VERSION}.zip"

RUN apk add --no-cache tzdata

RUN apk add wget
RUN mkdir ${ROOT_DIR}; \
	cd ${ROOT_DIR}; \
	wget ${TAR_URL} -O ${TAR_FILE}; \
	apk add unzip; \
	unzip ${TAR_FILE}; \
	rm -f ${TAR_FILE}

RUN apk add mysql mysql-client; \
	sed -e "s/\[mysqld\]/\[mysqld\]\n\ndatadir=\/var\/lib\/mysql\nport=3306\ninnodb_flush_log_at_trx_commit=2\nmax_connections=256\nmax_prepared_stmt_count=65535\nbinlog_cache_size=1M\nbinlog_stmt_cache_size=1M\nthread_cache_size=32\nbinlog_expire_logs_seconds=1209600\n\n/" /etc/my.cnf > /tmp/my.cnf; \
    cp /tmp/my.cnf /etc/my.cnf; \
    sed -e "s/skip-networking/#skip-networking/" /etc/my.cnf.d/mariadb-server.cnf > /tmp/mariadb-server.cnf; \
    cp /tmp/mariadb-server.cnf /etc/my.cnf.d/mariadb-server.cnf; \
    mysql_install_db --user=mysql
RUN mysqld_safe --user=mysql & \
    sleep 5; \
    mysql -uroot -hlocalhost --execute="ALTER USER 'root'@'localhost' IDENTIFIED BY '123456';"

RUN echo -e "#!/usr/bin/env sh\n\nmysqld_safe --user=mysql &\n/usr/local/goedge/edge-admin/bin/edge-admin\n" > ${ROOT_DIR}/run.sh; \
    chmod u+x ${ROOT_DIR}/run.sh

EXPOSE 7788
EXPOSE 8001
EXPOSE 3306

ENTRYPOINT [ "/usr/local/goedge/run.sh" ]