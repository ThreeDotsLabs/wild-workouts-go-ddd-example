<template>
    <login-layout>
        <div class="ww-login">
            <div class="ww-login__brand">
                <div class="ww-login__logo">
                    <svg width="63" height="34" viewBox="0 0 48 26" fill="none">
                        <path d="M3 5.5 L8.74 23 L13.5 1.8 L18.26 23 L24 5.5" stroke="#E0491F" stroke-width="3.6"
                              stroke-linecap="round" stroke-linejoin="miter" stroke-miterlimit="8"/>
                        <path d="M3 5.5 L8.74 23 L13.5 1.8 L18.26 23 L24 5.5" transform="translate(21 0)"
                              stroke="#F7F5F1" stroke-width="3.6"
                              stroke-linecap="round" stroke-linejoin="miter" stroke-miterlimit="8"/>
                    </svg>
                    <span class="ww-display ww-login__wordmark">Wild Workouts</span>
                </div>
                <div>
                    <h1 class="ww-display ww-login__headline">Train wild.<br><span>Book in seconds.</span></h1>
                    <p class="ww-login__tagline">Personal training with your coach.<br>Pick an open hour, show up, and lift.</p>
                </div>
                <div class="ww-login__about">
                    <p>Wild Workouts is the open-source example application for
                        <a href="https://threedots.tech/go-with-the-domain/?utm_source=wild-workouts" target="_blank"
                           rel="noopener">Go with the Domain</a> &mdash;
                        Three Dots Labs&rsquo; book on DDD &amp; Clean Architecture in Go.</p>
                    <a class="ww-login__repo" href="https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example"
                       target="_blank" rel="noopener">
                        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
                            <path fill-rule="evenodd" clip-rule="evenodd"
                                  d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8Z"/>
                        </svg>
                        wild-workouts-go-ddd-example
                    </a>
                    <span class="ww-login__copyright">&copy; Wild Workouts 2017&ndash;2026</span>
                </div>
            </div>

            <div class="ww-login__panel">
                <form class="ww-login__form" v-on:submit.prevent="submit">
                    <div>
                        <h2 class="ww-display ww-login__title">Welcome back</h2>
                        <p class="ww-login__sub">Sign in to manage your trainings.</p>
                    </div>

                    <div class="ww-card ww-login__demo">
                        <span class="ww-step-label">Demo accounts &mdash; tap to fill</span>
                        <div class="ww-login__demo-chips">
                            <a href="#" class="ww-login__chip" v-for="user in getTestUsers()" :key="user.login"
                               v-bind:title="user.login + ':' + user.password"
                               v-bind:data-login="user.login" v-bind:data-password="user.password"
                               v-bind:class="{'is-selected': login === user.login}"
                               v-on:click.prevent="loadCredentials">
                                <span class="ww-login__chip-dot"
                                      v-bind:class="{'ww-login__chip-dot--attendee': user.role === 'attendee'}"></span>
                                {{ user.role.charAt(0).toUpperCase() + user.role.slice(1) }} &middot; {{ user.name }}
                            </a>
                        </div>
                    </div>

                    <div>
                        <label class="ww-label" for="inputEmail">Email</label>
                        <input type="email" id="inputEmail" class="ww-input" v-model="login"
                               placeholder="you@example.com" required autofocus>
                    </div>
                    <div>
                        <label class="ww-label" for="inputPassword">Password</label>
                        <input type="password" id="inputPassword" class="ww-input" v-model="password"
                               placeholder="Your password" required>
                    </div>
                    <label class="ww-login__remember">
                        <input type="checkbox" value="remember-me"> Remember me
                    </label>
                    <button class="ww-btn ww-btn--primary ww-btn--lg ww-btn--block" type="submit">
                        Sign in
                        <span v-if="showLoader" class="ww-spinner"></span>
                    </button>
                </form>
            </div>
        </div>
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
                let target = event.currentTarget;
                this.login = target.getAttribute('data-login');
                this.password = target.getAttribute('data-password');
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
    .ww-login {
        display: flex;
        min-height: 100vh;
    }

    .ww-login__brand {
        width: 40%;
        min-width: 380px;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        padding: 48px;
        background: var(--ww-ink);
        color: var(--ww-bg);
        background-image: repeating-linear-gradient(-55deg, rgba(255, 255, 255, 0.025) 0 2px, transparent 2px 26px);
    }

    .ww-login__logo {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .ww-login__wordmark {
        font-size: 22px;
        font-weight: 700;
        letter-spacing: 0.06em;
    }

    .ww-login__headline {
        font-size: 68px;
        font-weight: 700;
        line-height: 0.95;
    }

    .ww-login__headline span {
        color: var(--ww-accent);
    }

    .ww-login__tagline {
        margin-top: 18px;
        max-width: 340px;
        font-size: 16px;
        color: #A9A398;
    }

    .ww-login__about {
        display: flex;
        flex-direction: column;
        gap: 10px;
        font-size: 13px;
        color: #A9A398;
    }

    .ww-login__about p {
        max-width: 340px;
        line-height: 1.5;
    }

    .ww-login__repo {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        align-self: flex-start;
        padding: 7px 12px;
        border-radius: 8px;
        border: 1px solid rgba(255, 255, 255, 0.2);
        color: var(--ww-bg);
        font-size: 13px;
        font-weight: 600;
    }

    .ww-login__repo:hover {
        border-color: rgba(255, 255, 255, 0.55);
        color: #FFFFFF;
    }

    .ww-login__copyright {
        font-size: 13px;
        color: var(--ww-faint);
    }

    .ww-login__panel {
        flex-grow: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 48px 32px;
    }

    .ww-login__form {
        width: 100%;
        max-width: 400px;
        display: flex;
        flex-direction: column;
        gap: 20px;
    }

    .ww-login__title {
        font-size: 38px;
        font-weight: 700;
        line-height: 1;
    }

    .ww-login__sub {
        margin-top: 6px;
        font-size: 15px;
        color: var(--ww-muted);
    }

    .ww-login__demo {
        display: flex;
        flex-direction: column;
        gap: 10px;
        padding: 14px 16px;
        border-radius: 12px;
    }

    .ww-login__demo-chips {
        display: flex;
        gap: 8px;
        flex-wrap: wrap;
    }

    .ww-login__chip {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 14px;
        border-radius: 999px;
        border: 1px solid var(--ww-line);
        background: var(--ww-bg);
        font-size: 13px;
        font-weight: 600;
        color: var(--ww-ink);
    }

    .ww-login__chip:hover {
        border-color: var(--ww-ink);
        color: var(--ww-ink);
    }

    .ww-login__chip.is-selected {
        border-color: var(--ww-ink);
        background: var(--ww-ink);
        color: #FFFFFF;
    }

    .ww-login__chip-dot {
        width: 8px;
        height: 8px;
        border-radius: 999px;
        background: var(--ww-green);
    }

    .ww-login__chip-dot--attendee {
        background: var(--ww-accent);
    }

    .ww-login__remember {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 14px;
        color: var(--ww-muted);
    }

    @media (max-width: 900px) {
        .ww-login {
            flex-direction: column;
        }

        .ww-login__brand {
            width: 100%;
            min-width: 0;
            gap: 32px;
        }

        .ww-login__headline {
            font-size: 48px;
        }
    }
</style>
