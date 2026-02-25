(function() {
"use strict";
function ZtAdmTaskView() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmTaskView);
  return el;
}
ZtAdmTaskView.prototype = Object.create(HTMLElement.prototype);
ZtAdmTaskView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = '<div class="alert alert-danger">Missing id</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  fetch("/api/adm/task/" + encodeURIComponent(id), { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(t) {
      function esc(s) { return String(s||"").replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;"); }
      function badgeClass(s) {
        var v = (s||"").toLowerCase();
        if (v === "todo") return "secondary";
        if (v === "in-progress") return "primary";
        if (v === "done") return "success";
        if (v === "error") return "danger";
        return "secondary";
      }
      var status = t.Status || t.status || "";
      var bc = badgeClass(status);
      var createdAt = t.CreatedAt ? esc(t.CreatedAt) : "—";
      var updatedAt = t.UpdatedAt ? esc(t.UpdatedAt) : "—";
      var doneAt = t.DoneAt ? esc(t.DoneAt) : "Not done yet";
      var params = t.Parameters || {};
      var paramsHtml = Object.keys(params).map(function(k) { return '"' + esc(k) + '": ' + esc(params[k]); }).join("\n") || "—";

      var retryBtn = "";
      if (status.toLowerCase() === "error") {
        retryBtn = ' <button class="btn btn-warning btn-sm" style="float:right" type="button" data-zt-retry><i class="fas fa-sync-alt"></i> Restart task in error</button>';
      }

      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="tasks"></zt-adm-tabs></div><div class="col-md-9 col-lg-9">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="far fa-check-square"></i></span><h3>Task details</h3><hr /></div>';
      html += '<div class="row"><div class="col-md-12"><table class="table"><tbody>';
      html += '<tr><td>Task ID</td><td><code>' + esc(t.ID || t.id) + '</code></td></tr>';
      html += '<tr><td>Status</td><td><span class="badge text-bg-' + bc + '">' + esc(status) + '</span>' + retryBtn + '</td></tr>';
      html += '<tr><td>Task type</td><td>' + esc(t.Name || t.name) + '</td></tr>';
      html += '<tr><td>Step</td><td>' + esc(t.Step || t.step) + '</td></tr>';
      html += '<tr><td>Created at</td><td>' + createdAt + '</td></tr>';
      html += '<tr><td>Last update</td><td>' + updatedAt + '</td></tr>';
      html += '<tr><td>Done at</td><td>' + doneAt + '</td></tr>';
      html += '<tr><td>Parameters</td><td><pre>' + paramsHtml.replace(/\n/g, "<br>") + '</pre></td></tr>';
      html += '</tbody></table></div></div></div></div>';
      self.innerHTML = html;

      var retryEl = self.querySelector("[data-zt-retry]");
      if (retryEl) {
        retryEl.addEventListener("click", function() {
          fetch("/api/adm/task/" + id + "/retry", { method: "POST", credentials: "same-origin" })
            .then(function(r) { return r.json(); })
            .then(function(data) {
              if (data.redirect && typeof loadPage === "function") loadPage(data.redirect);
              else window.location.reload();
            });
        });
      }
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Not found.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-task-view", ZtAdmTaskView);
})();
