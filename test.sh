echo test new_proxy
curl -X POST http://localhost:8000/new_proxy \
     -H "Content-Type: application/json" \
     -d '{
           "content": {
             "user": {
               "user": "testUser",
               "metas": {"key1": "value1", "key2": "value2"},
               "run_id": "testRunID"
             },
             "proxy_name": "TestProxy",
             "proxy_type": "http",
             "use_encryption": true,
             "use_compression": false,
             "bandwidth_limit": "100M",
             "bandwidth_limit_mode": "per_user",
             "group": "group1",
             "group_key": "groupKey123",
             "custom_domains": ["example.com"],
             "subdomain": "sub.example.com",
             "locations": "/location1",
             "http_user": "httpUser",
             "http_pwd": "httpPassword",
             "host_header_rewrite": "example.com",
             "headers": {"header1": "value1", "header2": "value2"}
           }
         }'

echo \n
echo test new_work_conn
curl -X POST http://localhost:8000/new_work_conn \
     -H "Content-Type: application/json" \
     -d '{
           "content": {
             "user": {
               "user": "testUser",
               "metas": {"key1": "value1", "key2": "value2"},
               "run_id": "testRunID"
             },
             "run_id": "runID123",
             "timestamp": 1617183584,
             "privilege_key": "secretKey123"
           }
         }'

echo \n
echo test new_user_conn
curl -X POST http://localhost:8000/new_user_conn \
     -H "Content-Type: application/json" \
     -d '{
            "content": {
              "user": {
                "user": "testUser",
                "metas": {"key1": "value1", "key2": "value2"},
                "run_id": "testRunID"
              },
              "proxy_name": "TestProxy",
              "proxy_type": "tcp",
              "remote_addr": "192.168.0.106:1234"
            }
         }'

echo \n
echo test access
curl -X GET http://localhost:8000/access
