<template>
    <app-layout>
        <div class="py-5 text-center">
            <h2>Your trainings</h2>
            <p class="lead">Below is an example form built entirely with Bootstrap’s form controls. Each required form
                group
                has a validation state that can be triggered by attempting to submit the form without completing it.</p>
        </div>
        <br><br>
        <table class="table">
            <thead>
            <tr>
                <th scope="col">#</th>
                <th scope="col">When</th>
                <th scope="col">Notes</th>
                <th scope="col" v-if="isTrainer">Attendee</th>
                <th scope="col">Actions</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="(training, idx) in calendar" :key="training.uuid"
                v-bind:class="{'table-info': training.requireRescheduleApproval}">
                <th scope="row">{{ idx+1 }}</th>
                <td>
                    <span v-bind:class="{'old-date': training.proposedTime}">{{ formatDateTime(training.time) }}</span>
                    <span v-if="training.proposedTime" v-bind:title="'proposed by ' + training.moveProposedBy"><br>{{ formatDateTime(training.proposedTime) }}</span>
                </td>
                <th>{{ training.notes }}</th>
                <th v-if="isTrainer">{{ training.user }}</th>
                <td>
                    <button type="button"
                            v-bind:class="training.canBeCancelled ? 'btn btn-warning' : 'btn btn-danger'"
                            v-bind:title="training.canBeCancelled ? 'Your training balance will be returned' : 'Your training balance will be not returned because it\s less than 24h before training'"
                            @click="cancelTraining"
                            v-bind:data-training-uuid="training.uuid"
                    >
                        Cancel
                    </button>
                    &nbsp;
                    <router-link tag="button" class="btn btn-info"
                                 :to="{ name: 'proposeNewDate', params: { trainingID: training.uuid }}"
                                 v-if="training.moveRequiresAccept"
                    >
                        Propose new time
                    </router-link>
                    &nbsp;
                    <router-link tag="button" class="btn btn-primary" v-if="!training.moveRequiresAccept"
                                 :to="{ name: 'rescheduleTraining', params: { trainingID: training.uuid }}"
                    >
                        Move
                    </router-link>

                    <div v-if="training.proposedTime">
                        <br>
                        <button type="button" class="btn btn-warning" @click="acceptReschedule"
                                v-bind:data-training-uuid="training.uuid"
                                v-if="userRole !== training.moveProposedBy"
                        >
                            Approve reschedule
                        </button>
                        &nbsp;
                        <button type="button" class="btn btn-warning" @click="rejectReschedule"
                                v-bind:data-training-uuid="training.uuid"
                        >
                            <span v-if="userRole !== training.moveProposedBy">Reject reschedule</span>
                            <span v-if="userRole === training.moveProposedBy">Cancel reschedule request</span>
                        </button>
                    </div>
                </td>
            </tr>
            </tbody>
        </table>
    </app-layout>
</template>

<script>
    import AppLayout from '../layouts/App.vue'

    import {approveReschedule, cancelTraining, getCalendar, rejectReschedule} from '../repositories/trainings'
    import {getUserRole, Trainer} from "../repositories/user";
    import {formatDateTime} from "../date";

    export default {
        components: {
            AppLayout,
        },
        data: function () {
            return {
                'calendar': null,
                'isTrainer': null,
                'userRole': null,
            }
        },
        mounted() {
            let self = this
            getCalendar(function (calendar) {
                self.calendar = calendar
            })
            this.isTrainer = getUserRole() === Trainer;
            this.userRole = getUserRole()
        },
        methods: {
            cancelTraining(event) {
                let self = this

                let trainingUUID = event.target.getAttribute('data-training-uuid');
                let training = self.calendar.find(t => t.uuid === trainingUUID);

                let msg = 'Are you sure you want to cancel training?';


                let opts = {
                    title: msg,
                    html: true,
                    loader: true,
                }

                if (!training.canBeCancelled) {
                    opts.body = "<b>It's less than 24h before training, so you will not receive your credits back.</b>"
                } else {
                    opts.body = "Your training balance will be returned."
                }

                this.$dialog.confirm(opts)
                    .then(dialog => {
                        cancelTraining(trainingUUID, function () {
                            getCalendar(function (calendar) {
                                self.calendar = calendar
                            })
                            self.$toast.info('Training cancelled');
                            dialog.close()
                        }, function () {
                            self.$toast.error('Failed to cancel training');
                            dialog.close()
                        })
                    })
                    .catch(function () {
                        console.log('Clicked on cancel')
                    })
            },
            acceptReschedule(event) {
                let self = this;
                let trainingUUID = event.target.getAttribute('data-training-uuid');

                this.$dialog.confirm("Are you sure you want to accept?")
                    .then(dialog => {
                        approveReschedule(trainingUUID, function () {
                            getCalendar(function (calendar) {
                                self.calendar = calendar
                            })
                            self.$toast.info('Reschedule accepted');
                            dialog.close()
                        }, function () {
                            self.$toast.error('Failed to accept reschedule');
                            dialog.close()
                        })
                    })
                    .catch(function () {
                        console.log('Clicked on cancel')
                    })
            },
            rejectReschedule(event) {
                let self = this;
                let trainingUUID = event.target.getAttribute('data-training-uuid');

                this.$dialog.confirm("Are you sure you want to reject?")
                    .then(dialog => {
                        rejectReschedule(trainingUUID, function () {
                            getCalendar(function (calendar) {
                                self.calendar = calendar
                            })

                            self.$toast.info('Reschedule rejected');
                            dialog.close()
                        }, function () {
                            self.$toast.error('Failed to reject reschedule');
                            dialog.close()
                        })
                    })
                    .catch(function () {
                        console.log('Clicked on cancel')
                    })
            },
            formatDateTime,
        },
    }
</script>

<style scoped>
    h3 {
        margin: 40px 0 0;
    }

    ul {
        list-style-type: none;
        padding: 0;
    }

    li {
        display: inline-block;
        margin: 0 10px;
    }

    a {
        color: #42b983;
    }

    body {
        background-color: #f5f5f5;
    }

    .old-date {
        text-decoration: line-through;
        opacity: 0.5;
    }
</style>
