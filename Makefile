profile ?= $(shell bash -c 'read -p "Profile: " profile; echo $$profile')

countdown:
	npm run --prefix internal/javascript build
	
watch:
	npm run --prefix internal/javascript watch

run:  
	aws-vault exec --no-session ${profile} -- go run cmd/poker/*.go

apply:
	cd terraform && make apply profile=${profile}
	# aws-vault exec --no-session ${profile} -- terraform -chdir=terraform apply