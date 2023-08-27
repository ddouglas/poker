profile ?= $(shell bash -c 'read -p "Profile: " profile; echo $$profile')

run:
	aws-vault exec ${profile} -- go run cmd/poker/*.go

apply:
	aws-vault exec ${profile} -- terraform -chdir=terraform apply