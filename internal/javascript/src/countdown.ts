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
    private interval: ReturnType<typeof setTimeout> | null;

    constructor({ initialValue, showHour, emitter, onComplete }: CountdownOpts) {
        this.initialValue = initialValue;
        this.showHour = showHour;
        this.emitter = emitter;
        this.onComplete = onComplete;
        this.countdownValue = initialValue;
        this.interval = null;
    }

    private decrementCountdown() {
        // console.debug("Countdown.decrementCountdown() start")
        this.countdownValue--;

        this.emitter(this.format(this.countdownValue))
        if (this.countdownValue === 0) {
            this.stop();
            this.onComplete();
        }
        // console.debug("Countdown.decrementCountdown() stop")
    }

    private format(duration: number): string {
        // console.debug("Countdown.format() start")
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

        // console.debug("Countdown.format() stop")
        return bits.join(':')

    }

    public start() {
        console.debug("Countdown.start()")
        if (this.interval) {
            console.debug("Countdown.start() if this.interval")
            this.reset()
        }

        if (this.countdownValue === 0) {
            console.error("countdownValue is 0, shutdown down")
            this.stop()
            return
        }
        console.debug("Countdown.start() this.interval")
        this.interval = setInterval(() => this.decrementCountdown(), 1000)
        this.isRunning = true
        console.debug("Countdown.start() done", this.interval)

    }

    public continue() {
        this.stop()
        this.start()
    }

    public stop() {
        console.debug("Countdown.stop() start", this.interval)
        if (this.interval) {
            console.debug("Countdown.stop() clearInterval")

            clearInterval(this.interval);
        }

        this.interval = null
        this.isRunning = false
        console.debug("Countdown.stop() stop")
    }

    public toggle() {
        console.debug("Countdown.toggle() start")
        this.isRunning ? this.stop() : this.start()
        console.debug("Countdown.toggle() stop")
        return
    }

    public getIsRunning() {
        return this.isRunning
    }

    public reset() {
        console.debug("Countdown.reset() start")
        this.stop();
        this.countdownValue = this.initialValue; // Reset the countdown value
        console.log("reset :: ", this)
        console.debug("Countdown.reset() stop")
    }
}

export default Countdown