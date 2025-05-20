FROM confluentinc/cp-kafka-connect:7.2.1

USER root

RUN confluent-hub install --no-prompt confluentinc/kafka-connect-s3:10.4.1

USER appuser
