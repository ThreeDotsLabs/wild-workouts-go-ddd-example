<template>
    <app-layout>
        <div class="py-5 text-center">
            <h2>Trainer's schedule</h2>
            <p class="lead">Below is an example form built entirely with Bootstrap’s form controls. Each required form
                group
                has a validation state that can be triggered by attempting to submit the form without completing it.</p>
        </div>

        <FullCalendar defaultView="timeGridWeek" :plugins="calendarPlugins" :header="calendarHeader"
                      :events="calendarEvents" navLinks="true"/>

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

    function getScheduleCalendarEvents(callback) {
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
                            start: start.toISOString(),
                            end: end.toISOString(),
                        })
                    }
                }
            }

            return callback(scheduleEvents)
        });
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
                    }
                })

                getScheduleCalendarEvents(function (scheduleEvents) {
                    self.calendarEvents = events.concat(scheduleEvents)
                });
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
            }
        }
    }
</script>

<style lang='scss'>
    @import '~@fullcalendar/core/main.css';
    @import '~@fullcalendar/daygrid/main.css';
    @import '~@fullcalendar/timegrid/main.css';
    @import '~@fullcalendar/list/main.css';

    .fc-unthemed td.fc-today {
        background: #ffffff;
    }
</style>