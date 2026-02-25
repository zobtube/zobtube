(function() {
"use strict";
function esc(s) { return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;"); }
function ZtAdmConfigProvider() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmConfigProvider);
  return el;
}
ZtAdmConfigProvider.prototype = Object.create(HTMLElement.prototype);
ZtAdmConfigProvider.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/config/provider", { credentials: "same-origin" })
    .then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
    .then(function(d) {
      var providers = d.providers || [];
      var loaded = d.provider_loaded || {};
      var offline = d.offline_mode === true;
      var html = '<style>.center{text-align:center}</style><div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="providers"></zt-adm-tabs></div><div class="col-md-9 col-lg-9">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-plug"></i></span><h3>Providers</h3><hr /></div>';
      html += '<div class="row"><div class="col-md-12 mb-4"><p>Providers are used to retrieve information and pictures of actors</p></div>';
      html += '<div class="col-md-12 mb-4"><table class="table"><thead><tr>';
      html += '<th scope="col">Provider</th><th class="center" scope="col">Status</th><th class="center" scope="col">Loaded</th>';
      html += '<th class="center" scope="col">Able to search actor</th><th class="center" scope="col">Able to scrape actor\'s picture</th></tr></thead><tbody>';
      providers.forEach(function(p) {
        var id = p.ID || p.id;
        var enabled = p.Enabled || p.enabled;
        var canSearch = p.AbleToSearchActor || p.ableToSearchActor;
        var canScrape = p.AbleToScrapePicture || p.ableToScrapePicture;
        var loadedName = loaded[id];
        html += '<tr><td>' + esc(id) + '</td><td class="center">';
        html += enabled ? '<span class="badge text-bg-success">Enabled</span>' : '<span class="badge text-bg-warning">Disabled</span>';
        html += ' <a href="/api/adm/config/provider/' + encodeURIComponent(id) + '/switch" class="zt-provider-switch" data-id="' + esc(id) + '"><i class="fas fa-sync"></i></a></td><td>';
        html += loadedName ? '<span class="badge text-bg-success">' + esc(loadedName) + ' loaded</span>' : '<span class="badge text-bg-danger">Not loaded</span>';
        html += '</td><td class="center">';
        if (canSearch) html += offline ? '<span class="badge text-bg-danger">Yes but disabled by Offline Mode</span>' : '<span class="badge text-bg-success">Yes</span>';
        else html += '<span class="badge text-bg-warning">No</span>';
        html += '</td><td class="center">';
        if (canScrape) html += offline ? '<span class="badge text-bg-danger">Yes but disabled by Offline Mode</span>' : '<span class="badge text-bg-success">Yes</span>';
        else html += '<span class="badge text-bg-warning">No</span>';
        html += '</td></tr>';
      });
      html += '</tbody></table></div></div>';
      html += '</div></div>';
      self.innerHTML = html;
      self.querySelectorAll(".zt-provider-switch").forEach(function(a) {
        a.addEventListener("click", function(e) {
          e.preventDefault();
          var id = a.getAttribute("data-id");
          fetch("/api/adm/config/provider/" + encodeURIComponent(id) + "/switch", { method: "GET", credentials: "same-origin" })
            .then(function(r) {
              if (r.ok && typeof loadPage === "function") loadPage("/adm/config/provider");
              else if (r.ok) window.location.reload();
            });
        });
      });
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function(e) {
      if (e && e.message === "403") self.innerHTML = '<div class="alert alert-danger">Forbidden</div>';
      else self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-adm-config-provider", ZtAdmConfigProvider);
})();
