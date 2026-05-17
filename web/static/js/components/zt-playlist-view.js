(function() {
"use strict";
function escapeHtml(s) {
  return String(s).replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
function ZtPlaylistView() {
  var el = Reflect.construct(HTMLElement, [], ZtPlaylistView);
  return el;
}
ZtPlaylistView.prototype = Object.create(HTMLElement.prototype);
ZtPlaylistView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) {
    self.innerHTML = '<div class="alert alert-danger">Missing playlist id.</div>';
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    return;
  }
  fetch("/api/playlists/" + encodeURIComponent(id), { credentials: "same-origin" })
    .then(function(r) {
      if (r.status === 401) {
        window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
        return Promise.reject();
      }
      if (!r.ok) throw new Error(String(r.status));
      return r.json();
    })
    .then(function(data) {
      var pl = data.playlist || {};
      var videos = data.videos || [];
      var name = escapeHtml(pl.name || "Playlist");
      var html = '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-list"></i></span>';
      html += '<h3 id="zt-playlist-title">' + name + '</h3>';
      html += '<p class="text-muted"><a href="/profile/playlists">Your playlists</a></p><hr /></div>';
      html += '<div class="mb-3 d-flex flex-wrap gap-2 align-items-center">';
      html += '<form id="zt-playlist-rename-form" class="d-flex gap-2" style="max-width:400px">';
      html += '<input type="text" class="form-control form-control-sm" id="zt-playlist-rename-input" value="' + name + '">';
      html += '<button type="submit" class="btn btn-sm btn-outline-primary">Rename</button></form>';
      html += '<button type="button" class="btn btn-primary btn-sm" id="zt-playlist-play-btn"' + (videos.length === 0 ? ' disabled' : '') + '><i class="fa fa-play"></i> Play playlist</button>';
      html += '<button type="button" class="btn btn-sm btn-outline-danger" id="zt-playlist-delete-btn">Delete playlist</button>';
      html += '<span id="zt-playlist-rename-error" class="text-danger small" style="display:none"></span></div>';
      html += '<div id="zt-playlist-videos" class="row">';
      if (videos.length === 0) {
        html += '<p class="text-muted">No videos in this playlist yet. Add videos from a video or clip page.</p>';
      } else {
        videos.forEach(function(v) {
          var vid = v.ID || v.id;
          html += '<div class="col-md-3 mb-4" data-video-id="' + escapeHtml(vid) + '">';
          html += '<zt-video-tile data-video="' + String(JSON.stringify(v)).replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;") + '"></zt-video-tile>';
          html += '<button type="button" class="btn btn-sm btn-outline-danger w-100 mt-1 zt-playlist-remove-video" data-video-id="' + escapeHtml(vid) + '">Remove</button>';
          html += '</div>';
        });
      }
      html += '</div>';
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);

      var playBtn = self.querySelector("#zt-playlist-play-btn");
      if (playBtn && videos.length > 0) {
        playBtn.addEventListener("click", function() {
          var url = window.ztPlaylistPlayUrl(videos[0], id, { autoplay: true });
          if (window.navigate) window.navigate(url);
          else window.location.href = url;
        });
      }

      self.querySelectorAll("#zt-playlist-videos zt-video-tile a").forEach(function(a) {
        a.addEventListener("click", function(e) {
          var col = a.closest("[data-video-id]");
          if (!col) return;
          var vid = col.getAttribute("data-video-id");
          var v = videos.find(function(x) { return (x.ID || x.id) === vid; });
          if (!v) return;
          e.preventDefault();
          var url = window.ztPlaylistPlayUrl(v, id, {});
          if (window.navigate) window.navigate(url);
          else window.location.href = url;
        });
      });

      self.querySelector("#zt-playlist-rename-form").addEventListener("submit", function(e) {
        e.preventDefault();
        var newName = self.querySelector("#zt-playlist-rename-input").value.trim();
        var errEl = self.querySelector("#zt-playlist-rename-error");
        errEl.style.display = "none";
        if (!newName) {
          errEl.textContent = "Name is required.";
          errEl.style.display = "inline";
          return;
        }
        fetch("/api/playlists/" + encodeURIComponent(id), {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          credentials: "same-origin",
          body: JSON.stringify({ name: newName })
        }).then(function(r) { return r.json().then(function(body) { return { ok: r.ok, body: body }; }); })
          .then(function(res) {
            if (res.ok) {
              self.querySelector("#zt-playlist-title").textContent = newName;
              if (typeof sendToast === "function") sendToast("Renamed", "", "bg-success", newName);
            } else {
              errEl.textContent = (res.body && res.body.error) || "Rename failed.";
              errEl.style.display = "inline";
            }
          });
      });

      self.querySelector("#zt-playlist-delete-btn").addEventListener("click", function() {
        if (!confirm("Delete this playlist?")) return;
        fetch("/api/playlists/" + encodeURIComponent(id), { method: "DELETE", credentials: "same-origin" })
          .then(function(r) {
            if (r.status === 204 || r.ok) {
              if (window.navigate) window.navigate("/profile/playlists");
              else window.location.href = "/profile/playlists";
            }
          });
      });

      self.querySelectorAll(".zt-playlist-remove-video").forEach(function(btn) {
        btn.addEventListener("click", function() {
          var videoId = btn.getAttribute("data-video-id");
          if (!videoId) return;
          fetch("/api/playlists/" + encodeURIComponent(id) + "/videos/" + encodeURIComponent(videoId), {
            method: "DELETE",
            credentials: "same-origin"
          }).then(function(r) {
            if (r.status === 204 || r.ok) {
              var col = btn.closest("[data-video-id]");
              if (col) col.remove();
              var grid = self.querySelector("#zt-playlist-videos");
              if (grid && !grid.querySelector("[data-video-id]")) {
                grid.innerHTML = '<p class="text-muted">No videos in this playlist yet. Add videos from a video or clip page.</p>';
              }
            }
          });
        });
      });

      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Playlist not found.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-playlist-view", ZtPlaylistView);
})();
