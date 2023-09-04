interface CountdownOpts {
    initialValue: number
    showHour: boolean
    emitter: (current: string) => void
    onComplete: () => void
}

class Countdown {
    private initialValue: number
    private showHour: boolean
    private isRunning: boolean = false
    private emitter: (current: string) => void
    private onComplete: () => void
    private countdownValue: number;
    private interval: number | null;

    constructor({ initialValue, showHour, emitter, onComplete }: CountdownOpts) {
        this.initialValue = initialValue;
        this.showHour = showHour;
        this.emitter = emitter;
        this.onComplete = onComplete;
        this.countdownValue = initialValue;
        this.interval = null;
    }

    private decrementCountdown() {
        this.countdownValue--;

        this.emitter(this.format(this.countdownValue))
        if (this.countdownValue === 0) {
            this.stop();
            this.onComplete();
        }
    }

    private format(duration: number): string {

        const parts = {
            hours: Math.floor((duration / (60 * 60)) % 24),
            minutes: Math.floor((duration / (60)) % 60),
            seconds: Math.floor(duration % 60)
        }

        const bits: string[] = []
        if (parts.hours > 0 && this.showHour) {
            let bit: string = ""
            if (parts.hours < 10) {
                bit = `0`
            }
            bit = `${bit}${parts.hours}`
        }

        if (parts.minutes == 0) {
            bits.push(`00`)
        } else if (parts.minutes > 0) {
            let bit: string = ""
            if (parts.minutes < 10) {
                bit = `0`
            }
            bit = `${bit}${parts.minutes}`
            bits.push(bit)
        }


        if (parts.seconds == 0) {
            bits.push(`00`)
        } else if (parts.seconds > 0) {
            let bit: string = ""
            if (parts.seconds < 10) {
                bit = `0`
            }
            bit = `${bit}${parts.seconds}`
            bits.push(bit)
        }

        return bits.join(':')

    }

    public start() {
        if (this.interval) {
            this.stop()
            this.reset()
        }

        this.interval = setInterval(() => this.decrementCountdown(), 1000)
        this.isRunning = true

    }

    public stop() {
        if (!this.interval) {
            console.warn("failed to stop interval")
            return
        }
        clearInterval(this.interval);
        this.interval = null
        this.isRunning = false
    }

    public toggle() {
        this.isRunning ? this.stop() : this.start()
        return
    }

    public getIsRunning() {
        return this.isRunning
    }

    public reset() {
        this.stop();
        this.countdownValue = this.initialValue; // Reset the countdown value
    }
}

declare var htmx: any

let cd: Countdown | null
let timerContainer = document.getElementById('timer-container')
let timer = document.getElementById('timer')
let timerToggle = document.getElementById("toggle-timer-button")
let nextTimerButton = document.getElementById("trigger-next-timer-level")

function bootstrap() {
    if (!timerContainer) {
        console.error("Failed to fetch timer-container by id")
        return
    }

    if (!timer) {
        console.error("failed to fetch timer element by id")
        return
    }

    if (!timerToggle) {
        console.error("failed to fetch toggle-timer-button element by id")
        return
    }
    if (!nextTimerButton) {
        console.error("failed to fetch trigger-next-timer-level element by id")
        return
    }

    const nextLevelURI = nextTimerButton.getAttribute("hx-get")
    if (!nextLevelURI) {
        console.error("trigger-next-timer-level element is missing attribute hx-get")
        return
    }

    const durationSecStr = timer.getAttribute("data-level-duration-sec")
    if (!durationSecStr) {
        console.error("trigger-next-timer-level element is missing attribute data-level-duration-sec")
        return
    }

    const durationSec = parseInt(durationSecStr)

    cd = new Countdown({
        initialValue: durationSec,
        showHour: false,
        emitter: function (text: string) {
            if (!timer) {
                console.error("failed to fetch timer element by id")
                return
            }

            timer.innerHTML = text
            console.log(`Receiving Emitted Text ${text}`)
        },
        onComplete: function () {
            htmx.ajax(
                'GET',
                nextLevelURI,
                timerContainer
            ).then(() => {
                // timerToggle.removeEventListener("click", toggleTimerFunc)
                // timerContainer.removeEventListener("htmx:afterSettle", cd.start)
                console.log("Swapped in timer-container")
            })

            console.log("countdown is complete")
        },
    })

    console.log("htmx:afterSettle timerContainer.addEventListener")
    timerContainer.addEventListener("htmx:afterSettle", resetCD)


}


function toggleTimer() {
    if (!cd) {
        console.error("toggleTimer :: cd is not set yet")
        return
    }
    cd.toggle()
    if (!cd.getIsRunning()) {
        htmx.removeClass(timerToggle, "fa-circle-stop")
        htmx.addClass(timerToggle, "fa-circle-play")
    } else {
        htmx.removeClass(timerToggle, "fa-circle-play")
        htmx.addClass(timerToggle, "fa-circle-stop")
    }
}

function resetCD() {
    if (!cd || !timerContainer) {
        console.error("resetCD :: cd and/or timerContainer are not set yet")
        return
    }
    // timerContainer.removeEventListener("htmx:afterSettle", resetCD)
    // document.body.removeEventListener("htmx:load", bootstrap)

    bootstrap()
    console.log("cd.start()")
    cd.start()
}

document.body.addEventListener("htmx:load", bootstrap)

if (timerToggle) {
    timerToggle.addEventListener("click", toggleTimer)
}

if (timerContainer) {
    timerContainer.addEventListener("htmx:afterSettle", resetCD)
}