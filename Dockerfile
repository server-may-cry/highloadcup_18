FROM busybox

ADD entrypoint.sh .
ADD server .

ENV GOGC=1

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/server"]
