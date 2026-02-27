package spec

import (
	"encoding/json"
	"fmt"
	"strings"
)

// backoffPreset defines a named backoff configuration preset.
type backoffPreset struct {
	Label   string  `json:"label"`
	Jitter  string  `json:"jitter"`
	FirstMs int64   `json:"firstMs"`
	MaxMs   int64   `json:"maxMs"`
	Factor  float64 `json:"factor"`
}

var presets = []backoffPreset{
	{Label: "Standard", Jitter: "none", FirstMs: 1000, MaxMs: 5000, Factor: 2.0},
	{Label: "Aggressive", Jitter: "full", FirstMs: 500, MaxMs: 30000, Factor: 3.0},
	{Label: "Gentle", Jitter: "equal", FirstMs: 2000, MaxMs: 10000, Factor: 1.5},
}

// builderXData returns the Alpine x-data expression for the task spec builder.
func builderXData(agentsEndpoint string) string {
	presetsJSON, _ := json.Marshal(presets)
	return fmt.Sprintf(`{
  name: '', slot: '', kind_type: 'subprocess',
  timeout_ms: 30000, restart_type: 'never', interval_ms: 0,
  admission: 'dropIfRunning',
  backoff_preset: 'standard',
  jitter: 'none', backoff_first_ms: 1000, backoff_max_ms: 5000, backoff_factor: 2.0,

  cmd: '', args: '', env_rows: [], cwd: '', fail_on_non_zero: true,
  wasm_json: '{}', container_json: '{}',

  target_mode: 'agents',
  agents: [], agents_opts: [], agents_open: false,
  label_rows: [],

  runner_label_rows: [],

  presets: %s,
  submitting: false,
  agents_endpoint: '%s',

  get kindConfig() {
    if (this.kind_type === 'subprocess') {
      const cfg = {};
      if (this.cmd) cfg.command = this.cmd;
      const a = this.args.split(/\s+/).filter(Boolean);
      if (a.length) cfg.args = a;
      const envMap = {};
      for (const r of this.env_rows) { const k = r.key.trim(); if (k) envMap[k] = r.value; }
      if (Object.keys(envMap).length) cfg.env = envMap;
      if (this.cwd) cfg.cwd = this.cwd;
      cfg.failOnNonZero = this.fail_on_non_zero;
      return cfg;
    }
    if (this.kind_type === 'wasm') { try { return JSON.parse(this.wasm_json); } catch { return {}; } }
    if (this.kind_type === 'container') { try { return JSON.parse(this.container_json); } catch { return {}; } }
    return {};
  },

  get targetLabels() {
    const out = {};
    for (const r of this.label_rows) { const k = r.key.trim(); if (k) out[k] = r.value; }
    return out;
  },

  get runnerLabels() {
    const out = {};
    for (const r of this.runner_label_rows) { const k = r.key.trim(); if (k) out[k] = r.value; }
    return out;
  },

  get createSpec() {
    const spec = {
      name: this.name, slot: this.slot,
      kind_type: this.kind_type, kind_config: this.kindConfig,
      timeout_ms: Number(this.timeout_ms),
      restart_type: this.restart_type,
      jitter: this.jitter,
      backoff_first_ms: Number(this.backoff_first_ms),
      backoff_max_ms: Number(this.backoff_max_ms),
      backoff_factor: Number(this.backoff_factor),
      admission: this.admission,
    };
    if (this.restart_type === 'always' && this.interval_ms > 0) {
      spec.interval_ms = Number(this.interval_ms);
    }
    if (this.target_mode === 'agents' && this.agents.length) {
      spec.targets = this.agents;
    }
    const tl = this.targetLabels;
    if (Object.keys(tl).length) spec.target_labels = tl;
    const rl = this.runnerLabels;
    if (Object.keys(rl).length) spec.runner_labels = rl;
    return spec;
  },

  get previewJSON() {
    return JSON.stringify(this.createSpec, null, 2);
  },

  applyPreset(name) {
    const p = this.presets.find(x => x.label.toLowerCase() === name);
    if (p) {
      this.jitter = p.jitter;
      this.backoff_first_ms = p.firstMs;
      this.backoff_max_ms = p.maxMs;
      this.backoff_factor = p.factor;
    }
    this.backoff_preset = name;
  }
}`, string(presetsJSON), agentsEndpoint)
}

// builderInitExpr returns the Alpine x-init expression that loads active agents.
func builderInitExpr() string {
	return `fetch(agents_endpoint).then(r => r.json()).then(d => {
  agents_opts = (d.items || []).map(a => a.id);
}).catch(() => {})`
}

// builderSubmitExpr returns the Alpine submit expression for the builder form.
func builderSubmitExpr(action string) string {
	return fmt.Sprintf(
		`submitting = true;
fetch('%s', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json', 'HX-Request': 'true' },
  body: JSON.stringify(createSpec)
}).then(r => {
  if (r.ok) {
    const redirect = r.headers.get('HX-Redirect');
    if (redirect) { window.location.href = redirect; return; }
    show = false;
    htmx.trigger(document.body, 'taskspec_update');
  }
}).catch(() => {}).finally(() => submitting = false)`,
		strings.ReplaceAll(action, "'", "\\'"),
	)
}
