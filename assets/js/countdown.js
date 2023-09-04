var Countdown = /** @class */ (function () {
    function Countdown(_a) {
        var initialValue = _a.initialValue, showHour = _a.showHour, emitter = _a.emitter, onComplete = _a.onComplete;
        this.isRunning = false;
        this.initialValue = initialValue;
        this.showHour = showHour;
        this.emitter = emitter;
        this.onComplete = onComplete;
        this.countdownValue = initialValue;
        this.interval = null;
    }
    Countdown.prototype.decrementCountdown = function () {
        this.countdownValue--;
        this.emitter(this.format(this.countdownValue));
        if (this.countdownValue === 0) {
            this.stop();
            this.onComplete();
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
        this.isRunning ? this.stop() : this.start();
        return;
    };
    Countdown.prototype.getIsRunning = function () {
        return this.isRunning;
    };
    Countdown.prototype.reset = function () {
        this.stop();
        this.countdownValue = this.initialValue; // Reset the countdown value
    };
    return Countdown;
}());
var cd;
var timerContainer = document.getElementById('timer-container');
var timer = document.getElementById('timer');
var timerToggle = document.getElementById("toggle-timer-button");
var nextTimerButton = document.getElementById("trigger-next-timer-level");
function bootstrap() {
    if (!timerContainer) {
        console.error("Failed to fetch timer-container by id");
        return;
    }
    if (!timer) {
        console.error("failed to fetch timer element by id");
        return;
    }
    if (!timerToggle) {
        console.error("failed to fetch toggle-timer-button element by id");
        return;
    }
    if (!nextTimerButton) {
        console.error("failed to fetch trigger-next-timer-level element by id");
        return;
    }
    var nextLevelURI = nextTimerButton.getAttribute("hx-get");
    if (!nextLevelURI) {
        console.error("trigger-next-timer-level element is missing attribute hx-get");
        return;
    }
    var durationSecStr = timer.getAttribute("data-level-duration-sec");
    if (!durationSecStr) {
        console.error("trigger-next-timer-level element is missing attribute data-level-duration-sec");
        return;
    }
    var durationSec = parseInt(durationSecStr);
    cd = new Countdown({
        initialValue: durationSec,
        showHour: false,
        emitter: function (text) {
            if (!timer) {
                console.error("failed to fetch timer element by id");
                return;
            }
            timer.innerHTML = text;
            console.log("Receiving Emitted Text ".concat(text));
        },
        onComplete: function () {
            htmx.ajax('GET', nextLevelURI, timerContainer).then(function () {
                // timerToggle.removeEventListener("click", toggleTimerFunc)
                // timerContainer.removeEventListener("htmx:afterSettle", cd.start)
                console.log("Swapped in timer-container");
            });
            console.log("countdown is complete");
        },
    });
    console.log("htmx:afterSettle timerContainer.addEventListener");
    timerContainer.addEventListener("htmx:afterSettle", resetCD);
}
function toggleTimer() {
    if (!cd) {
        console.error("toggleTimer :: cd is not set yet");
        return;
    }
    cd.toggle();
    if (!cd.getIsRunning()) {
        htmx.removeClass(timerToggle, "fa-circle-stop");
        htmx.addClass(timerToggle, "fa-circle-play");
    }
    else {
        htmx.removeClass(timerToggle, "fa-circle-play");
        htmx.addClass(timerToggle, "fa-circle-stop");
    }
}
function resetCD() {
    if (!cd || !timerContainer) {
        console.error("resetCD :: cd and/or timerContainer are not set yet");
        return;
    }
    // timerContainer.removeEventListener("htmx:afterSettle", resetCD)
    // document.body.removeEventListener("htmx:load", bootstrap)
    bootstrap();
    console.log("cd.start()");
    cd.start();
}
document.body.addEventListener("htmx:load", bootstrap);
if (timerToggle) {
    timerToggle.addEventListener("click", toggleTimer);
}
if (timerContainer) {
    timerContainer.addEventListener("htmx:afterSettle", resetCD);
}
