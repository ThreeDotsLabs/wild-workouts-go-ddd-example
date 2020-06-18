# Wild Workouts

Wild Workouts is an example project that we created to show how to build Go applications that are **easy to develop, maintain, and fun to work with, especially in the long term!**

No application is perfect from the beginning. With over a dozen coming articles, we will uncover what issues you can find in the current implementation. We will also show how to fix these issues and achieve clean implementation by refactoring.

### Articles

1. [**Too modern Go application? Building a serverless application with Google Cloud Run and Firebase**](https://threedots.tech/post/serverless-cloud-run-firebase-modern-go-application/?utm_source=github.com)
2. [**A complete Terraform setup of a serverless application on Google Cloud Run and Firebase**](https://threedots.tech/post/complete-setup-of-serverless-application/?utm_source=github.com)
3. [**Robust gRPC communication on Google Cloud Run (but not only!)**](https://threedots.tech/post/robust-grpc-google-cloud-run/?utm_source=github.com)
4. [**You should not build your own authentication. Let Firebase do it for you.**](https://threedots.tech/post/firebase-cloud-run-authentication/?utm_source=github.com)
5. [**Business Applications in Go: Things to know about DRY**](https://threedots.tech/post/things-to-know-about-dry/?utm_source=github.com)
6. *More articles are on the way!*

### Directories

- [api](api/) OpenAPI and gRPC definitions
- [docker](docker/) Dockerfiles
- [internal](internal/) application code
- [scripts](scripts/) deployment and development scripts
- [terraform](terraform/) - infrastructure defintion
- [web](web/) - frontend JavaScript code

### Live Demo

The example application is available at [https://threedotslabs-wildworkouts.web.app/](https://threedotslabs-wildworkouts.web.app/).

### Running locally

```go
> docker-compose up

# ...

web_1             |  INFO  Starting development server...
web_1             |  DONE  Compiled successfully in 6315ms11:18:26 AM
web_1             |
web_1             |
web_1             |   App running at:
web_1             |   - Local:   http://localhost:8080/
web_1             |
web_1             |   It seems you are running Vue CLI inside a container.
web_1             |   Access the dev server via http://localhost:<your container's external mapped port>/
web_1             |
web_1             |   Note that the development build is not optimized.
web_1             |   To create a production build, run yarn build.
```

### Google Cloud Deployment

```go
> cd terraform/
> make

Fill all required parameters:
	project [current: wild-workouts project]:       # <----- put your Wild Workouts Google Cloud project name here (it will be created) 
	user [current: email@gmail.com]:                # <----- put your Google (Gmail, G-suite etc.) e-mail here
	billing_account [current: My billing account]:  # <----- your billing account name, can be found here https://console.cloud.google.com/billing
	region [current: europe-west1]: 
	firebase_location [current: europe-west]: 

# it may take a couple of minutes...

The setup is almost done!

Now you need to enable Email/Password provider in the Firebase console.
To do this, visit https://console.firebase.google.com/u/0/project/[your-project]/authentication/providers

You can also downgrade the subscription plan to Spark (it's set to Blaze by default).
The Spark plan is completely free and has all features needed for running this project.

Congratulations! Your project should be available at: https://[your-project].web.app

If it's not, check if the build finished successfully: https://console.cloud.google.com/cloud-build/builds?project=[your-project]

If you need help, feel free to contact us at https://threedots.tech
```

### Screenshots

![Wild Workouts login](https://threedots.tech/media/serverless-cloud-run-firebase-modern-go-app/login.png "Logo Title Text 1")
![Wild Workouts trainer's schedule](https://threedots.tech/media/serverless-cloud-run-firebase-modern-go-app/schedule.png "Logo Title Text 1")
![Wild Workouts schedule training](https://threedots.tech/media/serverless-cloud-run-firebase-modern-go-app/new-training.png "Logo Title Text 1")
