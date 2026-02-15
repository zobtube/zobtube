(function() {
"use strict";
var admHomeStyles = "h4:not(.first-h4){margin-top:30px}.s16{width:30px;height:16px}.card-adm{display:flex;justify-content:space-between}.card-adm-btn{display:inline;align-self:center}dt{float:left;clear:left;margin-right:30px;width:160px;text-align:right}dd{margin-left:0}";
function esc(s){ return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;"); }
function ZtAdmHome() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmHome);
  return el;
}
ZtAdmHome.prototype = Object.create(HTMLElement.prototype);
ZtAdmHome.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm", { credentials: "same-origin" })
    .then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
    .then(function(d) {
      var build = d.build || d.Build || {};
      var vc = d.video_count ?? d.VideoCount ?? 0;
      var ac = d.actor_count ?? d.ActorCount ?? 0;
      var cc = d.channel_count ?? d.ChannelCount ?? 0;
      var uc = d.user_count ?? d.UserCount ?? 0;
      var catc = d.category_count ?? d.CategoryCount ?? 0;
      var healthErrors = d.health_errors || d.HealthErrors || [];
      var version = build.Version || build.version || "?";
      var commit = build.Commit || build.commit || "none";
      var buildDate = build.BuildDate || build.build_date || "";
      var golang = d.golang_version || d.GolangVersion || "";
      var dbDriver = d.db_driver || d.DBDriver || "";
      var binaryPath = esc(d.binary_path || d.BinaryPath || "");
      var startupDir = esc(d.startup_directory || d.StartupDirectory || "");
      var commitShort = commit === "none" ? "none" : String(commit).substring(0, 6);

      var html = '<style>'+admHomeStyles+'</style><div class="row"><div class="col-12"><zt-adm-tabs data-active="overview"></zt-adm-tabs></div>';
      html += '<div class="col-md-12"><h4 class="first-h4">Overview</h4><hr/><div class="row">';
      html += '<div class="col-md-4 mb-3"><div class="card"><div class="card-body card-adm"><span><div class="d-flex align-items-center"><i class="fas fa-video s16"></i><h3 class="gl-m-0 gl-ml-3">'+vc+'</h3></div><div class="gl-mt-3 text-uppercase">Videos</div></span><div class="card-adm-btn"><a class="btn btn-primary" href="/adm/videos">View</a></div></div></div></div>';
      html += '<div class="col-md-4 mb-3"><div class="card"><div class="card-body card-adm"><span><div class="d-flex align-items-center"><i class="far fa-user s16"></i><h3 class="gl-m-0 gl-ml-3">'+ac+'</h3></div><div class="gl-mt-3 text-uppercase">Actors</div></span><div class="card-adm-btn"><a class="btn btn-primary" href="/adm/actors">View</a></div></div></div></div>';
      html += '<div class="col-md-4 mb-3"><div class="card"><div class="card-body card-adm"><span><div class="d-flex align-items-center"><i class="fas fa-podcast s16"></i><h3 class="gl-m-0 gl-ml-3">'+cc+'</h3></div><div class="gl-mt-3 text-uppercase">Channels</div></span><div class="card-adm-btn"><a class="btn btn-primary" href="/adm/channels">View</a></div></div></div></div>';
      html += '<div class="col-md-4 mb-3"><div class="card"><div class="card-body card-adm"><span><div class="d-flex align-items-center"><i class="fas fa-user-circle s16"></i><h3 class="gl-m-0 gl-ml-3">'+uc+'</h3></div><div class="gl-mt-3 text-uppercase">Users</div></span><div class="card-adm-btn"><a class="btn btn-primary" href="/adm/users">View</a></div></div></div></div>';
      html += '<div class="col-md-4 mb-3"><div class="card"><div class="card-body card-adm"><span><div class="d-flex align-items-center"><i class="fas fa-certificate s16"></i><h3 class="gl-m-0 gl-ml-3">'+catc+'</h3></div><div class="gl-mt-3 text-uppercase">Categories</div></span><div class="card-adm-btn"><a class="btn btn-primary" href="/adm/categories">View</a></div></div></div></div>';
      html += '</div></div>';

      html += '<div class="col-md-12"><h4>Health</h4><hr/>';
      if (healthErrors.length === 0) {
        html += '<div class="alert alert-success">No issues found in your configuration.</div>';
      } else {
        html += '<div class="alert alert-warning">Some errors were detected:<ul style="list-style:disc;margin-left:18px">';
        healthErrors.forEach(function(e){ html += '<li>'+esc(e)+'</li>'; });
        html += '</ul></div>';
      }
      html += '</div>';

      html += '<div class="col-md-12"><h4>About</h4><hr/><dl><dt>Version</dt><dd>'+esc(version)+' (commit <code>'+esc(commitShort)+'</code>)</dd><dt>Built at</dt><dd>'+esc(buildDate)+'</dd><dt>Golang</dt><dd>'+esc(golang)+'</dd><dt>Database</dt><dd>'+esc(dbDriver)+'</dd><dt>Binary Path</dt><dd>'+binaryPath+'</dd><dt>Working Path</dt><dd>'+startupDir+'</dd></dl></div>';

      html += '<div class="col-md-12"><h4>More Info</h4><hr/><dl><dt>Home Page</dt><dd><a href="https://zobtube.com">zobtube.com</a></dd><dt>Reddit</dt><dd><a href="https://www.reddit.com/r/zobtube/">www.reddit.com/r/zobtube</a></dd><dt>Source</dt><dd><a href="https://github.com/zobtube/zobtube">github.com/zobtube/zobtube</a></dd><dt>Feature Requests</dt><dd><a href="https://github.com/zobtube/zobtube/issues">github.com/zobtube/zobtube/issues</a></dd></dl></div>';
      html += '</div>';

      self.innerHTML = html;
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function(e) {
      if (e && e.message === "403") { self.innerHTML = '<div class="alert alert-danger">Forbidden</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
      self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-adm-home", ZtAdmHome);
})();
