import Countdown from "./countdown"
import { fetchElements } from "./elements"
import "./events"

declare var htmx: any

var countdown: Countdown | null


export function initCountdown() {

    console.debug("initCountdown :: start")

    // Fetch all the elements that we're going to be interacting with on the page
    const elements = fetchElements()
    if (!elements) return

    const {
        // Endpoint that HTMX will use to reach out and fetch the next level
        nextLevelURI,
        // A String representation of the number of seconds that the time needs to run for
        durationSecStr,
        // The HTMLElement representing the text of our timer
        timer,
    } = elements

    // One scenario that can occur is when the timer is complete, meaning all levels have been run through,
    // no seconds are returns. The attribute is not set on the element, so here we just make sure that we
    // didn't receive an empty string
    let parsedDuractionSec: number = 0
    if (durationSecStr) {
        parsedDuractionSec = parseInt(durationSecStr)
    }

    countdown = new Countdown({
        initialValue: parsedDuractionSec,
        showHour: parsedDuractionSec > 3600,
        emitter: (text: string) => {
            timer.innerHTML = text
            console.debug(`received emitted value ${text}`)
        },
        onComplete: () => {

            const nextLevelURIProceed = `${nextLevelURI}?proceed=true`

            htmx.ajax(
                'GET',
                nextLevelURIProceed,
                htmx.find('#timer-container')
            )

        }
    })

    console.debug("initCountdown :: complete")

}

export function toggleCountdown() {

    console.debug("toggleCountdown :: start")

    if (!countdown) {
        console.error("timerToggleEventClick :: countdown not set")
        return
    }

    const elements = fetchElements()
    if (!elements) {
        console.error("failed to fetch elements, unable to register click event on timer toggle")
        return
    }

    const { timerToggle } = elements

    countdown.toggle()

    if (!countdown.getIsRunning()) {
        htmx.removeClass(timerToggle, "fa-circle-stop")
        htmx.addClass(timerToggle, "fa-circle-play")
    } else {
        htmx.removeClass(timerToggle, "fa-circle-play")
        htmx.addClass(timerToggle, "fa-circle-stop")
    }

    console.debug("toggleCountdown :: complete")

}

export function stopCountdown() {

    console.debug("stopCountdown :: start")


    if (!countdown) {
        console.error("failed to stop countdown, countdown is undefined", countdown)
        return
    }

    countdown.stop()

    const elements = fetchElements()
    if (!elements) {
        console.error("failed to fetch elements, unable to register click event on timer toggle")
        return
    }

    const { timerToggle } = elements

    htmx.removeClass(timerToggle, "fa-circle-stop")
    htmx.addClass(timerToggle, "fa-circle-play")

    console.debug("stopCountdown :: complete")

}

export function startCountdown() {

    console.debug("startCountdown :: start")

    const elements = fetchElements()
    if (!elements) {
        console.error("failed to fetch elements, unable to register click event on timer toggle")
        return
    }

    const { timerToggle } = elements

    htmx.removeClass(timerToggle, "fa-circle-stop")
    htmx.addClass(timerToggle, "fa-circle-play")
    countdown?.start()

    console.debug("startCountdown :: complete")

}

