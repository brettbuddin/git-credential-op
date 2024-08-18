FROM alpine:3.20

RUN apk add --no-cache \
  nginx \
  fcgiwrap \
  openssl \
  git-daemon \
  git \
  openrc \
  mdevd-openrc

# Install credentials: username=foo, password=bar
RUN printf "foo:$(echo -n "bar" | openssl passwd -apr1 -stdin)\n" >> /etc/nginx/.gitpasswd

# Set up the "test" repository.
RUN mkdir -p /srv/git/test \
	&& cd /srv/git/test \
	&& git init . --bare --shared \
	&& git update-server-info
COPY gitconfig /srv/git/test/config

# Install nginx configuration.
COPY nginx.conf /etc/nginx/http.d/git.conf

# Allow fcgiwrap to deal with the repository data.
RUN chown -R fcgiwrap:www-data /srv/git

# Start nginx and fcgiwrap on boot of the container
RUN sed -i '/getty/d' /etc/inittab \
	&& rc-update add nginx default \
	&& rc-update add fcgiwrap default

CMD ["/sbin/init"]
