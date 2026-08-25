<template>
    <app-layout>
        <div class="ww-page-head">
            <div>
                <h1>Set your availability</h1>
                <p class="ww-page-sub">Tap an hour to open or close it &mdash; changes save instantly</p>
            </div>
            <div class="ww-seg">
                <button type="button" class="ww-seg__opt" v-for="period in periods" :key="period.from+period.to"
                        v-bind:class="{'is-active': selectedPeriod === selectedPeriodValue(period)}"
                        @click="changedPeriod(period)">
                    {{ periodLabel(period) }}
                </button>
            </div>
        </div>

        <div class="ww-legend">
            <div><span class="ww-legend__swatch"></span>Closed</div>
            <div><span class="ww-legend__swatch ww-legend__swatch--open"></span>Open for booking</div>
            <div><span class="ww-legend__swatch ww-legend__swatch--booked"></span>Training booked &mdash; locked</div>
        </div>

        <div class="ww-week">
            <div class="ww-week__day" v-for="day in schedule" :key="day.date.toString()">
                <div class="ww-week__head">
                    <span class="ww-week__weekday">{{ weekdayShort(day.date) }}</span>
                    <span class="ww-display ww-week__num">{{ dayOfMonth(day.date) }}</span>
                    <button type="button" class="ww-week__all" @click="selectAllInDay(day)">Open all</button>
                </div>

                <label class="ww-slot" v-for="hour in day.hours" :key="formatDateTime(hour.hour)"
                       v-bind:class="{'is-open': hour.available && !hour.hasTrainingScheduled, 'is-booked': hour.hasTrainingScheduled}"
                       v-bind:title="hour.hasTrainingScheduled ? 'Training scheduled on this date' : ''">
                    <input type="checkbox" autocomplete="off" v-model="hour.available"
                           @change.prevent="toggleHour($event, hour)"
                           v-bind:data-hour="formatHour(hour.hour)" v-bind:data-date="formatDate(day.date)"
                           v-bind:disabled="hour.hasTrainingScheduled">
                    <span>{{ formatHour(hour.hour) }}</span>
                    <small v-if="hour.hasTrainingScheduled">booked</small>
                </label>
            </div>
        </div>
    </app-layout>
</template>

<script>
    import AppLayout from '../layouts/App.vue'
    import {getPeriods, getSchedule, setHourAvailability} from "../repositories/trainings";
    import {dayOfMonth, formatDate, formatDateTime, formatHour, monthShort, weekdayShort} from "../date";

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
            this.periods = getPeriods()
            this.selectedDateFrom = this.periods[0].from
            this.selectedDateTo = this.periods[0].to
            this.selectedPeriod = this.selectedPeriodValue(this.periods[0])

            this.refreshSchedule()
        },
        methods: {
            refreshSchedule() {
                let self = this
                getSchedule(this.selectedDateFrom, this.selectedDateTo, function (data) {
                    self.schedule = data
                }, function () {
                    self.$toast.error('Failed to load the schedule')
                })
            },
            selectAllInDay(day) {
                let self = this
                let updates = []

                for (let idx in day.hours) {
                    let d = day.hours[idx].hour;
                    updates.push([formatDate(d), formatHour(d)])
                }

                setHourAvailability(updates, true, function () {
                    self.refreshSchedule()
                }, function () {
                    self.$toast.error('Failed to update schedule')
                })
            },
            changedPeriod(period) {
                this.selectedDateFrom = period.from
                this.selectedDateTo = period.to
                this.selectedPeriod = this.selectedPeriodValue(period)

                this.refreshSchedule()
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
                        hour.available = !event.target.checked
                        self.$toast.error("Failed to update schedule")
                    },
                )
            },
            selectedPeriodValue(period) {
                return period.from + '/' + period.to
            },
            periodLabel(period) {
                let from = new Date(period.from + 'T00:00')
                let to = new Date(period.to + 'T00:00')

                let label = monthShort(from) + ' ' + from.getDate() + ' – '
                if (from.getMonth() !== to.getMonth()) {
                    label += monthShort(to) + ' '
                }
                return label + to.getDate()
            },
            formatDate,
            formatDateTime,
            formatHour,
            weekdayShort,
            dayOfMonth,
        },
    }
</script>

<style scoped>
    .ww-legend {
        display: flex;
        align-items: center;
        gap: 18px;
        flex-wrap: wrap;
        padding-bottom: 16px;
        font-size: 13px;
        color: var(--ww-muted);
    }

    .ww-legend > div {
        display: flex;
        align-items: center;
        gap: 7px;
    }

    .ww-legend__swatch {
        width: 12px;
        height: 12px;
        border-radius: 4px;
        border: 1px solid var(--ww-line);
        background: var(--ww-surface);
    }

    .ww-legend__swatch--open {
        border-color: var(--ww-green-line);
        background: var(--ww-green-soft);
    }

    .ww-legend__swatch--booked {
        border-color: var(--ww-ink);
        background: var(--ww-ink);
    }

    .ww-week {
        display: flex;
        gap: 12px;
        align-items: flex-start;
        overflow-x: auto;
        padding-bottom: 8px;
    }

    .ww-week__day {
        display: flex;
        flex-direction: column;
        gap: 8px;
        flex: 1 1 0;
        min-width: 108px;
    }

    .ww-week__head {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 2px;
        padding-bottom: 6px;
    }

    .ww-week__weekday {
        font-size: 11px;
        font-weight: 600;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: var(--ww-muted);
    }

    .ww-week__num {
        font-size: 24px;
        font-weight: 700;
        line-height: 1;
    }

    .ww-week__all {
        border: none;
        background: none;
        padding: 2px 4px;
        font-size: 12px;
        font-weight: 600;
        color: var(--ww-accent);
        cursor: pointer;
    }

    .ww-week__all:hover {
        color: var(--ww-accent-hover);
    }

    .ww-slot {
        height: 42px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 6px;
        border-radius: 10px;
        border: 1px solid var(--ww-line);
        background: var(--ww-surface);
        font-size: 13px;
        font-weight: 500;
        color: var(--ww-faint);
        cursor: pointer;
        margin: 0;
    }

    .ww-slot:hover {
        border-color: var(--ww-ink);
    }

    .ww-slot input {
        position: absolute;
        opacity: 0;
        pointer-events: none;
    }

    .ww-slot.is-open {
        border-color: var(--ww-green-line);
        background: var(--ww-green-soft);
        color: var(--ww-green-deep);
        font-weight: 600;
    }

    .ww-slot.is-booked {
        flex-direction: column;
        gap: 1px;
        border-color: var(--ww-ink);
        background: var(--ww-ink);
        color: #FFFFFF;
        font-weight: 600;
        cursor: not-allowed;
    }

    .ww-slot.is-booked small {
        font-size: 10px;
        font-weight: 500;
        color: #B9B4AA;
        line-height: 1;
    }
</style>
