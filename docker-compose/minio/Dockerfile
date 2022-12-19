FROM minio/mc:latest

RUN mkdir /init
RUN mkdir /policy

COPY init.sh ./init/
COPY app-user-policy.json ./policy/

ENTRYPOINT ["./init/init.sh"]