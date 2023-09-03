profile ?= $(shell bash -c 'read -p "Profile: " profile; echo $$profile')

countdown:
	cd ./internal/javascript && npm run build

generate:
	templ generate

run: generate 
	aws-vault exec --no-session ${profile} -- go run cmd/poker/*.go

apply:
	aws-vault exec --no-session ${profile} -- terraform -chdir=terraform apply