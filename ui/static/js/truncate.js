document.addEventListener("alpine:init", () => {
    Alpine.data("truncate", () => ({
        full: "",
        init() {
            this.full = this.$el.dataset.full || "";

            const display = this.$refs.display;
            const measure = this.$refs.measure;

            const sync = () => {
                const cs = getComputedStyle(display);
                measure.style.font = cs.font;
                measure.style.letterSpacing = cs.letterSpacing;
            };

            const update = () => {
                sync();

                const avail = display.clientWidth;
                if (!avail) return;

                measure.textContent = this.full;
                if (measure.offsetWidth <= avail) {
                    display.textContent = this.full;
                    return;
                }

                const len = this.full.length;
                let lo = 1, hi = Math.floor(len / 2);
                while (lo < hi) {
                    const mid = Math.ceil((lo + hi) / 2);
                    measure.textContent =
                        this.full.slice(0, mid) + "\u2026" + this.full.slice(-mid);
                    if (measure.offsetWidth <= avail) {
                        lo = mid;
                    } else {
                        hi = mid - 1;
                    }
                }
                display.textContent =
                    this.full.slice(0, lo) + "\u2026" + this.full.slice(-lo);
            };

            const ro = new ResizeObserver(update);
            ro.observe(this.$el);
            update();
        },
    }));
});
