base_image:=gcr.io/cloudrun/hello
load_envs:=source ./.env

all: check init apply push_repo done

.PHONY: check
check:
	@command -v terraform >/dev/null || ( echo "terraform is not installed"; exit 1 )
	@command -v gcloud >/dev/null || ( echo "gcloud is not installed"; exit 1 )

.PHONY: init
init:
	@./set-envs.sh
	terraform init

.PHONY: apply
apply:
	${load_envs} && terraform apply

.PHONY: push_repo
push_repo:
	git config --global credential.https://source.developers.google.com.helper gcloud.sh  
	repo_url=$(shell terraform output -no-color | grep repo_url | awk '{print $$3}') && \
	if [ -z "$(shell git remote | grep google)" ]; then \
		git remote add google "$$repo_url"; \
	else \
		git remote set-url google "$$repo_url"; \
	fi
	git push google master

.PHONY: done
done:
	@echo
	@printf "=%.0s" {1..80}
	@echo
	@echo
	@printf "\e[93mThe setup is almost done!\e[0m\n"
	@echo
	@printf "Now you need to enable \e[1mEmail/Password\e[0m provider in the Firebase console.\n"
	@${load_envs} && echo "To do this, visit https://console.firebase.google.com/u/0/project/$$TF_VAR_project/authentication/providers"
	@echo
	@printf "You can also downgrade the subscription plan to \e[1mSpark\e[0m (it's set to Blaze by default).\n"
	@echo "The Spark plan is completely free and has all features needed for running this project."
	@echo
	@${load_envs} && printf "\e[92mCongratulations!\e[0m Your project should be available at: https://$$TF_VAR_project.web.app \n"
	@echo
	@${load_envs} && echo "If it's not, check if the build finished successfully: https://console.cloud.google.com/cloud-build/builds?project=$$TF_VAR_project"
	@echo
	@echo "If you need help, feel free to contact us at https://threedots.tech"

.PHONY: set_project
set_project:
	${load_envs} && gcloud config set project "$$TF_VAR_project"

.PHONY: firestore
firestore: set_project
	${load_envs} && gcloud alpha firestore databases create "--region=$$TF_VAR_firebase_location"

.PHONY: firebase_builder
firebase_builder: set_project
	[ -d cloud-builders-community ] || git clone https://github.com/GoogleCloudPlatform/cloud-builders-community.git
	cd cloud-builders-community/firebase; gcloud builds submit .

.PHONY: docker_images
docker_images:
	docker pull "${base_image}"

	${load_envs} && for service in trainer trainings users; do \
		tag="gcr.io/$$TF_VAR_project/$$service"; \
		docker tag "${base_image}" "$$tag"; \
		docker push "$$tag"; \
	done

.PHONY: destroy
destroy:
	terraform state rm "google_project_iam_member.owner"
	terraform state rm "google_project_service.container_registry"
	terraform state rm "google_project_service.cloud_run"
	${load_envs} && terraform destroy
	git remote rm google
