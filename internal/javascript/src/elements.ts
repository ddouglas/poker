
interface ElementsAndAttributes {
    timerContainer: HTMLElement
    timer: HTMLElement
    timerToggle: HTMLElement
    nextTimerButton: HTMLElement
    nextLevelURI: string
    durationSecStr: string
}

export function fetchElements(): ElementsAndAttributes | null {
    const timerContainer = document.getElementById('timer-container')
    const timer = document.getElementById('timer')
    const timerToggle = document.getElementById("toggle-timer-button")
    const nextTimerButton = document.getElementById("trigger-next-timer-level")

    if (!timerContainer) {
        console.error("Failed to fetch timer-container by id")
        return null
    }

    if (!timer) {
        console.error("failed to fetch timer element by id")
        return null
    }

    if (!timerToggle) {
        console.error("failed to fetch toggle-timer-button element by id")
        return null
    }
    if (!nextTimerButton) {
        console.error("failed to fetch trigger-next-timer-level element by id")
        return null
    }

    let nextLevelURI = nextTimerButton.getAttribute("hx-get")
    if (!nextLevelURI) {
        console.error("trigger-next-timer-level element is missing attribute hx-get")
        nextLevelURI = ""
    }

    let durationSecStr = timer.getAttribute("data-level-duration-sec")
    if (!durationSecStr) {
        console.error("trigger-next-timer-level element is missing attribute data-level-duration-sec")
        durationSecStr = "0"
        // return null
    }

    return { timer, timerContainer, timerToggle, nextTimerButton, nextLevelURI, durationSecStr }

}