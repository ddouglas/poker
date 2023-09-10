import { fetchElements } from "./elements"
import { initCountdown, startCountdown, stopCountdown, toggleCountdown } from "./main"

var abort: AbortController

function initAbort() {
    abort = new AbortController()
}

initAbort()

document.addEventListener("DOMContentLoaded", () => {

    console.log("DOMContentLoaded :: start")
    initCountdown()
    initTimerToggleEventClick()
    console.log("DOMContentLoaded :: complete")
})

document.body.addEventListener("countdown::proceed", () => {
    console.debug("countdown::proceed :: start")
    resetCountdown()
    startCountdown()
    console.debug("countdown::proceed :: complete")
})

document.body.addEventListener("countdown::reset", () => {
    console.debug("countdown::reset :: start")
    resetCountdown()
    console.debug("countdown::reset :: complete")
})

function resetCountdown() {
    stopCountdown()
    abort.abort()
    initAbort()
    initCountdown()
    initTimerToggleEventClick()
}

function initTimerToggleEventClick() {

    const elements = fetchElements()
    if (!elements) {
        console.error("failed to fetch elements, unable to register click event on timer toggle")
        return
    }

    const { timerToggle } = elements
    timerToggle.addEventListener("click", () => toggleCountdown(), { signal: abort.signal })

}

