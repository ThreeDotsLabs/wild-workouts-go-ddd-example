<template>
    <app-layout>
        <div class="ww-page-head">
            <div>
                <h1>Your trainings</h1>
                <p class="ww-page-sub" v-if="!isTrainer">Upcoming sessions with your trainer</p>
                <p class="ww-page-sub" v-if="isTrainer">Upcoming sessions with your attendees</p>
            </div>
            <div class="ww-trainings-head-actions" v-if="!isTrainer">
                <div class="ww-card ww-balance" v-if="trainingBalance !== null">
                    <span class="ww-display ww-balance__value">{{ trainingBalance }}</span>
                    <span class="ww-balance__label">Trainings<br>left</span>
                </div>
                <router-link tag="button" class="ww-btn ww-btn--primary ww-btn--lg"
                             :to="{ name: 'scheduleTraining' }">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                        <path d="M8 3v10M3 8h10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    </svg>
                    Schedule training
                </router-link>
            </div>
        </div>

        <div class="ww-trainings">
            <div class="ww-card ww-training" v-for="training in calendar" :key="training.uuid"
                 v-bind:class="{'ww-training--pending': training.proposedTime}">
                <div class="ww-training__date">
                    <span class="ww-training__weekday">{{ weekdayShort(training.time) }}</span>
                    <span class="ww-display ww-training__day">{{ dayOfMonth(training.time) }}</span>
                    <span class="ww-training__month">{{ monthShort(training.time) }}</span>
                </div>
                <div class="ww-training__divider"></div>
                <div class="ww-training__body">
                    <div class="ww-training__time-row">
                        <template v-if="!training.proposedTime">
                            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                                <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.6"/>
                                <path d="M8 5v3l2 1.5" stroke="currentColor" stroke-width="1.6"
                                      stroke-linecap="round"/>
                            </svg>
                            <span class="ww-training__time">{{ formatHourRange(training.time) }}</span>
                        </template>
                        <template v-if="training.proposedTime">
                            <span class="ww-training__time-old">{{ formatHourRange(training.time) }}</span>
                            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                                <path d="M2 8h11M9 4l4 4-4 4" stroke="#B97E10" stroke-width="1.8"
                                      stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                            <span class="ww-training__time">{{ formatDayLong(training.proposedTime) }},
                                {{ formatHourRange(training.proposedTime) }}</span>
                            <span class="ww-badge ww-badge--amber">
                                Reschedule proposed by {{ training.moveProposedBy }}</span>
                        </template>
                    </div>
                    <span class="ww-training__notes" v-if="training.notes">{{ training.notes }}</span>
                    <span class="ww-training__attendee" v-if="isTrainer">with {{ training.user }}</span>
                </div>
                <div class="ww-training__actions">
                    <router-link tag="button" class="ww-btn ww-btn--ghost" v-if="!training.moveRequiresAccept"
                                 :to="{ name: 'rescheduleTraining', params: { trainingID: training.uuid }}">
                        Move
                    </router-link>
                    <router-link tag="button" class="ww-btn ww-btn--ghost"
                                 :to="{ name: 'proposeNewDate', params: { trainingID: training.uuid }}"
                                 v-if="training.moveRequiresAccept">
                        Propose new time
                    </router-link>
                    <button type="button" class="ww-btn ww-btn--danger"
                            v-bind:title="training.canBeCancelled ? 'Your training balance will be returned' : 'Your training balance will not be returned because it\'s less than 24h before the training'"
                            @click="cancelTraining"
                            v-bind:data-training-uuid="training.uuid">
                        Cancel
                    </button>
                    <template v-if="training.proposedTime">
                        <button type="button" class="ww-btn ww-btn--success" @click="acceptReschedule"
                                v-bind:data-training-uuid="training.uuid"
                                v-if="userRole !== training.moveProposedBy">
                            Approve
                        </button>
                        <button type="button" class="ww-btn ww-btn--ghost" @click="rejectReschedule"
                                v-bind:data-training-uuid="training.uuid">
                            <span v-if="userRole !== training.moveProposedBy">Reject</span>
                            <span v-if="userRole === training.moveProposedBy">Cancel request</span>
                        </button>
                    </template>
                </div>
            </div>

            <div class="ww-card ww-trainings-empty" v-if="calendar !== null && calendar.length === 0">
                <svg class="ww-trainings-empty__art" width="150" height="32" viewBox="0 0 150 32" fill="none"
                     aria-hidden="true">
                    <rect x="6" y="11.5" width="20" height="9" rx="4.5" fill="#191713"/>
                    <rect x="124" y="11.5" width="20" height="9" rx="4.5" fill="#191713"/>
                    <rect x="24" y="8" width="9" height="16" rx="3" fill="#191713"/>
                    <rect x="117" y="8" width="9" height="16" rx="3" fill="#191713"/>
                    <rect x="30" y="13" width="90" height="6" rx="3" fill="#191713"/>
                </svg>
                <span class="ww-display ww-trainings-empty__title">Nothing on the bar yet</span>
                <p v-if="!isTrainer">Your trainer is waiting. Grab an hour and get to work.</p>
                <p v-if="isTrainer">No sessions booked so far. Open some hours so attendees can book you.</p>
                <router-link tag="button" class="ww-btn ww-btn--primary ww-btn--lg" v-if="!isTrainer"
                             :to="{ name: 'scheduleTraining' }">
                    Schedule your first training
                </router-link>
                <router-link tag="button" class="ww-btn ww-btn--ghost ww-btn--lg" v-if="isTrainer"
                             :to="{ name: 'setSchedule' }">
                    Set your availability
                </router-link>
            </div>
        </div>
    </app-layout>
</template>

