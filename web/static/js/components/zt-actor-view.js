(function() {
"use strict";
function esc(s) {
  return String(s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/"/g, "&quot;");
}
function ZtActorView() {
  var el = Reflect.construct(HTMLElement, [], ZtActorView);
  return el;
}
ZtActorView.prototype = Object.create(HTMLElement.prototype);

function photosetCardsHtml(items) {
  if (!items || items.length === 0) {
    return '<div class="col-md-12"><div class="alert alert-warning">No photosets for this actor.</div></div>';
  }
  var html = "";
  items.forEach(function(ps) {
    var pid = ps.ID || ps.id;
    var name = esc(ps.Name || ps.name || "Untitled");
    var status = ps.Status || ps.status || "";
    var cover = "/api/photoset/" + pid + "/cover";
    html += '<div class="col-md-3 col-sm-6 mb-4"><a href="/photoset/' + pid + '" class="text-decoration-none text-dark">';
    html += '<div class="card h-100"><div class="ratio ratio-4x3 bg-light">';
    html += '<img class="lazy" data-src="' + cover + '" alt="' + name + '" style="object-fit:cover;width:100%;height:100%">';
    html += '</div></div><div class="card-body p-2"><h6 class="card-title mb-0">' + name + '</h6>';
    if (status && status !== "ready") html += '<small class="text-muted">' + esc(status) + "</small>";
    html += "</div></div></a></div>";
  });
  return html;
}

ZtActorView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) {
    self.innerHTML = "Missing id";
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    return;
  }
  var actorUrl = "/api/actor/" + encodeURIComponent(id);
  var photosetsUrl = actorUrl + "/photosets";
  Promise.all([
    fetch(actorUrl, { credentials: "same-origin" }).then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    }),
    fetch(photosetsUrl, { credentials: "same-origin" }).then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
  ])
    .then(function(results) {
      var a = results[0];
      var psData = results[1];
      var photosets = psData.items || [];
      var photosetTotal = psData.total != null ? psData.total : photosets.length;
      var admin = window.__USER__ && window.__USER__.admin;
      var aliases = (a.Aliases || a.aliases || []).map(function(x) {
        return x.Name || x.name;
      }).filter(Boolean).join(" / ");
      var catsHtml = (a.Categories || a.categories || []).map(function(s) {
        var n = esc(s.Name || s.name || "");
        return '<div class="category">' + n + "</div>";
      }).join("");
      var links = a.Links || a.links || [];
      var videos = a.Videos || a.videos || [];
      var videoCount = videos.length;
      var name = esc(a.Name || a.name || "");
      var desc = esc(a.Description || a.description || "").replace(/\n/g, "<br>");
      var linksHtml = links.map(function(l) {
        return '<a href="' + (l.URL || l.url) + '" target="_blank" rel="noopener noreferrer"><img class="img-rounded" src="/static/images/provider-' + (l.Provider || l.provider) + '.png" style="height:80px;width:80px;margin-top:5px"></a>';
      }).join("");
      var initialTab = "videos";
      try {
        if (new URLSearchParams(window.location.search).get("tab") === "photosets") initialTab = "photosets";
      } catch (e) { /* ignore */ }

      var html = '<style>';
      html += ".actor_name{font-size:3rem}";
      html += ".bio_detail_label{color:#6b6b6b}";
      html += ".zt-actor-tabs{display:flex;gap:0;border-bottom:1px solid #ddd;background:#f0f0f0;margin:0 0 1.5rem;padding:0}";
      html += ".zt-actor-tab{display:flex;align-items:center;gap:.5rem;padding:.75rem 1.25rem;border:none;background:transparent;cursor:pointer;font-weight:600;color:#252525;position:relative}";
      html += ".zt-actor-tab:hover{color:#167ac6}";
      html += ".zt-actor-tab.active::after{content:'';position:absolute;left:0;right:0;bottom:0;height:4px;background:#f44336}";
      html += ".zt-actor-tab-badge{background:#ff9800;color:#fff;font-size:.75rem;font-weight:600;padding:.15rem .45rem;border-radius:4px;min-width:1.5rem;text-align:center}";
      html += "</style>";
      html += '<div style="display:flex"><div style="width:250px;margin-right:25px"><img class="img-rounded" src="/api/actor/' + id + '/thumb" style="height:250px;width:250px"></div>';
      html += '<div id="bio" style="flex-grow:1;margin:0"><h2 class="card-title actor_name">' + name + "</h2>";
      if (aliases) html += "<h4>aka " + esc(aliases) + "</h4>";
      if (admin) html += ' <a href="/actor/' + id + '/edit"><i>Edit profile</i></a>';
      html += '<div class="categories mb-4" style="padding-top:25px">' + catsHtml + "</div><div>" + desc + "</div></div>";
      html += '<div style="display:flex"><div style="width:170px"><div style="margin-top:20px;float:right">' + linksHtml + "</div></div></div></div>";
      html += "<hr /><br />";
      html += '<div class="zt-actor-tabs" role="tablist">';
      html += '<button type="button" class="zt-actor-tab' + (initialTab === "videos" ? " active" : "") + '" data-tab="videos" role="tab" aria-selected="' + (initialTab === "videos" ? "true" : "false") + '">';
      html += '<i class="fa fa-film"></i> Videos <span class="zt-actor-tab-badge">' + videoCount + "</span></button>";
      html += '<button type="button" class="zt-actor-tab' + (initialTab === "photosets" ? " active" : "") + '" data-tab="photosets" role="tab" aria-selected="' + (initialTab === "photosets" ? "true" : "false") + '">';
      html += '<i class="fa fa-images"></i> Photosets <span class="zt-actor-tab-badge">' + photosetTotal + "</span></button>";
      html += "</div>";
      html += '<div id="zt-actor-tab-videos" class="zt-actor-panel"' + (initialTab !== "videos" ? ' style="display:none"' : "") + ">";
      html += '<div class="row row-cols-1 row-cols-md-6 g-4">';
      videos.forEach(function(v) {
        html += (window.ztThumbPreviewHtml || function() { return ""; })(v);
      });
      html += "</div></div>";
      html += '<div id="zt-actor-tab-photosets" class="zt-actor-panel"' + (initialTab !== "photosets" ? ' style="display:none"' : "") + '><div class="row">';
      html += photosetCardsHtml(photosets);
      html += "</div></div>";

      self.innerHTML = html;

      var photosetsLazyDone = initialTab === "photosets";
      function lazyPhotosetsIfNeeded() {
        if (photosetsLazyDone) return;
        photosetsLazyDone = true;
        if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self.querySelector("#zt-actor-tab-photosets") || self);
      }
      if (initialTab === "photosets") lazyPhotosetsIfNeeded();

      function setTab(tab) {
        var tabs = self.querySelectorAll(".zt-actor-tab");
        var panels = {
          videos: self.querySelector("#zt-actor-tab-videos"),
          photosets: self.querySelector("#zt-actor-tab-photosets")
        };
        tabs.forEach(function(btn) {
          var on = btn.getAttribute("data-tab") === tab;
          btn.classList.toggle("active", on);
          btn.setAttribute("aria-selected", on ? "true" : "false");
        });
        Object.keys(panels).forEach(function(key) {
          if (panels[key]) panels[key].style.display = key === tab ? "" : "none";
        });
        if (tab === "photosets") lazyPhotosetsIfNeeded();
        var url = new URL(window.location.href);
        if (tab === "photosets") url.searchParams.set("tab", "photosets");
        else url.searchParams.delete("tab");
        window.history.replaceState(null, "", url.pathname + url.search + url.hash);
      }

      self.querySelectorAll(".zt-actor-tab").forEach(function(btn) {
        btn.addEventListener("click", function() {
          setTab(btn.getAttribute("data-tab"));
        });
      });

      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Not found.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-actor-view", ZtActorView);
})();
