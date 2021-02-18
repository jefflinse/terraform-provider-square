.PHONY: default
default:

.PHONY: clean
clean:
	rm -rf .terraform/
	rm -f .terraform.lock.hcl
	rm -f terraform.tfstate

.PHONY: init
init: clean
	terraform init
