<template>
    <app-layout>
        <div class="ww-page-head">
            <div>
                <h1 v-if="!trainingToReschedule">Schedule a training</h1>
                <h1 v-if="trainingToReschedule && !isPropose">Move your training</h1>
                <h1 v-if="trainingToReschedule && isPropose">Propose a new time</h1>
                <p class="ww-page-sub">Pick any open hour in your trainer&rsquo;s week</p>
            </div>
            <div class="ww-balance-chip" v-if="!trainingToReschedule && trainingBalance !== null">
                <span class="ww-display ww-balance-chip__value">{{ trainingBalance }}</span>
                <span class="ww-balance-chip__label">trainings left</span>
            </div>
        </div>

        <div class="ww-card ww-alert" v-if="isPropose">
            It&rsquo;s less than 24h until the training, so the other side has to accept the new time.
        </div>

        <div class="ww-schedule">
            <div class="ww-card ww-schedule__steps">
                <div class="ww-schedule__step">
                    <span class="ww-step-label">1 &middot; Pick a day</span>
                    <div class="ww-days">
                        <button type="button" class="ww-day" v-for="day in calendar" :key="formatDate(day.date)"
                                v-bind:class="{'is-selected': formatDate(day.date) === trainingData.date}"
                                @click="selectDay(day)">
                            <span class="ww-day__weekday">{{ weekdayShort(day.date) }}</span>
                            <span class="ww-display ww-day__num">{{ dayOfMonth(day.date) }}</span>
                            <span class="ww-day__month">{{ monthShort(day.date) }}</span>
                        </button>
                    </div>
                    <span class="ww-schedule__hint" v-if="calendar.length === 0">
                        No open hours in the next three weeks &mdash; check back later.</span>
                </div>

                <div class="ww-schedule__step">
                    <span class="ww-step-label">2 &middot; Pick a time</span>
                    <span class="ww-schedule__hint" v-if="!trainingData.date">Pick a day first.</span>
                    <div class="ww-hours" v-if="trainingData.date">
                        <template v-for="hour in availableHours">
                            <button type="button" class="ww-hour" :key="formatHour(hour.hour)"
                                    v-if="!hour.hasTrainingScheduled"
                                    v-bind:class="{'is-selected': formatHour(hour.hour) === trainingData.hour}"
                                    @click="trainingData.hour = formatHour(hour.hour)">
                                {{ formatHour(hour.hour) }}
                            </button>
                        </template>
                    </div>
                </div>

                <div class="ww-schedule__step">
                    <span class="ww-step-label">3 &middot; Notes for the trainer
                        <span class="ww-step-label__optional">(optional)</span></span>
                    <textarea class="ww-textarea" rows="3" v-model="trainingData.notes" maxlength="1000"
                              placeholder="Anything the trainer should know before the session…"></textarea>
                </div>
            </div>

            <div class="ww-summary">
                <span class="ww-step-label ww-summary__label">Your session</span>
                <div class="ww-summary__when" v-if="selectedDay && trainingData.hour">
                    <span class="ww-display ww-summary__date">{{ formatDayLong(selectedDay.date) }}</span>
                    <span class="ww-display ww-summary__time">{{ selectedHourRange }}</span>
                </div>
                <div class="ww-summary__when" v-if="!selectedDay || !trainingData.hour">
                    <span class="ww-summary__placeholder">Pick a day and a time to see your session here.</span>
                </div>
                <div class="ww-summary__facts">
                    <div><span>Duration</span><span>60 min</span></div>
                    <div v-if="!trainingToReschedule && trainingBalance !== null">
                        <span>Cost</span><span>1 credit &middot; {{ trainingBalance - 1 }} left after</span>
                    </div>
                </div>
                <button class="ww-btn ww-btn--primary ww-btn--lg ww-btn--block" type="button"
                        v-bind:disabled="!trainingData.date || !trainingData.hour"
                        @click="scheduleNewTraining">
                    <span v-if="!trainingToReschedule">Confirm training</span>
                    <span v-if="trainingToReschedule && !isPropose">Move training</span>
                    <span v-if="trainingToReschedule && isPropose">Send proposal</span>
                    <span v-if="showLoader" class="ww-spinner"></span>
                </button>
                <span class="ww-summary__note">Free cancellation until 24h before the session.</span>
            </div>
        </div>
    </app-layout>
</template>

<script>
    import AppLayout from '../layouts/App.vue'
    import {getTrainingBalance} from "../repositories/user";
    import {getAvailableDates, rescheduleTraining, scheduleTraining} from "../repositories/trainings";
    import {apiErrorMessage} from "../repositories/errors";
    import {dayOfMonth, formatDate, formatDayLong, formatHour, monthShort, weekdayShort} from "../date";

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
            }, function () {
                self.$toast.error('Failed to load available dates')
            })
            this.trainingToReschedule = this.$route.params['trainingID'];
            this.isReschedule = this.$attrs.isReschedule;
            this.isPropose = this.$attrs.isPropose;
            getTrainingBalance(
                balance => self.trainingBalance = balance,
                () => self.$toast.error('Failed to load training balance'),
            );
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
        computed: {
            selectedDay() {
                return this.calendar.find(obj => formatDate(obj.date) === this.trainingData.date)
            },
            selectedHourRange() {
                if (!this.trainingData.hour) {
                    return ''
                }
                let end = parseInt(this.trainingData.hour.split(':')[0]) + 1
                return this.trainingData.hour + ' – ' + (end.toString().length == 2 ? end : '0' + end) + ':00'
            },
        },
        methods: {
            selectDay(day) {
                this.trainingData.date = formatDate(day.date)
                this.trainingData.hour = ''
                this.trainingDayChange()
            },
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
                            self.$toast.error(apiErrorMessage(err, "Failed to reschedule training"));
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
                            self.$toast.error(apiErrorMessage(err, "Failed to add training"));
                            console.error(err)
                        }
                    )
                }
            },
            formatDate,
            formatHour,
            weekdayShort,
            monthShort,
            dayOfMonth,
            formatDayLong,
        },
    }
