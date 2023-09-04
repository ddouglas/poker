"use strict";
(() => {
  // src/countdown.ts
  var Countdown = class {
    constructor({ initialValue, showHour, emitter, onComplete }) {
      this.isRunning = false;
      this.initialValue = initialValue;
      this.showHour = showHour;
      this.emitter = emitter;
      this.onComplete = onComplete;
      this.countdownValue = initialValue;
      this.interval = null;
    }
    decrementCountdown() {
      console.debug("Countdown.decrementCountdown() start");
      this.countdownValue--;
      this.emitter(this.format(this.countdownValue));
      if (this.countdownValue === 0) {
        this.stop();
        this.onComplete();
      }
      console.debug("Countdown.decrementCountdown() stop");
    }
    format(duration) {
      console.debug("Countdown.format() start");
      const parts = {
        hours: Math.floor(duration / (60 * 60) % 24),
        minutes: Math.floor(duration / 60 % 60),
        seconds: Math.floor(duration % 60)
      };
      const bits = [];
      if (parts.hours > 0 && this.showHour) {
        let bit = "";
        if (parts.hours < 10) {
          bit = `0`;
        }
        bit = `${bit}${parts.hours}`;
      }
      if (parts.minutes == 0) {
        bits.push(`00`);
      } else if (parts.minutes > 0) {
        let bit = "";
        if (parts.minutes < 10) {
          bit = `0`;
        }
        bit = `${bit}${parts.minutes}`;
        bits.push(bit);
      }
      if (parts.seconds == 0) {
        bits.push(`00`);
      } else if (parts.seconds > 0) {
        let bit = "";
        if (parts.seconds < 10) {
          bit = `0`;
        }
        bit = `${bit}${parts.seconds}`;
        bits.push(bit);
      }
      console.debug("Countdown.format() stop");
      return bits.join(":");
    }
    start() {
      console.debug("Countdown.start()");
      if (this.interval) {
        console.debug("Countdown.start() if this.interval");
        this.reset();
      }
      if (this.countdownValue === 0) {
        console.error("countdownValue is 0, shutdown down");
        this.stop();
        return;
      }
      console.debug("Countdown.start() this.interval");
      this.interval = setInterval(() => this.decrementCountdown(), 1e3);
      this.isRunning = true;
      console.debug("Countdown.start() done");
    }
    continue() {
      this.stop();
      this.start();
    }
    stop() {
      console.debug("Countdown.stop() start");
      if (this.interval) {
        clearInterval(this.interval);
      }
      this.interval = null;
      this.isRunning = false;
      console.debug("Countdown.stop() stop");
    }
    toggle() {
      console.debug("Countdown.toggle() start");
      this.isRunning ? this.stop() : this.start();
      console.debug("Countdown.toggle() stop");
      return;
    }
    getIsRunning() {
      return this.isRunning;
    }
    reset() {
      console.debug("Countdown.reset() start");
      this.stop();
      this.countdownValue = this.initialValue;
      console.log("reset :: ", this);
      console.debug("Countdown.reset() stop");
    }
  };
  var countdown_default = Countdown;

  // src/elements.ts
  function fetchElements() {
    const timerContainer = document.getElementById("timer-container");
    const timer = document.getElementById("timer");
    const timerToggle = document.getElementById("toggle-timer-button");
    const nextTimerButton = document.getElementById("trigger-next-timer-level");
    if (!timerContainer) {
      console.error("Failed to fetch timer-container by id");
      return null;
    }
    if (!timer) {
      console.error("failed to fetch timer element by id");
      return null;
    }
    if (!timerToggle) {
      console.error("failed to fetch toggle-timer-button element by id");
      return null;
    }
    if (!nextTimerButton) {
      console.error("failed to fetch trigger-next-timer-level element by id");
      return null;
    }
    const nextLevelURI = `${nextTimerButton.getAttribute("hx-get")}?continue=true`;
    if (!nextLevelURI) {
      console.error("trigger-next-timer-level element is missing attribute hx-get");
      return null;
    }
    let durationSecStr = timer.getAttribute("data-level-duration-sec");
    if (!durationSecStr) {
      console.error("trigger-next-timer-level element is missing attribute data-level-duration-sec");
      durationSecStr = "0";
    }
    return { timer, timerContainer, timerToggle, nextTimerButton, nextLevelURI, durationSecStr };
  }

  // src/events.ts
  var abort;
  function initAbort() {
    abort = new AbortController();
  }
  initAbort();
  document.body.addEventListener("htmx:load", () => {
    initCountdown();
    initTimerToggleEventClick();
  });
  document.body.addEventListener("countdown::proceed", () => {
    console.debug("countdown::proceed :: start");
    resetCountdown();
    startCountdown();
    console.debug("countdown::proceed :: complete");
  });
  document.body.addEventListener("countdown::reset", () => {
    console.debug("countdown::reset :: start");
    resetCountdown();
    console.debug("countdown::reset :: complete");
  });
  function resetCountdown() {
    stopCountdown();
    abort.abort();
    initAbort();
    initCountdown();
    initTimerToggleEventClick();
  }
  function initTimerToggleEventClick() {
    const elements = fetchElements();
    if (!elements) {
      console.error("failed to fetch elements, unable to register click event on timer toggle");
      return;
    }
    const { timerToggle } = elements;
    timerToggle.addEventListener("click", () => toggleCountdown(), { signal: abort.signal });
  }

  // src/main.ts
  var countdown;
  function initCountdown() {
    console.debug("initCountdown :: start");
    const elements = fetchElements();
    if (!elements)
      return;
    const {
      // Endpoint that HTMX will use to reach out and fetch the next level
      nextLevelURI,
      // A String representation of the number of seconds that the time needs to run for
      durationSecStr,
      // The HTMLElement representing the text of our timer
      timer
      // // The HTMLElement representing the button that is used to start and stop the timer
      // timerToggle
    } = elements;
    let parsedDuractionSec = 0;
    if (durationSecStr) {
      parsedDuractionSec = parseInt(durationSecStr);
    }
    countdown = new countdown_default({
      initialValue: parsedDuractionSec,
      showHour: parsedDuractionSec > 3600,
      emitter: (text) => {
        timer.innerHTML = text;
        console.debug(`received emitted value ${text}`);
      },
      onComplete: () => {
        const nextLevelURIContinue = `${nextLevelURI}?continue=true`;
        htmx.ajax(
          "GET",
          nextLevelURIContinue,
          htmx.find("#timer-container")
        );
      }
    });
    console.debug("initCountdown :: complete");
  }
  function toggleCountdown() {
    console.debug("toggleCountdown :: start");
    if (!countdown) {
      console.error("timerToggleEventClick :: countdown not set");
      return;
    }
    const elements = fetchElements();
    if (!elements) {
      console.error("failed to fetch elements, unable to register click event on timer toggle");
      return;
    }
    const { timerToggle } = elements;
    countdown.toggle();
    if (!countdown.getIsRunning()) {
      htmx.removeClass(timerToggle, "fa-circle-stop");
      htmx.addClass(timerToggle, "fa-circle-play");
    } else {
      htmx.removeClass(timerToggle, "fa-circle-play");
      htmx.addClass(timerToggle, "fa-circle-stop");
    }
    console.debug("toggleCountdown :: complete");
  }
  function stopCountdown() {
    countdown?.stop();
    const elements = fetchElements();
    if (!elements) {
      console.error("failed to fetch elements, unable to register click event on timer toggle");
      return;
    }
    const { timerToggle } = elements;
    htmx.removeClass(timerToggle, "fa-circle-stop");
    htmx.addClass(timerToggle, "fa-circle-play");
  }
  function startCountdown() {
    const elements = fetchElements();
    if (!elements) {
      console.error("failed to fetch elements, unable to register click event on timer toggle");
      return;
    }
    const { timerToggle } = elements;
    htmx.removeClass(timerToggle, "fa-circle-stop");
    htmx.addClass(timerToggle, "fa-circle-play");
    countdown?.start();
  }
})();
//# sourceMappingURL=main.js.map
