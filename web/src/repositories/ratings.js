import RatingsDefaultApi from './clients/ratings/src/api/DefaultApi'
import {ApiClient as RatingsApiClient, PostRating} from './clients/ratings/src'

import {ApiClient as TrainerApiClient} from './clients/trainer/src'

const serverSettings = {
    hostname: window.location.hostname,
};
export let ratingsClient = new RatingsApiClient()
ratingsClient.basePath = ratingsClient.getBasePathFromSettings(0, serverSettings);
let ratingsAPI = new RatingsDefaultApi(ratingsClient)

export let trainerClient = new TrainerApiClient()
trainerClient.basePath = trainerClient.getBasePathFromSettings(0, serverSettings);

if (process.env.NODE_ENV === 'development') {
    ratingsClient.basePath = "http://localhost:3003/api"
}

export function getRatings(callback) {
    ratingsAPI.getRatings((error, data) => {
        if (error) {
            console.error(error);
        } else {
            callback(data.ratings)
        }
    });
}

export function canRate(training, callback) {
    ratingsAPI.getTrainingRating(training.uuid, (error, data) => {
        if (error) {
            console.error(error);
        } else {
            callback(data.canRate)
        }
    });
}

export function submitRating(trainingUUID, rating, comment, successCallback, errorCallback) {
    ratingsAPI.postTrainingRating(trainingUUID, new PostRating(parseInt(rating), comment), (error) => {
        if (error) {
            errorCallback(error)
        } else {
            successCallback()
        }
    })
}