ARG RUBY_VERSION=3.1.2

#Copy geodb
FROM amazon/aws-cli:2.5.4 AS geodb
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY

RUN aws configure set aws_access_key_id ${AWS_ACCESS_KEY_ID} \
    && aws configure set aws_secret_access_key ${AWS_SECRET_ACCESS_KEY}

RUN aws s3 cp s3://appodeal-geo-db/prod/GeoIP2-City_current.mmdb /opt/GeoIP/ \
    && aws s3 cp s3://appodeal-geo-db/prod/SxGeoCity.dat /opt/GeoIP/

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

# Configure bundler
ENV LANG=C.UTF-8 \
  BUNDLE_JOBS=4 \
  BUNDLE_RETRY=3

# Store Bundler settings in the project's root
ENV BUNDLE_APP_CONFIG=.bundle

# Uncomment this line if you want to run binstubs without prefixing with `bin/` or `bundle exec`
ENV PATH /app/bin:$PATH

# Upgrade RubyGems and install the latest Bundler version
RUN gem update --system && gem install bundler

# Create a directory for the app code
RUN mkdir -p /app
WORKDIR /app

ENV RUBYOPT="--enable-frozen_string_literal,yjit"

# Document that we're going to expose port 3000
EXPOSE 3000
# Use Bash as the default command
CMD ["/bin/bash"]

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

ENV LANG=C.UTF-8 \
  BUNDLE_JOBS=4 \
  BUNDLE_RETRY=3

# Change user to deploy
USER deploy

WORKDIR /app

COPY --chown=deploy:deploy Gemfile .
COPY --chown=deploy:deploy Gemfile.lock .

RUN bundle install && rm -rf /data/bundle/cache
ENV RUBYOPT="--enable-frozen_string_literal,yjit"
ENV LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so
# Check that Ruby uses jemalloc
# MALLOC_CONF=stats_print:true ruby -e "exit"

COPY --chown=deploy:deploy . .

COPY --from=geodb /opt/GeoIP/* /opt/apps/appodeal/public/system/
