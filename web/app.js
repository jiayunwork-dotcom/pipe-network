const $ = (id) => document.getElementById(id);
const COLORS = ["#2f6feb", "#e05d2a", "#1e9e6a", "#8b4fd6", "#c2a020", "#3a8fa8"];

async function fetchJSON(url, method, body) {
  const opt = { method: method || "GET" };
  if (body) {
    opt.headers = { "Content-Type": "application/json" };
    opt.body = JSON.stringify(body);
  }
  const r = await fetch(url, opt);
  return r;
}

function layout(nodes) {
  const pos = {};
  const cx = 300, cy = 180, R = 130;
  const sources = nodes.filter((n) => n.is_source);
  const others = nodes.filter((n) => !n.is_source);
  sources.forEach((n, i) => {
    pos[n.id] = { x: cx - R - 60, y: cy + (i - (sources.length - 1) / 2) * 60 };
  });
  others.forEach((n, i) => {
    const a = (Math.PI * 2 * i) / Math.max(1, others.length) - Math.PI / 2;
    pos[n.id] = { x: cx + R * Math.cos(a) * 0.7 + 40, y: cy + R * Math.sin(a) * 0.7 };
  });
  return pos;
}

function drawDiagram(net, result) {
  const svg = $("diagram");
  svg.innerHTML = "";
  const pos = layout(net.nodes);
  net.pipes.forEach((p) => {
    const a = pos[p.from], b = pos[p.to];
    if (!a || !b) return;
    const line = document.createElementNS("http://www.w3.org/2000/svg", "line");
    line.setAttribute("x1", a.x); line.setAttribute("y1", a.y);
    line.setAttribute("x2", b.x); line.setAttribute("y2", b.y);
    line.setAttribute("class", "edge");
    svg.appendChild(line);
    const q = result ? (result.flow[p.id] || 0) : 0;
    const t = document.createElementNS("http://www.w3.org/2000/svg", "text");
    t.setAttribute("x", (a.x + b.x) / 2);
    t.setAttribute("y", (a.y + b.y) / 2 - 4);
    t.setAttribute("class", "edge-label");
    t.setAttribute("text-anchor", "middle");
    t.textContent = p.id + " " + q.toFixed(4);
    svg.appendChild(t);
  });
  net.nodes.forEach((n) => {
    const p = pos[n.id];
    const c = document.createElementNS("http://www.w3.org/2000/svg", "circle");
    c.setAttribute("cx", p.x); c.setAttribute("cy", p.y); c.setAttribute("r", 14);
    c.setAttribute("class", n.is_source ? "node-dot node-source" : "node-dot");
    svg.appendChild(c);
    const t = document.createElementNS("http://www.w3.org/2000/svg", "text");
    t.setAttribute("x", p.x); t.setAttribute("y", p.y + 4);
    t.setAttribute("class", "node-label");
    t.setAttribute("text-anchor", "middle");
    t.setAttribute("fill", "#fff");
    t.textContent = n.id;
    svg.appendChild(t);
  });
}

function fillTables(net, result) {
  const ft = $("flows").querySelector("tbody");
  ft.innerHTML = "";
  net.pipes.forEach((p) => {
    const q = result.flow[p.id] || 0;
    const hf = (result.head[p.from] || 0) - (result.head[p.to] || 0);
    const tr = document.createElement("tr");
    tr.innerHTML = `<td>${p.id}</td><td>${q.toFixed(6)}</td><td>${hf.toFixed(4)}</td>`;
    ft.appendChild(tr);
  });
  const ht = $("heads").querySelector("tbody");
  ht.innerHTML = "";
  net.nodes.forEach((n) => {
    const h = result.head[n.id] || 0;
    const ph = h - (n.elevation || 0);
    const tr = document.createElement("tr");
    tr.innerHTML = `<td>${n.id}${n.is_source ? " (源)" : ""}</td><td>${h.toFixed(4)}</td><td>${ph.toFixed(4)}</td>`;
    ht.appendChild(tr);
  });
}

async function loadExample() {
  const name = $("example").value;
  if (!name) return;
  const r = await fetchJSON("/example/" + name);
  $("netlist").value = await r.text();
  $("error").classList.add("hidden");
}

async function solve() {
  const text = $("netlist").value;
  let net;
  try {
    net = JSON.parse(text);
  } catch (e) {
    showError("JSON 解析失败: " + e.message);
    return;
  }
  const r = await fetchJSON("/api/solve", "POST", net);
  const data = await r.json();
  if (!r.ok) {
    showError((data.message || "求解失败") + " [" + (data.error || "") + "]");
    return;
  }
  $("error").classList.add("hidden");
  drawDiagram(net, data);
  fillTables(net, data);
  $("status").textContent =
    `method=${data.method} iterations=${data.iterations} converged=${data.converged} ` +
    `maxMassResidual=${data.max_mass_residual?.toExponential(2)} loopClosure=${data.loop_closure?.toExponential(2)}`;
}

function showError(msg) {
  const el = $("error");
  el.textContent = msg;
  el.classList.remove("hidden");
}

(async function init() {
  const meta = await fetchJSON("/api/meta");
  const m = await meta.json();
  const sel = $("example");
  (m.examples || []).forEach((name) => {
    const o = document.createElement("option");
    o.value = name; o.textContent = name;
    sel.appendChild(o);
  });
  $("load").onclick = loadExample;
  $("solve").onclick = solve;
  if (sel.options.length) await loadExample();
})();
