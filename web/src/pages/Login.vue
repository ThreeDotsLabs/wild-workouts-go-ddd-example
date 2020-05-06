<template>
    <login-layout>
        <form class="form-signin" v-on:submit.prevent="submit">
            <span style="font-size: 60px;">&#128023;</span>
            <h1 class="h3 mb-3 font-weight-normal">Please sign in</h1>
            <div class="alert alert-primary" role="alert">
                <p v-for="user in getTestUsers()" :key="user.login">
                    {{ user.role.charAt(0).toUpperCase() + user.role.slice(1) }} credentials
                    <a href="#" v-bind:data-login="user.login" v-bind:data-password="user.password"
                       v-on:click="loadCredentials" :key="user.login">{{ user.login }}:{{ user.password }}</a>
                </p>
            </div>

            <label for="inputEmail" class="sr-only">Email address</label>
            <input type="email" id="inputEmail" class="form-control" v-model="login" placeholder="Email address"
                   required autofocus>
            <label for="inputPassword" class="sr-only">Password</label>
            <input type="password" id="inputPassword" class="form-control" v-model="password" placeholder="Password"
                   required>
            <div class="checkbox mb-3">
                <label>
                    <input type="checkbox" value="remember-me"> Remember me
                </label>
            </div>
            <button class="btn btn-lg btn-primary btn-block" type="submit">
                Sign in
                <span v-if="showLoader" class="spinner-grow spinner-grow-sm" style="width: 1.3rem; height: 1.3rem;"
                      role="status" aria-hidden="true"></span>
            </button>
            <p class="mt-5 mb-3 text-muted">&copy; 2017-2020</p>
        </form>
    </login-layout>
</template>

<script>
    import LoginLayout from '../layouts/Login'
    import {getTestUsers, loginUser} from '../repositories/user'
    import {Auth} from "../repositories/auth";


    export default {
        name: "Login",
        components: {
            LoginLayout,
        },
        mounted() {
            if (Auth.isLoggedIn()) {
                this.$router.push({name: 'trainingsList'});
            }
        },
        methods: {
            submit: function () {
                let self = this
                this.showLoader = true

                loginUser(this.login, this.password)
                    .then(function () {
                        self.$toast.info("Hey buddy!")
                        self.$router.push({name: 'trainingsList'})
                    })
                    .catch(error => {
                        self.$toast.error("Failed to log in")
                        console.error(error)
                        self.showLoader = false
                    })
            },
            loadCredentials(event) {
                this.login = event.target.getAttribute('data-login');
                this.password = event.target.getAttribute('data-password');
            },
            getTestUsers,
        }
        ,
        data: function () {
            return {
                'login': '',
                'password': '',
                'showLoader': false,
            }
        }
    }
</script>

<style scoped>
    .form-signin {
        width: 100%;
        width: 330px;
        padding: 15px;
        margin: 0 auto;
    }

    .form-signin .checkbox {
        font-weight: 400;
    }

    .form-signin .form-control {
        position: relative;
        box-sizing: border-box;
        height: auto;
        padding: 10px;
        font-size: 16px;
    }

    .form-signin .form-control:focus {
        z-index: 2;
    }

    .form-signin input[type="email"] {
        margin-bottom: -1px;
        border-bottom-right-radius: 0;
        border-bottom-left-radius: 0;
    }

    .form-signin input[type="password"] {
        margin-bottom: 10px;
        border-top-left-radius: 0;
        border-top-right-radius: 0;
    }
</style>