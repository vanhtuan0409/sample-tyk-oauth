reload:
	curl -s -H "x-tyk-authorization: myabcsecret" localhost:8081/tyk/reload | jq

list-apis:
	curl -s -H "x-tyk-authorization: myabcsecret" localhost:8081/tyk/apis | jq

list-apis-oas:
	curl -s -H "x-tyk-authorization: myabcsecret" localhost:8081/tyk/apis/oas | jq

import-oas:
	curl -s -X POST \
		-H "x-tyk-authorization: myabcsecret" \
		-H "content-type: application/json" \
		-d "@./conf/tyk/newapps/myapp.json" \
		localhost:8081/tyk/apis/oas
