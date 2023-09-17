
interface ElementsAndAttributes {
    timerContainer: HTMLElement
    timer: HTMLElement
    timerToggle: HTMLElement
    nextTimerButton: HTMLElement | null
    nextLevelURI: string
    durationSecStr: string
    audioPlay: HTMLAudioElement | null
    audioContinue: HTMLAudioElement | null
    audioBeep: HTMLAudioElement | null
}

export function fetchElements(): ElementsAndAttributes | null {
    const timerContainer = document.getElementById('timer-container')
    const timer = document.getElementById('timer')
    const timerToggle = document.getElementById("toggle-timer-button")
    const audioPlay = document.getElementById("audio-play") as HTMLAudioElement | null
    const audioContinue = document.getElementById("audio-continue") as HTMLAudioElement | null
    const audioBeep = document.getElementById("audio-beep") as HTMLAudioElement | null
    const nextTimerButton = document.getElementById("trigger-next-timer-level")

    if (!timerContainer) {
        console.error("failed to fetch timer-container by id")
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

    let nextLevelURI: string = ""
    if (nextTimerButton) {
        console.error("failed to fetch trigger-next-timer-level element by id")
        nextLevelURI = nextTimerButton.getAttribute("hx-get") || ""
    }

    let durationSecStr = timer.getAttribute("data-level-duration-sec")
    if (!durationSecStr) {
        console.error("trigger-next-timer-level element is missing attribute data-level-duration-sec")
        durationSecStr = "0"
    }


    return {
        timer,
        timerContainer,
        timerToggle,
        nextTimerButton,
        nextLevelURI,
        durationSecStr,
        audioPlay,
        audioContinue,
        audioBeep
    }

}