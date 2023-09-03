interface CountdownOpts {
    initialValue: number
    showHour: boolean
    emitter: (current: string) => void
}

class Countdown {
    initialValue: number
    showHour: boolean
    isRunning: boolean = false
    emitter: (current: string) => void
    private countdownValue: number;
    private interval: number | null;

    constructor({ initialValue, showHour, emitter }: CountdownOpts) {
        this.initialValue = initialValue;
        this.showHour = showHour;
        this.emitter = emitter;
        this.countdownValue = initialValue;
        this.interval = null;
    }

    private decrementCountdown() {
        this.countdownValue--;

        this.emitter(this.format(this.countdownValue))
        if (this.countdownValue === 0) {
            this.stop();
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

    public reset() {
        this.stop();
        this.countdownValue = this.initialValue; // Reset the countdown value
    }
}

declare var countdownServerData: any

document.body.addEventListener("htmx:load", function () {

    console.log("countdownServerData from countdown.js", countdownServerData)

    console.log("event listener from countdown.js");




})
// const countdown = new Countdown(10); 