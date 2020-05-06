import TrainingsDefaultApi from './clients/trainings/src/api/DefaultApi'
import {ApiClient as TrainingsApiClient} from './clients/trainings/src'

import TrainerDefaultApi from './clients/trainer/src/api/DefaultApi'
import {ApiClient as TrainerApiClient} from './clients/trainer/src'
import HourUpdate from "./clients/trainer/src/model/HourUpdate";
import {formatDate} from "../date";
import PostTraining from "./clients/trainings/src/model/PostTraining";

const serverSettings = {
    hostname: window.location.hostname,
};
export let trainingsClient = new TrainingsApiClient()
trainingsClient.basePath = trainingsClient.getBasePathFromSettings(0, serverSettings);
let trainingsAPI = new TrainingsDefaultApi(trainingsClient)

export let trainerClient = new TrainerApiClient()
trainerClient.basePath = trainerClient.getBasePathFromSettings(0, serverSettings);
let trainerAPI = new TrainerDefaultApi(trainerClient)

if (process.env.NODE_ENV === 'development') {
    trainingsClient.basePath = "http://localhost:3001/api"
    trainerClient.basePath = "http://localhost:3000/api"
}

export function getSchedule(dateFrom, dateTo, callback) {
    trainerAPI.getTrainerAvailableHours(dateFrom, dateTo, (error, data) => {
        if (error) {
            console.error(error)
        } else {
            callback(data)
        }
    })
}

export function getAvailableDates(callback) {
    let to = new Date()
    to.setDate(to.getDate() + (7 * 3))

    getSchedule(formatDate(new Date()), formatDate(to), function (data) {
        callback(data.filter(obj => obj.hasFreeHours))
    })
}

export function setHourAvailability(updates, availability, callback, errorCallback) {
    let hourUpdates = []

    updates.forEach(function (val) {
        hourUpdates.push(new Date(val[0] + 'T' + val[1]))
    })

    let hourUpdate = new HourUpdate(hourUpdates)

    if (availability) {
        trainerAPI.makeHourAvailable(hourUpdate, (error) => {
            if (error) {
                console.error(error)
                errorCallback && errorCallback()
            } else {
                callback && callback()
            }
        })
    } else {
        trainerAPI.makeHourUnavailable(hourUpdate, (error) => {
            if (error) {
                console.error(error)
            } else {
                callback && callback()
            }
        })
    }
}

export function getCalendar(callback) {
    trainingsAPI.getTrainings((error, data) => {
        if (error) {
            console.error(error);
        } else {
            callback(data.trainings)
        }
    });
}

export function getPeriods() {
    let periods = []

    for (let week = 0; week <= 2; week++) {
        let from = new Date()
        from.setDate(from.getDate() + (week * 8))

        let to = new Date()
        to.setDate(to.getDate() + (week * 8) + 7)

        periods.push({'from': formatDate(from), 'to': formatDate(to)})
    }

    return periods
}

export function scheduleTraining(notes, date, hour, successCallback, errorCallback) {
    let req = new PostTraining(notes, new Date(date + 'T' + hour));

    trainingsAPI.createTraining(req, (error) => {
        if (error) {
            errorCallback(error)
        } else {
            successCallback()
        }
    })
}

export function rescheduleTraining(trainingUUID, notes, date, hour, successCallback, errorCallback) {
    let req = new PostTraining(notes, new Date(date + 'T' + hour));

    trainingsAPI.rescheduleTraining(trainingUUID, req, (error) => {
        if (error) {
            errorCallback(error)
        } else {
            successCallback()
        }
    })
}

export function cancelTraining(uuid, successCallback, errorCallback) {
    trainingsAPI.cancelTraining(uuid, (error) => {
        if (error) {
            errorCallback(error)
        } else {
            successCallback()
        }
    })
}

export function approveReschedule(uuid, successCallback, errorCallback) {
    trainingsAPI.approveRescheduleTraining(uuid, (error) => {
        if (error) {
            errorCallback(error)
        } else {
            successCallback()
        }
    })
}

export function rejectReschedule(uuid, successCallback, errorCallback) {
    trainingsAPI.rejectRescheduleTraining(uuid, (error) => {
        if (error) {
            errorCallback(error)
        } else {
            successCallback()
        }
    })
}
