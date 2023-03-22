ARG RUBY_VERSION=3.1.2

FROM ruby:$RUBY_VERSION-slim AS dev

# Common dependencies
RUN apt-get update -qq \
  && DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends \
    build-essential \
    gnupg2 \
    curl \
    less \
    git \
    vim \
    libpq-dev \
  && apt-get clean \
  && rm -rf /var/cache/apt/archives/* \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
  && truncate -s 0 /var/log/*log

#ARG NODE_MAJOR
#RUN curl -sL https://deb.nodesource.com/setup_$NODE_MAJOR.x | bash -
#RUN apt-get update -qq && DEBIAN_FRONTEND=noninteractive apt-get -yq dist-upgrade && \
#  DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends \
#    nodejs \
#    && apt-get clean \
#    && rm -rf /var/cache/apt/archives/* \
#    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
#    && truncate -s 0 /var/log/*log
#
#ARG YARN_VERSION
#RUN npm install -g yarn@$YARN_VERSION

# Configure bundler
ENV LANG=C.UTF-8 \
  BUNDLE_JOBS=4 \
  BUNDLE_RETRY=3

# Store Bundler settings in the project's root
ENV BUNDLE_APP_CONFIG=.bundle

# Uncomment this line if you want to run binstubs without prefixing with `bin/` or `bundle exec`
ENV PATH=/app/bin:$PATH

# Upgrade RubyGems and install the latest Bundler version
RUN gem update --system && \
    gem install bundler

# Create a directory for the app code
RUN mkdir -p /app
WORKDIR /app

ENV RUBYOPT="--enable-frozen_string_literal,yjit"

# Document that we're going to expose port 3000
EXPOSE 3000
# Use Bash as the default command
CMD ["/usr/bin/bash"]

FROM ruby:$RUBY_VERSION-slim AS prod

RUN apt-get update -qq \
  && DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends \
    build-essential \
    gnupg2 \
    libpq-dev \
    libjemalloc-dev \
    vim \
  && apt-get clean \
  && rm -rf /var/cache/apt/archives/* \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
  && truncate -s 0 /var/log/*log

RUN useradd -ms /usr/sbin/nologin -u 1000 deploy

# Configure bundler
ENV LANG=C.UTF-8 \
  BUNDLE_JOBS=4 \
  BUNDLE_RETRY=3

USER deploy

WORKDIR /app

COPY --chown=deploy:deploy Gemfile .
COPY --chown=deploy:deploy Gemfile.lock .

RUN bundle install && rm -rf /data/bundle/cache

COPY --chown=deploy:deploy . .

ENV PG_DATABASE= \
  PG_DATABASE= \
  PG_USERNAME= \
  PG_PASSWORD= \
  PG_HOST= \
  PG_PORT= \
  APPODEAL_PG_DATABASE= \
  APPODEAL_PG_USERNAME= \
  APPODEAL_PG_PASSWORD= \
  APPODEAL_PG_HOST= \
  APPODEAL_PG_PORT= \
  ALLOWED_HOSTS=

ENV RUBYOPT="--enable-frozen_string_literal,yjit"
ENV LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so
# Check that Ruby uses jemalloc
# MALLOC_CONF=stats_print:true ruby -e "exit"

# Assets
RUN bundle exec rails assets:precompile
