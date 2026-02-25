(function() {
"use strict";
function ZtAdmConfigOffline() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmConfigOffline);
  return el;
}
ZtAdmConfigOffline.prototype = Object.create(HTMLElement.prototype);
ZtAdmConfigOffline.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/config/offline", { credentials: "same-origin" })
    .then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
    .then(function(d) {
      var enabled = d.offline_mode === true;
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="offline"></zt-adm-tabs></div><div class="col-md-9 col-lg-9">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-power-off"></i></span><h3>Offline mode</h3><hr /></div>';
      html += '<div class="row"><div class="col-md-12 mb-4">';
      html += '<p>By design, ZobTube already makes almost no calls to the outside world.</p>';
      html += '<p>Only calls are towards providers to retrieve information and pictures.</p>';
      html += '<p>If you want to completely disable any external call, you can enable the <i>Offline Mode</i>.</p>';
      html += '<p>It will disable provider calls and any future calls that could come with later development.</p></div>';
      html += '<div class="col-md-12 mb-4"><h4 class="mb-4">Offline mode</h4>';
      html += 'The Offline Mode is currently: ';
      html += enabled ? '<span class="badge text-bg-success">Enabled</span>' : '<span class="badge text-bg-warning">Disabled</span>';
      html += '</div><div class="col-md-12 mb-4">';
      html += enabled ? '<button class="btn btn-xl btn-danger" id="zt-offline-disable">Disable Offline Mode</button>' : '<button class="btn btn-xl btn-danger" id="zt-offline-enable">Enable Offline Mode</button>';
      html += '</div></div>';
      html += '</div></div>';
      self.innerHTML = html;
      var btn = self.querySelector(enabled ? "#zt-offline-disable" : "#zt-offline-enable");
      if (btn) {
        btn.addEventListener("click", function() {
          var action = enabled ? "disable" : "enable";
          fetch("/api/adm/config/offline/" + action, { method: "GET", credentials: "same-origin" })
            .then(function(r) {
              if (r.ok) {
                if (typeof loadPage === "function") loadPage("/adm/config/offline");
                else window.location.reload();
              }
            });
        });
      }
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function(e) {
      if (e && e.message === "403") self.innerHTML = '<div class="alert alert-danger">Forbidden</div>';
      else self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-adm-config-offline", ZtAdmConfigOffline);
})();