</script>

<style scoped>
    .ww-balance-chip {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 8px 16px;
        border-radius: 999px;
        background: var(--ww-accent-soft);
    }

    .ww-balance-chip__value {
        font-size: 20px;
        font-weight: 700;
        color: var(--ww-accent);
        line-height: 1;
    }

    .ww-balance-chip__label {
        font-size: 13px;
        font-weight: 600;
        color: var(--ww-accent-deep);
    }

    .ww-alert {
        margin-bottom: 20px;
        padding: 14px 18px;
        border-color: var(--ww-amber-line);
        background: var(--ww-amber-soft);
        color: var(--ww-amber-deep);
        font-size: 14px;
        font-weight: 500;
    }

    .ww-schedule {
        display: flex;
        gap: 20px;
        align-items: flex-start;
    }

    .ww-schedule__steps {
        flex-grow: 1;
        min-width: 0;
        display: flex;
        flex-direction: column;
        gap: 26px;
        padding: 26px 28px;
    }

    .ww-schedule__step {
        display: flex;
        flex-direction: column;
        gap: 12px;
    }

    .ww-step-label__optional {
        text-transform: none;
        letter-spacing: 0;
        font-weight: 500;
        color: var(--ww-faint);
    }

    .ww-schedule__hint {
        font-size: 13px;
        color: var(--ww-faint);
    }

    .ww-days {
        display: flex;
        gap: 10px;
        overflow-x: auto;
        padding-bottom: 4px;
    }

    .ww-day {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1px;
        min-width: 76px;
        padding: 10px 0 8px 0;
        border-radius: 12px;
        border: 1px solid var(--ww-line);
        background: var(--ww-surface);
        cursor: pointer;
        flex-shrink: 0;
    }

    .ww-day:hover {
        border-color: var(--ww-ink);
    }

    .ww-day__weekday, .ww-day__month {
        font-size: 11px;
        font-weight: 600;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: var(--ww-muted);
    }

    .ww-day__num {
        font-size: 26px;
        font-weight: 700;
        line-height: 1.1;
    }

    .ww-day.is-selected {
        border-color: var(--ww-ink);
        background: var(--ww-ink);
        color: #FFFFFF;
    }

    .ww-day.is-selected .ww-day__weekday, .ww-day.is-selected .ww-day__month {
        color: #B9B4AA;
    }

    .ww-hours {
        display: grid;
        grid-template-columns: repeat(4, minmax(0, 1fr));
        gap: 10px;
    }

    .ww-hour {
        padding: 13px 0;
        text-align: center;
        border-radius: 10px;
        border: 1px solid var(--ww-line);
        background: var(--ww-surface);
        font-size: 14px;
        font-weight: 600;
        color: var(--ww-ink);
        cursor: pointer;
    }

    .ww-hour:hover {
        border-color: var(--ww-ink);
    }

    .ww-hour.is-selected {
        border-color: var(--ww-accent);
        background: var(--ww-accent);
        color: #FFFFFF;
    }

    .ww-summary {
        width: 360px;
        flex-shrink: 0;
        display: flex;
        flex-direction: column;
        gap: 18px;
        padding: 26px 28px;
        border-radius: 14px;
        background: var(--ww-ink);
        color: var(--ww-bg);
    }

    .ww-summary__label {
        color: var(--ww-faint);
    }

    .ww-summary__when {
        display: flex;
        flex-direction: column;
        gap: 4px;
        min-height: 74px;
    }

    .ww-summary__date {
        font-size: 36px;
        font-weight: 700;
        line-height: 1;
    }

    .ww-summary__time {
        font-size: 36px;
        font-weight: 700;
        line-height: 1;
        color: var(--ww-accent);
    }

    .ww-summary__placeholder {
        font-size: 14px;
        color: var(--ww-faint);
    }

    .ww-summary__facts {
        display: flex;
        flex-direction: column;
        gap: 10px;
        border-top: 1px solid rgba(255, 255, 255, 0.12);
        padding-top: 16px;
    }

    .ww-summary__facts > div {
        display: flex;
        justify-content: space-between;
        font-size: 14px;
    }

    .ww-summary__facts > div span:first-child {
        color: #A9A398;
    }

    .ww-summary__facts > div span:last-child {
        font-weight: 600;
    }

    .ww-summary .ww-btn:disabled {
        opacity: 0.45;
        cursor: not-allowed;
    }

    .ww-summary__note {
        font-size: 12px;
        color: var(--ww-faint);
    }

    @media (max-width: 980px) {
        .ww-schedule {
            flex-direction: column;
        }

        .ww-summary {
            width: 100%;
        }
    }
</style>
