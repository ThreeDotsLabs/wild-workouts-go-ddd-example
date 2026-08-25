export function formatDate(d) {
    return d.getFullYear().toString() + "-" + ((d.getMonth() + 1).toString().length == 2 ? (d.getMonth() + 1).toString() : "0" + (d.getMonth() + 1).toString()) + "-" + (d.getDate().toString().length == 2 ? d.getDate().toString() : "0" + d.getDate().toString());
}

export function formatHour(d) {
    return (d.getHours().toString().length == 2 ? d.getHours().toString() : "0" + d.getHours().toString()) + ":" + ((parseInt(d.getMinutes() / 5) * 5).toString().length == 2 ? (parseInt(d.getMinutes() / 5) * 5).toString() : "0" + (parseInt(d.getMinutes() / 5) * 5).toString());
}

export function formatDateTime(d) {
    return formatDate(d) + " " + formatHour(d)
}

const WEEKDAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

export function weekdayShort(d) {
    return WEEKDAYS[d.getDay()]
}

export function monthShort(d) {
    return MONTHS[d.getMonth()]
}

export function dayOfMonth(d) {
    return (d.getDate().toString().length == 2 ? d.getDate().toString() : "0" + d.getDate().toString())
}

export function formatHourRange(d) {
    let end = new Date(d.getTime());
    end.setHours(end.getHours() + 1);
    return formatHour(d) + " – " + formatHour(end)
}

export function formatDayLong(d) {
    return weekdayShort(d) + " " + dayOfMonth(d) + " " + monthShort(d)
}