<script>
    import AppLayout from '../layouts/App.vue'

    import {approveReschedule, cancelTraining, getCalendar, rejectReschedule} from '../repositories/trainings'
    import {getTrainingBalance, getUserRole, Trainer} from "../repositories/user";
    import {apiErrorMessage} from "../repositories/errors";
    import {dayOfMonth, formatDayLong, formatHourRange, monthShort, weekdayShort} from "../date";

    export default {
        components: {
            AppLayout,
        },
        data: function () {
            return {
                'calendar': null,
                'isTrainer': null,
                'userRole': null,
                'trainingBalance': null,
            }
        },
        mounted() {
            let self = this
            this.refreshTrainings()
            this.isTrainer = getUserRole() === Trainer;
            this.userRole = getUserRole()

            if (!this.isTrainer) {
                getTrainingBalance(
                    balance => self.trainingBalance = balance,
                    () => self.$toast.error('Failed to load training balance'),
                );
            }
        },
        methods: {
            refreshTrainings() {
                let self = this
                getCalendar(function (calendar) {
                    self.calendar = calendar
                }, function () {
                    self.$toast.error('Failed to load trainings')
                })
            },
            cancelTraining(event) {
                let self = this

                let trainingUUID = event.currentTarget.getAttribute('data-training-uuid');
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
                            self.refreshTrainings()
                            self.$toast.info('Training cancelled');
                            dialog.close()
                        }, function (err) {
                            self.$toast.error(apiErrorMessage(err, 'Failed to cancel training'));
                            dialog.close()
                        })
                    })
                    .catch(function () {
                        console.log('Clicked on cancel')
                    })
            },
            acceptReschedule(event) {
                let self = this;
                let trainingUUID = event.currentTarget.getAttribute('data-training-uuid');

                this.$dialog.confirm({title: 'Are you sure you want to accept?'})
                    .then(dialog => {
                        approveReschedule(trainingUUID, function () {
                            self.refreshTrainings()
                            self.$toast.info('Reschedule accepted');
                            dialog.close()
                        }, function (err) {
                            self.$toast.error(apiErrorMessage(err, 'Failed to accept reschedule'));
                            dialog.close()
                        })
                    })
                    .catch(function () {
                        console.log('Clicked on cancel')
                    })
            },
            rejectReschedule(event) {
                let self = this;
                let trainingUUID = event.currentTarget.getAttribute('data-training-uuid');

                this.$dialog.confirm({title: 'Are you sure you want to reject?'})
                    .then(dialog => {
                        rejectReschedule(trainingUUID, function () {
                            self.refreshTrainings()

                            self.$toast.info('Reschedule rejected');
                            dialog.close()
                        }, function (err) {
                            self.$toast.error(apiErrorMessage(err, 'Failed to reject reschedule'));
                            dialog.close()
                        })
                    })
                    .catch(function () {
                        console.log('Clicked on cancel')
                    })
            },
            weekdayShort,
            monthShort,
            dayOfMonth,
            formatHourRange,
            formatDayLong,
        },
    }
</script>

<style scoped>
    .ww-trainings-head-actions {
        display: flex;
        align-items: center;
        gap: 14px;
    }

    .ww-balance {
        display: flex;
        align-items: center;
        gap: 12px;
        padding: 8px 18px;
        border-radius: 12px;
    }

    .ww-balance__value {
        font-size: 30px;
        font-weight: 700;
        color: var(--ww-accent);
        line-height: 1;
    }

    .ww-balance__label {
        font-size: 12px;
        color: var(--ww-muted);
        text-transform: uppercase;
        letter-spacing: 0.08em;
        line-height: 1.3;
    }

    .ww-trainings {
        display: flex;
        flex-direction: column;
        gap: 14px;
    }

    .ww-training {
        display: flex;
        align-items: center;
        gap: 22px;
        padding: 18px 24px;
    }

    .ww-training--pending {
        border-color: var(--ww-amber-line);
    }

    .ww-training__date {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 64px;
        flex-shrink: 0;
    }

    .ww-training__weekday, .ww-training__month {
        font-size: 11px;
        font-weight: 600;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: var(--ww-muted);
    }

    .ww-training__day {
        font-size: 34px;
        font-weight: 700;
        line-height: 1.05;
    }

    .ww-training__divider {
        width: 1px;
        height: 56px;
        background: var(--ww-line);
        flex-shrink: 0;
    }

    .ww-training__body {
        display: flex;
        flex-direction: column;
        gap: 6px;
        flex-grow: 1;
        min-width: 0;
    }

    .ww-training__time-row {
        display: flex;
        align-items: center;
        gap: 10px;
        flex-wrap: wrap;
    }

    .ww-training__time {
        font-size: 16px;
        font-weight: 600;
    }

    .ww-training__time-old {
        font-size: 15px;
        font-weight: 500;
        color: #A29C90;
        text-decoration: line-through;
    }

    .ww-training__notes {
        font-size: 14px;
        color: var(--ww-muted);
    }

    .ww-training__attendee {
        font-size: 14px;
        font-weight: 600;
    }

    .ww-training__actions {
        display: flex;
        gap: 8px;
        flex-wrap: wrap;
        justify-content: flex-end;
    }

    .ww-trainings-empty {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 6px;
        padding: 56px 48px 60px 48px;
        border-style: dashed;
        text-align: center;
        color: var(--ww-muted);
    }

    .ww-trainings-empty__art {
        margin-bottom: 14px;
    }

    .ww-trainings-empty__title {
        font-size: 30px;
        font-weight: 700;
        line-height: 1;
        color: var(--ww-ink);
    }

    .ww-trainings-empty p {
        max-width: 400px;
        font-size: 14px;
    }

    .ww-trainings-empty .ww-btn {
        margin-top: 16px;
    }
</style>
