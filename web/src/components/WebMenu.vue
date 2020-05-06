<template>
    <div class="d-flex flex-column flex-md-row align-items-center p-3 px-md-4 mb-3 bg-white border-bottom shadow-sm">
        <h5 class="my-0 mr-md-auto font-weight-normal">Wild Workouts &#128023;</h5>
        <nav class="my-2 my-md-0 mr-md-3">
            <router-link class="p-2 text-dark" :to="{ name: 'scheduleTraining' }" v-if="userType === 'attendee'">
                Schedule new training
            </router-link>
            <router-link class="p-2 text-dark" :to="{ name: 'trainingsList' }">Trainings</router-link>
            <router-link class="p-2 text-dark" :to="{ name: 'calendar' }">Calendar</router-link>
            <router-link class="p-2 text-dark" :to="{ name: 'setSchedule' }" v-if="userType === 'trainer'">Set
                schedule
            </router-link>
        </nav>
        <a class="btn btn-outline-primary" v-on:click="signOut" href="/login">Logout</a>
    </div>
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
        data: function () {
            return {
                'userType': getUserRole(),
            }
        }
    }
</script>

<style scoped>

</style>