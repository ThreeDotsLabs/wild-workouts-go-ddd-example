<template>
    <app-layout>
        <div class="ww-page-head">
            <div>
                <h1>Calendar</h1>
                <p class="ww-page-sub" v-if="isTrainer">Your open hours and booked sessions</p>
                <p class="ww-page-sub" v-if="!isTrainer">Your trainer&rsquo;s open hours and your sessions</p>
            </div>
        </div>

        <div class="ww-card ww-calendar">
            <FullCalendar defaultView="timeGridWeek" :plugins="calendarPlugins" :header="calendarHeader"
                          :events="calendarEvents" navLinks="true" :allDaySlot="false" :height="'auto'"
                          :firstDay="1" minTime="06:00:00" maxTime="22:00:00"/>

            <div class="ww-calendar__legend">
                <div><span class="ww-calendar__swatch ww-calendar__swatch--open"></span>Open hours</div>
                <div><span class="ww-calendar__swatch ww-calendar__swatch--booked"></span>Booked training</div>
            </div>
        </div>
    </app-layout>
</template>

<script>
    import AppLayout from '../layouts/App.vue'
    import FullCalendar from '@fullcalendar/vue'
    import interactionPlugin from '@fullcalendar/interaction'
    import dayGridPlugin from '@fullcalendar/daygrid'
    import timeGridPlugin from '@fullcalendar/timegrid'
    import listPlugin from '@fullcalendar/list'
    import {getCalendar, getSchedule} from "../repositories/trainings";
    import {getUserRole, Trainer} from "../repositories/user";

    function getScheduleCalendarEvents(callback, errorCallback) {
        let start = new Date();
        start.setMonth(start.getMonth() - 3);
        let end = new Date();
        end.setMonth(end.getMonth() + 3);

        getSchedule(start, end, function (schedule) {
            let scheduleEvents = []

            for (let idx in schedule) {
                let date = schedule[idx]

                for (let idx in date.hours) {
                    let hour = date.hours[idx]

                    if (hour.available) {
                        let start = hour.hour;
                        let end = new Date(hour.hour.getTime());
                        end.setHours(end.getHours() + 1)

                        scheduleEvents.push({
                            rendering: 'background',
                            backgroundColor: '#2F9E62',
                            start: start.toISOString(),
                            end: end.toISOString(),
                        })
                    }
                }
            }

            return callback(scheduleEvents)
        }, errorCallback);
    }

    export default {
        components: {
            AppLayout,
            FullCalendar,
        },
        mounted() {
            let self = this

            getCalendar(function (data) {
                let events = data.map(function (obj) {
                    let end = new Date(obj.time.getTime());
                    end.setHours(end.getHours()+1)

                    let isTrainer = getUserRole() === Trainer

                    return {
                        title: isTrainer ? obj.user : 'Training',
                        start: obj.time,
                        end: end.toISOString(),
                        backgroundColor: '#E0491F',
                        borderColor: '#E0491F',
                        textColor: '#FFFFFF',
                    }
                })

                getScheduleCalendarEvents(function (scheduleEvents) {
                    self.calendarEvents = events.concat(scheduleEvents)
                }, function () {
                    self.$toast.error('Failed to load the schedule')
                });
            }, function () {
                self.$toast.error('Failed to load trainings')
            })
        },
        data() {
            return {
                calendarPlugins: [interactionPlugin, dayGridPlugin, timeGridPlugin, listPlugin],
                calendarHeader: {
                    left: 'prev,next today',
                    center: 'title',
                    right: 'dayGridMonth,timeGridWeek,timeGridDay,listWeek'
                },
                calendarEvents: [],
                isTrainer: getUserRole() === Trainer,
            }
        }
    }
</script>

<style>
    @import '~@fullcalendar/core/main.css';
    @import '~@fullcalendar/daygrid/main.css';
    @import '~@fullcalendar/timegrid/main.css';
    @import '~@fullcalendar/list/main.css';

    .ww-calendar {
        padding: 20px 24px 24px 24px;
        margin-bottom: 8px;
    }

    .ww-calendar__legend {
        display: flex;
        align-items: center;
        gap: 18px;
        padding-top: 16px;
        font-size: 13px;
        color: var(--ww-muted);
    }

    .ww-calendar__legend > div {
        display: flex;
        align-items: center;
        gap: 7px;
    }

    .ww-calendar__swatch {
        width: 12px;
        height: 12px;
        border-radius: 4px;
    }

    .ww-calendar__swatch--open {
        background: rgba(47, 158, 98, 0.15);
        border: 1px dashed rgba(47, 158, 98, 0.5);
    }

    .ww-calendar__swatch--booked {
        background: var(--ww-accent);
    }

    /* FullCalendar restyling */

    .ww-calendar .fc {
        font-family: var(--ww-font-ui);
    }

    .ww-calendar .fc-toolbar {
        flex-wrap: wrap;
        gap: 12px;
    }

    .ww-calendar .fc-toolbar h2 {
        font-family: var(--ww-font-display);
        text-transform: uppercase;
        font-size: 26px;
        font-weight: 700;
        color: var(--ww-ink);
    }

    .ww-calendar .fc-button-primary {
        background: var(--ww-surface);
        border: 1px solid var(--ww-line);
        color: var(--ww-ink);
        font-size: 13px;
        font-weight: 600;
        text-transform: capitalize;
        padding: 7px 14px;
        border-radius: 9px;
        box-shadow: none;
    }

    .ww-calendar .fc-button-primary:not(:disabled):hover {
        background: var(--ww-bg);
        border-color: var(--ww-ink);
        color: var(--ww-ink);
    }

    .ww-calendar .fc-button-primary:disabled {
        background: var(--ww-disabled-bg);
        border-color: var(--ww-line-soft);
        color: var(--ww-disabled-text);
    }

    .ww-calendar .fc-button-primary:not(:disabled).fc-button-active,
    .ww-calendar .fc-button-primary:not(:disabled):active {
        background: var(--ww-ink);
        border-color: var(--ww-ink);
        color: #FFFFFF;
        box-shadow: none;
    }

    .ww-calendar .fc-button-primary:focus,
    .ww-calendar .fc-button-primary:not(:disabled).fc-button-active:focus,
    .ww-calendar .fc-button-primary:not(:disabled):active:focus {
        box-shadow: none;
    }

    .ww-calendar .fc th, .ww-calendar .fc td {
        border-color: var(--ww-line-soft);
    }

    .ww-calendar .fc-day-header {
        padding: 8px 0;
        font-size: 12px;
        font-weight: 600;
        letter-spacing: 0.05em;
        text-transform: uppercase;
        color: var(--ww-muted);
    }

    .ww-calendar .fc-day-header a {
        color: var(--ww-muted);
    }

    .ww-calendar .fc-day-header.fc-today, .ww-calendar .fc-day-header.fc-today a {
        color: var(--ww-accent);
    }

    .ww-calendar .fc-axis {
        font-size: 11px;
        color: var(--ww-faint);
    }

    .ww-calendar .fc-unthemed td.fc-today {
        background: rgba(224, 73, 31, 0.035);
    }

    .ww-calendar .fc-event {
        border-radius: 8px;
        border: none;
        padding: 3px 6px;
        font-size: 12px;
        font-weight: 600;
    }

    .ww-calendar .fc-bgevent {
        opacity: 0.12;
    }

    .ww-calendar .fc-list-item-title, .ww-calendar .fc-list-item-time, .ww-calendar .fc-list-heading td {
        font-size: 14px;
    }
</style>
