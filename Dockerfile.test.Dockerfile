FROM busybox
EXPOSE 80

ADD data /tmp/data
ADD entrypoint.sh .
ADD server .

ENV GOGC=1
ENV GODEBUG=gctrace=1

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/server"]
