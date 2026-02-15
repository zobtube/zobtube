(function() {
"use strict";
function ZtActorList() {
  var el = Reflect.construct(HTMLElement, [], ZtActorList);
  return el;
}
ZtActorList.prototype = Object.create(HTMLElement.prototype);
ZtActorList.prototype.connectedCallback = function() {
  var self = this;
  var admin = (window.__USER__ && window.__USER__.admin);
  fetch("/api/actor", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var items = data.items || [];
      var html = '<style>.card-text{display:flex;align-content:center;justify-content:space-between;color:#606060}</style>' +
        '<div class="themeix-section-h"><span class="heading-icon"><i class="far fa-user"></i></span><h3>Actors' +
        (admin ? ' <a href="/actor/new"><i class="fas fa-plus-circle"></i></a>' : '') + '</h3><hr /></div>' +
        '<div class="row"><div class="col-md-12"><div class="form-floating mb-3"><input id="actor-filter" class="form-control"><label for="actor-filter">Looking for someone?</label></div></div></div>' +
        '<div class="row row-cols-1 row-cols-md-4 g-4">';
      items.forEach(function(a) {
        var urlView = "/actor/" + (a.ID || a.id);
        var urlThumb = "/api/actor/" + (a.ID || a.id) + "/thumb";
        var sexIcon = (a.Sex || a.sex) === "f" ? "fa-venus" : (a.Sex || a.sex) === "m" ? "fa-mars" : (a.Sex || a.sex) === "s" ? "fa-mars-and-venus" : "fa-person-circle-question";
        var vlen = (a.Videos || a.videos || []).length;
        var llen = (a.Links || a.links || []).length;
        var name = (a.Name || a.name || "").replace(/&/g,"&amp;").replace(/</g,"&lt;");
        html += '<div class="col"><div class="card"><img data-src="' + urlThumb + '" class="card-img-top lazy"><div class="card-body"><h5 class="card-title">' +
          '<a class="stretched-link" href="' + urlView + '">' + name + '</a><span style="position:absolute;right:15px;"><i class="fa ' + sexIcon + '"></i></span></h5>' +
          '<div class="card-text"><a><i class="fas fa-film"></i> ' + vlen + '</a><a><i class="fas fa-link"></i> ' + llen + '</a></div></div></div></div>';
      });
      html += "</div>";
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (document.getElementById("actor-filter")) {
        document.getElementById("actor-filter").oninput = function() {
          var q = (this.value || "").toLowerCase();
          self.querySelectorAll(".col").forEach(function(col) {
            var txt = (col.textContent || "").toLowerCase();
            col.style.display = q === "" || txt.indexOf(q) >= 0 ? "" : "none";
          });
        };
      }
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-actor-list", ZtActorList);
})();
