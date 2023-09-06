SHELL:=/bin/bash

apply:
	source .env && aws-vault exec --no-session ${profile} -- terraform apply