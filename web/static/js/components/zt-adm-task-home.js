(function() {
"use strict";
function ZtAdmTaskHome() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmTaskHome);
  return el;
}
ZtAdmTaskHome.prototype = Object.create(HTMLElement.prototype);
ZtAdmTaskHome.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/task/home", { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(d) {
      var items = d.items || [];
      function esc(s) { return String(s||"").replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;"); }
      function badgeClass(s) {
        var v = (s||"").toLowerCase();
        if (v === "todo") return "secondary";
        if (v === "in-progress") return "primary";
        if (v === "done") return "success";
        if (v === "error") return "danger";
        return "secondary";
      }
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="tasks"></zt-adm-tabs></div><div class="col-md-9 col-lg-9"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-tasks"></i></span><h3>Task list</h3> <a href="/adm/tasks" style="margin-left:0;font-size:14px">View all tasks →</a><hr /></div><table class="table table-striped"><thead><tr><th>Task ID</th><th>Task</th><th>Status</th></tr></thead><tbody>';
      items.forEach(function(t) {
        var id = t.ID || t.id;
        var name = esc(t.Name || t.name || "");
        var status = esc(t.Status || t.status || "");
        var bc = badgeClass(t.Status || t.status);
        html += '<tr><td><a href="/adm/task/' + esc(id) + '"><code>' + esc(id) + '</code></a></td><td>' + name + '</td><td><span class="badge text-bg-' + bc + '">' + status + '</span></td></tr>';
      });
      html += '</tbody></table></div></div>';
      self.innerHTML = html;
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-task-home", ZtAdmTaskHome);
})();
