Start rails server in production mode and call wrk
```shell
wrk -c10 -d 20 -t5 -s wrk.lua http://0.0.0.0:3100/config
```
