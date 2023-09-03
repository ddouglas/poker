var Countdown = /** @class */ (function () {
    function Countdown(_a) {
        var initialValue = _a.initialValue, showHour = _a.showHour, emitter = _a.emitter;
        this.isRunning = false;
        this.initialValue = initialValue;
        this.showHour = showHour;
        this.emitter = emitter;
        this.countdownValue = initialValue;
        this.interval = null;
    }
    Countdown.prototype.decrementCountdown = function () {
        this.countdownValue--;
        this.emitter(this.format(this.countdownValue));
        if (this.countdownValue === 0) {
            this.stop();
        }
    };
    Countdown.prototype.format = function (duration) {
        var parts = {
            hours: Math.floor((duration / (60 * 60)) % 24),
            minutes: Math.floor((duration / (60)) % 60),
            seconds: Math.floor(duration % 60)
        };
        var bits = [];
        if (parts.hours > 0 && this.showHour) {
            var bit = "";
            if (parts.hours < 10) {
                bit = "0";
            }
            bit = "".concat(bit).concat(parts.hours);
        }
        if (parts.minutes == 0) {
            bits.push("00");
        }
        else if (parts.minutes > 0) {
            var bit = "";
            if (parts.minutes < 10) {
                bit = "0";
            }
            bit = "".concat(bit).concat(parts.minutes);
            bits.push(bit);
        }
        if (parts.seconds == 0) {
            bits.push("00");
        }
        else if (parts.seconds > 0) {
            var bit = "";
            if (parts.seconds < 10) {
                bit = "0";
            }
            bit = "".concat(bit).concat(parts.seconds);
            bits.push(bit);
        }
        return bits.join(':');
    };
    Countdown.prototype.start = function () {
        var _this = this;
        if (this.interval) {
            this.stop();
            this.reset();
        }
        this.interval = setInterval(function () { return _this.decrementCountdown(); }, 1000);
        this.isRunning = true;
    };
    Countdown.prototype.stop = function () {
        if (!this.interval) {
            console.warn("failed to stop interval");
            return;
        }
        clearInterval(this.interval);
        this.interval = null;
        this.isRunning = false;
    };
    Countdown.prototype.toggle = function () {
        // if (this.isRunning) {
        //     this.stop()
        // } else if (!this.isRunning) {
        //     this.start()
        // }
        this.isRunning ? this.stop() : this.start();
        return;
    };
    Countdown.prototype.reset = function () {
        this.stop();
        this.countdownValue = this.initialValue; // Reset the countdown value
    };
    return Countdown;
}());
document.body.addEventListener("htmx:load", function () {
    console.log("countdownServerData from countdown.js", countdownServerData);
    console.log("event listener from countdown.js");
});
// const countdown = new Countdown(10); 
