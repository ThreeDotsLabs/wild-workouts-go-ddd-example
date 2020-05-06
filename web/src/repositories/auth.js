import {getTestUsers, usersClient} from "./user";
import firebase from "firebase";
import {sign} from "jsonwebtoken";
import {trainerClient, trainingsClient} from "./trainings";

class FirebaseAuth {
    login(login, password) {
        return firebase.auth().signInWithEmailAndPassword(login, password)
    }

    waitForAuthReady() {
        return new Promise((resolve) => {
            firebase
                .auth()
                .onAuthStateChanged(function () {
                    resolve()
                });
        })
    }

    getJwtToken(required) {
        return new Promise((resolve, reject) => {
            if (!firebase.auth().currentUser) {
                if (required) {
                    reject('no user found')
                } else {
                    resolve(null)
                }
                return
            }

            firebase.auth().currentUser.getIdToken(false)
                .then(function (idToken) {
                    resolve(idToken)
                })
                .catch(function (error) {
                    reject(error)
                });
        })
    }

    logout() {
        return new Promise(resolve => {
            if (!firebase.auth().currentUser) {
                resolve()
                return
            }

            return firebase.auth().signOut()
        })
    }


    isLoggedIn() {
        return firebase.auth().currentUser != null
    }
}

class MockAuth {
    login(login, password) {
        return new Promise((resolve, reject) => {
            setTimeout(function () {
                let found = getTestUsers().filter(u => u.login === login && u.password === password);

                if (found) {
                    localStorage.setItem('_mock_user', JSON.stringify(found[0]));
                    resolve()
                } else {
                    reject('invalid login or password')
                }
            }, 500) // simulate http request
        })
    }

    waitForAuthReady() {
        return new Promise((resolve) => {
            setTimeout(resolve, 50)
        })
    }

    getJwtToken() {
        return new Promise((resolve) => {
            let user = this.currentMockUser()

            let claims = {
                'user_uuid': user.uuid,
                'email': user.login,
                'role': user.role,
                'name': user.name,
            }
            let token = sign(claims, 'mock_secret')
            resolve(token)
        })
    }

    currentMockUser() {
        let userStr = localStorage.getItem('_mock_user');
        if (!userStr) {
            return null
        }

        return JSON.parse(userStr)
    }

    logout() {
        return new Promise(resolve => {
            localStorage.setItem('_mock_user', '')

            setTimeout(resolve, 50)
        })
    }

    isLoggedIn() {
        return this.currentMockUser() !== null
    }
}

export function setApiClientsAuth(idToken) {
    usersClient.authentications['bearerAuth'].accessToken = idToken
    trainerClient.authentications['bearerAuth'].accessToken = idToken
    trainingsClient.authentications['bearerAuth'].accessToken = idToken
}

const MOCK_AUTH = process.env.NODE_ENV === 'development'
export let Auth

if (MOCK_AUTH) {
    Auth = new MockAuth()
} else {
    Auth = new FirebaseAuth()
}
