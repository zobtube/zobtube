(function() {
"use strict";
var profileStatsStyles = ".s16{width:30px;height:16px}.card-stat .card-body{padding:1.25rem}";
function formatViewTime(ns) {
  if (!ns) return "0:00";
  var nice = window.ztNiceDuration || function(n) {
    var s = Math.floor(n / 1e9), m = Math.floor(s / 60); s %= 60;
    var h = Math.floor(m / 60); m %= 60;
    return h > 0 ? String(h).padStart(2, "0") + ":" + String(m).padStart(2, "0") + ":" + String(s).padStart(2, "0")
      : String(m).padStart(2, "0") + ":" + String(s).padStart(2, "0");
  };
  var sec = Math.floor(ns / 1e9);
  if (sec < 86400) return nice(ns);
  var d = Math.floor(sec / 86400);
  sec %= 86400;
  var h = Math.floor(sec / 3600);
  sec %= 3600;
  var m = Math.floor(sec / 60);
  return d + "d " + h + "h " + m + "m";
}
function statCard(icon, value, label) {
  return '<div class="col-md-4 col-lg-4 mb-3"><div class="card card-stat"><div class="card-body">' +
    '<div class="d-flex align-items-center"><i class="' + icon + ' s16"></i><h3 class="gl-m-0 gl-ml-3 mb-0">' + value + '</h3></div>' +
    '<div class="gl-mt-3 text-uppercase text-muted small">' + label + '</div>' +
    '</div></div></div>';
}
function ZtProfileStats() {
  var el = Reflect.construct(HTMLElement, [], ZtProfileStats);
  return el;
}
ZtProfileStats.prototype = Object.create(HTMLElement.prototype);
ZtProfileStats.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/profile", { credentials: "same-origin" })
    .then(function(r) {
      if (!r.ok) throw new Error(String(r.status));
      return r.json();
    })
    .then(function(data) {
      var s = data.stats || data.Stats || {};
      var videosUnique = s.videos_unique ?? s.VideosUnique ?? 0;
      var videosTotal = s.videos_total ?? s.VideosTotal ?? 0;
      var actorsUnique = s.actors_unique ?? s.ActorsUnique ?? 0;
      var actorsTotal = s.actors_total ?? s.ActorsTotal ?? 0;
      var viewTimeNs = s.total_view_time_ns ?? s.TotalViewTimeNs ?? 0;
      var viewTime = formatViewTime(viewTimeNs);

      var html = '<style>' + profileStatsStyles + '</style><div class="row"><div class="col-md-3 col-lg-3"><zt-profile-tabs data-active="stats"></zt-profile-tabs></div><div class="col-md-9 col-lg-9">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-chart-bar"></i></span><h3>Your stats</h3><hr /></div>';
      html += '<div class="row">';
      html += statCard("fas fa-video", videosUnique, "Unique videos");
      html += statCard("fas fa-play-circle", videosTotal, "Total video views");
      html += statCard("fas fa-clock", viewTime, "Total view time");
      html += statCard("far fa-user", actorsUnique, "Unique actors");
      html += statCard("fas fa-users", actorsTotal, "Total actor views");
      html += '</div></div></div>';
      self.innerHTML = html;
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function(e) {
      if (e && e.message === "401") {
        window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
        return;
      }
      self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-profile-stats", ZtProfileStats);
})();
