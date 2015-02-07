#
# Simple Ubuntu docker file to run the deviation binary in a loop.
#
FROM       ubuntu:latest
MAINTAINER Mark Mandel <mark.mandel@gmail.com>

CMD ["/bin/bash", "-c", "while true; do /opt/bin/deviation; sleep 60; done"]
#CMD ["/opt/bin/deviation"]