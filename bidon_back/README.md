# DB sync guide

1. **dip provision services**
2. **CURL debezium connector**

   curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" 127.0.0.1:8083/connectors/ --data "@../debezium.json"

3. **Start karafka**

   dip bundle exec karafka server

4. **Start sync**

   dip rails appodeal:sync_test_apps

