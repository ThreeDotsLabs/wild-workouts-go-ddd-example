import UsersDefaultApi from './clients/users/src/api/DefaultApi'
import {ApiClient as UsersApiClient} from './clients/users/src'

import 'firebase/auth';

import {Auth, setApiClientsAuth} from "./auth";

export const Trainer = 'trainer';
export const Attendee = 'attendee';

const serverSettings = {
  hostname: window.location.hostname,
};
export let usersClient = new UsersApiClient()
usersClient.basePath = usersClient.getBasePathFromSettings(0, serverSettings);
let usersAPI = new UsersDefaultApi(usersClient)

if (process.env.NODE_ENV === 'development') {
    usersClient.basePath = "http://localhost:3002/api"
}

export function getUserRole() {
    return localStorage.getItem('role')
}

export function getTrainingBalance(callback) {
    return usersAPI.getCurrentUser((error, data) => {
        if (!error) {
            callback(data.balance)
            return
        }
        console.error(error)
    })
}

export function loginUser(login, password) {
    return Auth.login(login, password)
        .then(function () {
            return Auth.waitForAuthReady()
        })
        .then(function () {
            return Auth.getJwtToken(false)
        })
        .then(token => {
            setApiClientsAuth(token)
        })
        .then(function () {
            return new Promise(((resolve, reject) => {
                usersAPI.getCurrentUser((error, data) => {
                    if (!error) {
                        resolve(data)
                        return
                    }
                    reject(error)
                })
            }))
        })
        .then(data => {
            localStorage.setItem('role', data.role)
        })
}

export function getTestUsers() {
    return [
        {
            'uuid': '1',
            'login': 'trainer@threedots.tech',
            'password': '123456',
            'role': 'trainer',
            'name': 'Trainer',
        },
        {
            'uuid': '2',
            'login': 'attendee@threedots.tech',
            'password': '123456',
            'role': 'attendee',
            'name': 'Mock Arnie',
        },
    ]
}
