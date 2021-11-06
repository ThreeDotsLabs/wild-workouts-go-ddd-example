import 'bootstrap'
import 'bootstrap/dist/css/bootstrap.min.css'
import './assets/main.css'

import Vue from 'vue'

import VueRouter from 'vue-router'
import VueDialog from "vuejs-dialog"
import 'vuejs-dialog/dist/vuejs-dialog.min.css';
import VueToast from 'vue-toast-notification'
import 'vue-toast-notification/dist/theme-sugar.css'
import TrainingsList from './pages/TrainingsList'
import Calendar from './pages/Calendar'
import ScheduleTraining from './pages/ScheduleTraining'
import Login from './pages/Login'
import SetSchedule from './pages/SetSchedule'
import {loadFirebaseConfig} from "./firebase";
import {Auth, setApiClientsAuth} from "./repositories/auth";
import RatingsList from "./pages/RatingsList";
import RateTraining from "./pages/RateTraining";

Vue.use(VueRouter)

Vue.use(VueDialog, {
    html: true,
    loader: false,
    okText: 'Proceed',
    cancelText: 'Cancel',
    animation: 'fade'
})

Vue.use(VueToast, {
    position: 'top-right',
    duration: 5000,
})

Vue.config.productionTip = false


const routes = [
    {
        path: '/',
        redirect: 'login'
    },
    {
        path: '/login',
        component: Login,
        name: 'login',
    },
    {
        path: '/trainings',
        component: TrainingsList,
        name: 'trainingsList',
    },
    {
        path: '/calendar',
        component: Calendar,
        name: 'calendar',

    },
    {
        path: '/trainings/schedule',
        component: ScheduleTraining,
        name: 'scheduleTraining',
    },
    {
        path: '/trainings/reschedule/:trainingID',
        component: ScheduleTraining,
        name: 'rescheduleTraining',
    },
    {
        path: '/trainings/propose-new-date/:trainingID',
        component: ScheduleTraining,
        name: 'proposeNewDate',
        props: {isPropose: true},
    },
    {
        path: '/trainer/set-schedule',
        component: SetSchedule,
        name: 'setSchedule',
    },
    {
        path: '/ratings',
        component: RatingsList,
        name: 'ratingsList',
    },
    {
        path: '/ratings/:trainingID',
        component: RateTraining,
        name: 'rateTraining',
    },
]

const router = new VueRouter({
    routes,
    mode: 'history',
})

const app = new Vue({
    router,
})


loadFirebaseConfig()
    .then(function () {
        return Auth.waitForAuthReady()
    })
    .then(function () {
        return Auth.getJwtToken(false)
    })
    .then(token => {
        setApiClientsAuth(token)
    })
    .finally(function () {
        app.$mount('#app')
    })

