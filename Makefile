restart_docker:
	docker-compose down -v && docker-compose up -d && docker-compose logs -f tyk

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

list-oauth-app:
	curl -s -H "x-tyk-authorization: myabcsecret" localhost:8081/tyk/oauth/clients/myapp_oauth | jq

list-oauth-token:
	curl -s -H "x-tyk-authorization: myabcsecret" localhost:8081/tyk/oauth/clients/myapp_oauth/sample_oauth_app/tokens | jq

start-flow:
	curl -v -s -X POST \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "response_type=code&client_id=sample_oauth_app&redirect_uri=http://localhost:5050/redirect" \
		http://localhost:8080/oauth/api/oauth/authorize | jq
