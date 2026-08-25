<template>
    <header class="ww-nav">
        <div class="ww-nav__brand">
            <svg width="44" height="24" viewBox="0 0 48 26" fill="none">
                <path d="M3 5.5 L8.74 23 L13.5 1.8 L18.26 23 L24 5.5" stroke="#E0491F" stroke-width="3.6"
                      stroke-linecap="round" stroke-linejoin="miter" stroke-miterlimit="8"/>
                <path d="M3 5.5 L8.74 23 L13.5 1.8 L18.26 23 L24 5.5" transform="translate(21 0)"
                      stroke="#F7F5F1" stroke-width="3.6"
                      stroke-linecap="round" stroke-linejoin="miter" stroke-miterlimit="8"/>
            </svg>
            <span class="ww-display ww-nav__title">Wild Workouts</span>
        </div>
        <nav class="ww-nav__links">
            <router-link :to="{ name: 'trainingsList' }" exact active-class="is-active">Trainings</router-link>
            <router-link :to="{ name: 'calendar' }" active-class="is-active">Calendar</router-link>
            <router-link :to="{ name: 'scheduleTraining' }" active-class="is-active"
                         v-if="userType === 'attendee'">Schedule training
            </router-link>
            <router-link :to="{ name: 'setSchedule' }" active-class="is-active"
                         v-if="userType === 'trainer'">Set schedule
            </router-link>
        </nav>
        <div class="ww-nav__user">
            <div class="ww-nav__avatar" v-bind:class="{'ww-nav__avatar--trainer': userType === 'trainer'}">
                {{ initials }}
            </div>
            <span class="ww-nav__role">{{ roleName }}</span>
            <a class="ww-nav__logout" v-on:click="signOut" href="/login">Log out</a>
        </div>
    </header>
</template>

<script>
    import {getUserRole} from '../repositories/user'
    import {Auth} from "../repositories/auth";

    export default {
        name: "WebMenu",
        methods: {
            signOut: function () {
                Auth.logout().finally(function () {
                    self.$router.push({name: 'login'});
                })
            }
        },
        computed: {
            roleName: function () {
                if (!this.userType) {
                    return ''
                }
                return this.userType.charAt(0).toUpperCase() + this.userType.slice(1)
            },
            initials: function () {
                return this.userType ? this.userType.slice(0, 2).toUpperCase() : ''
            },
        },
        data: function () {
            return {
                'userType': getUserRole(),
            }
        }
    }
</script>

<style scoped>
    .ww-nav {
        display: flex;
        align-items: center;
        gap: 24px;
        height: 64px;
        padding: 0 32px;
        background: var(--ww-ink);
        color: var(--ww-bg);
    }

    .ww-nav__brand {
        display: flex;
        align-items: center;
        gap: 10px;
    }

    .ww-nav__title {
        font-size: 20px;
        font-weight: 700;
        letter-spacing: 0.06em;
    }

    .ww-nav__links {
        display: flex;
        gap: 4px;
        margin-left: 24px;
        flex-grow: 1;
    }

    .ww-nav__links a {
        padding: 8px 14px;
        border-radius: 8px;
        font-size: 14px;
        font-weight: 500;
        color: #B9B4AA;
    }

    .ww-nav__links a:hover {
        color: #FFFFFF;
    }

    .ww-nav__links a.is-active {
        font-weight: 600;
        color: #FFFFFF;
        background: rgba(255, 255, 255, 0.10);
    }

    .ww-nav__user {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .ww-nav__avatar {
        width: 30px;
        height: 30px;
        border-radius: 999px;
        background: var(--ww-accent);
        color: #FFFFFF;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 12px;
        font-weight: 600;
    }

    .ww-nav__avatar--trainer {
        background: var(--ww-green);
    }

    .ww-nav__role {
        font-size: 13px;
        font-weight: 600;
        color: var(--ww-bg);
    }

    .ww-nav__logout {
        padding: 7px 14px;
        border-radius: 8px;
        border: 1px solid rgba(255, 255, 255, 0.25);
        color: var(--ww-bg);
        font-size: 13px;
        font-weight: 500;
    }

    .ww-nav__logout:hover {
        border-color: rgba(255, 255, 255, 0.6);
        color: #FFFFFF;
    }
</style>
