<template>
    <app-layout>
        <div class="py-5 text-center">
            <p style="font-size: 59px;">&#128170;</p>
            <h2>
                <span v-if="!trainingToReschedule">Schedule training</span>
                <span v-if="trainingToReschedule && !isPropose">Re-schedule training {{trainingToReschedule}}</span>
                <span v-if="trainingToReschedule && isPropose">Propose re-schedule training {{trainingToReschedule}}</span>
            </h2>

            <br>
            <p class="lead">Below is an example form built entirely with Bootstrap’s form controls. Each required form
                group
                has a validation state that can be triggered by attempting to submit the form without completing it.</p>

            <div class="alert alert-warning" role="alert" v-if="isPropose">
                It's less than 24h left until the training. Proposition of re-schedule may be refused.
            </div>

            <div v-if="!isPropose">
                <span class="trainings-left">Trainings left: <b>{{ trainingBalance }}</b></span>
            </div>
        </div>
        <div class="row justify-content-md-center">
            <div class="col-md-8 order-md-1 l-md">
                <form class="needs-validation" @submit.prevent="scheduleNewTraining" novalidate>
                    <div class="row">
                        <div class="col-md-6 mb-3">
                            <label for="day">Day</label>
                            <select class="custom-select" size="7" id="day" v-model="trainingData.date"
                                    v-on:change="trainingDayChange">
                                <option v-for="day in calendar" :key="formatDate(day.date)"
                                        :value="formatDate(day.date)">
                                    {{ formatDate(day.date) }}
                                </option>
                            </select>
                        </div>
                        <div class="col-md-6 mb-3">
                            <label for="hour">Hour</label>
                            <select class="custom-select" size="7" id="hour" v-model="trainingData.hour">
                                <template v-for="hour in availableHours">
                                    <option :key="formatHour(hour.hour)" :value="formatHour(hour.hour)"
                                            v-if="!hour.hasTrainingScheduled"> {{ formatHour(hour.hour) }}
                                    </option>
                                </template>
                            </select>
                        </div>
                    </div>


                    <div class="form-group">
                        <label for="notes">Notes <small>(visible for trainer)</small></label>
                        <textarea class="form-control" id="notes" rows="3" v-model="trainingData.notes" maxlength="1000"></textarea>
                    </div>

                    <hr class="mb-4">
                    <button class="btn btn-primary btn-lg btn-block" type="submit">
                        Schedule training
                        <span v-if="showLoader" class="spinner-grow spinner-grow-sm"
                              style="width: 1.3rem; height: 1.3rem;"
                              role="status" aria-hidden="true"></span>
                    </button>
                </form>
            </div>
        </div>
    </app-layout>
</template>

<script>
    import AppLayout from '../layouts/App.vue'
    import {getTrainingBalance} from "../repositories/user";
    import {getAvailableDates, rescheduleTraining, scheduleTraining} from "../repositories/trainings";
    import {formatDate, formatHour} from "../date";

    export default {
        name: "ScheduleTraining",
        params: [],
        components: {
            AppLayout,
        },
        created() {
            let self = this

            getAvailableDates(function (data) {
                self.calendar = data
            })
            this.trainingToReschedule = this.$route.params['trainingID'];
            this.isReschedule = this.$attrs.isReschedule;
            this.isPropose = this.$attrs.isPropose;
            getTrainingBalance(balance => self.trainingBalance = balance);
        },
        data: function () {
            return {
                'trainingData': {
                    'date': '',
                    'hour': '',
                    'notes': '',
                },
                'isReschedule': null,
                'isPropose': null,
                'trainingToReschedule': null,
                'calendar': [],
                'availableHours': [],
                'trainingBalance': null,
                'showLoader': false,
            }
        },
        methods: {
            trainingDayChange() {
                const currentDate = this.calendar.find(obj => formatDate(obj.date) === this.trainingData.date);

                if (!currentDate) {
                    return
                }

                this.availableHours = currentDate.hours.filter(obj => obj.available === true);
            },
            scheduleNewTraining() {
                let self = this;

                self.showLoader = true

                if (self.trainingToReschedule != null) {
                    rescheduleTraining(
                        this.trainingToReschedule,
                        this.trainingData.notes,
                        this.trainingData.date,
                        this.trainingData.hour,
                        this.isPropose,
                        function () {
                            if (self.isPropose) {
                                self.$toast.success('Training reschedule proposal sent!');
                            } else {
                                self.$toast.success('Training rescheduled!');
                            }
                            self.showLoader = false
                            self.$router.push({name: 'trainingsList'});
                        },
                        function (err) {
                            self.showLoader = false
                            self.$toast.error("Failed to reschedule training");
                            console.error(err)
                        },
                    )
                } else {
                    scheduleTraining(
                        this.trainingData.notes,
                        this.trainingData.date,
                        this.trainingData.hour,
                        function () {
                            self.showLoader = false
                            self.$toast.success('Training added!');
                            self.$router.push({name: 'trainingsList'});
                        },
                        function (err) {
                            self.showLoader = false
                            self.$toast.error("Failed to add training");
                            console.error(err)
                        }
                    )
                }
            },
            formatDate,
            formatHour,
        },
    }
</script>

<style scoped>
    .trainings-left {
        font-size: 1.5rem;
    }
</style>