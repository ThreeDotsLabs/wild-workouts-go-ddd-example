/**
 * Extracts the message returned by the API for errors that were written for the user
 * (see the Error schema in api/openapi). Falls back to the given message for
 * unexpected failures, which don't carry a message meant to be displayed.
 *
 * @param error error passed to an API client callback
 * @param fallback message to show when the API didn't return one
 * @returns {String}
 */
export function apiErrorMessage(error, fallback) {
    const body = error && error.response && error.response.body

    if (body && body.message) {
        return body.message
    }

    return fallback
}
