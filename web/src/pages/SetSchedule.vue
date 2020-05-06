<template>
    <app-layout>
        <div class="py-5 text-center">
            <p style="font-size: 59px;">&#128170;</p>
            <h2>
                <span>Set schedule</span>
            </h2>

            <br>
            <p class="lead">Below is an example form built entirely with Bootstrap’s form controls. Each required form
                group
                has a validation state that can be triggered by attempting to submit the form without completing it.</p>

            <div class="btn-group btn-group-toggle">
                <label class="btn btn-outline-dark" v-for="period in periods" :key="period.from+period.to"
                       :class="{ 'active': selectedPeriod === selectedPeriodValue(period) }">
                    <input type="radio" name="options" v-bind:data-from="period.from" v-bind:data-to="period.to"
                           @change="changedPeriod" v-model="selectedPeriod" v-bind:value="selectedPeriodValue(period)">
                    {{ period.from }} - {{ period.to }}
                </label>
            </div>
        </div>

        <form class="needs-validation" novalidate>
            <div class="row">
                <div class="text-center schedule-column" v-for="day in schedule" :key="day.date.toString()">
                    <h4 class="mb-3">{{ formatDate(day.date) }}</h4>
                    <button type="button" class="btn btn-outline-primary btn-sm" v-bind:data-date="day.date"
                            @click="selectAllInDay">Select all
                    </button>
                    <br><br>

                    <div v-for="hour in day.hours" :key="formatDateTime(hour.hour)">
                        <div class="btn-group-toggle" data-toggle="buttons"
                             v-bind:title="hour.hasTrainingScheduled ? 'Training scheduled on this date' : ''">
                            <label v-bind:class="{'btn btn-lg': true, 'active': hour.available, 'btn-primary': !hour.hasTrainingScheduled, 'btn-secondary': hour.hasTrainingScheduled}">
                                <input type="checkbox" autocomplete="off" v-model="hour.available"
                                       @change.prevent="toggleHour($event, hour)"
                                       v-bind:data-hour="formatHour(hour.hour)" v-bind:data-date="formatDate(day.date)"
                                       v-bind:disabled="hour.hasTrainingScheduled">
                                {{ formatHour(hour.hour) }}
                            </label>
                        </div>
                        <br>
                    </div>
                </div>
            </div>
        </form>

    </app-layout>
</template>

<script>
    import AppLayout from '../layouts/App.vue'
    import {getPeriods, getSchedule, setHourAvailability} from "../repositories/trainings";
    import {formatDate, formatDateTime, formatHour} from "../date";

    export default {
        name: "SetSchedule",
        params: [],
        components: {
            AppLayout,
        },
        data: function () {
            return {
                'schedule': [],
                'periods': [],
                'selectedDateFrom': '',
                'selectedDateTo': '',
                'selectedPeriod': '',
            }
        },
        created() {
            let self = this

            this.periods = getPeriods()
            this.selectedDateFrom = this.periods[0].from
            this.selectedDateTo = this.periods[0].to
            this.selectedPeriod = this.selectedPeriodValue(this.periods[0])

            getSchedule(this.selectedDateFrom, this.selectedDateTo, function (data) {
                self.schedule = data
            })
        },
        methods: {
            selectAllInDay(event) {
                let date = event.target.getAttribute('data-date')
                let schedule = this.schedule;
                let self = this

                for (let scheduleIdx in schedule) {
                    let day = schedule[scheduleIdx]

                    if (date != day.date) {
                        continue
                    }

                    let updates = []

                    for (let idx in day.hours) {
                        let d = day.hours[idx].hour;
                        updates.push([formatDate(d), formatHour(d)])
                    }

                    setHourAvailability(updates, true, function () {
                        getSchedule(self.selectedDateFrom, self.selectedDateTo, function (data) {
                            self.schedule = data
                        })
                    })
                }

                this.schedule = schedule;
            },
            changedPeriod(event) {
                this.selectedDateFrom = event.target.getAttribute('data-from')
                this.selectedDateTo = event.target.getAttribute('data-to')

                let self = this
                getSchedule(this.selectedDateFrom, this.selectedDateTo, function (data) {
                    self.schedule = data
                })
            },
            toggleHour(event, hour) {
                let self = this

                let updates = [[event.target.getAttribute('data-date'), event.target.getAttribute('data-hour')]];
                setHourAvailability(
                    updates,
                    event.target.checked,
                    () => {
                    },
                    () => {
                        hour.available = false
                        self.$toast.error("Failed to update schedule")
                    },
                )
            },
            selectedPeriodValue(period) {
                return period.from + '/' + period.to
            },
            formatDate,
            formatDateTime,
            formatHour,
        },
    }
</script>

<style scoped>
    .btn-primary:not(.checked) {
        color: #007bff;
        background-color: transparent;
        background-image: none;
        border-color: #007bff;
    }

    .schedule-column {
        width: 140px;
    }
</style>