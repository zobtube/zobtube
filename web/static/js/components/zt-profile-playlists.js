(function() {
"use strict";
function escapeHtml(s) {
  return String(s).replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
function ZtProfilePlaylists() {
  var el = Reflect.construct(HTMLElement, [], ZtProfilePlaylists);
  return el;
}
ZtProfilePlaylists.prototype = Object.create(HTMLElement.prototype);
ZtProfilePlaylists.prototype.connectedCallback = function() {
  var self = this;
  var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-profile-tabs data-active="playlists"></zt-profile-tabs></div><div class="col-md-9 col-lg-9">';
  html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-list"></i></span><h3>Your playlists</h3><hr /></div>';
  html += '<p class="text-muted">Create playlists to save videos and clips. Only you can see your playlists.</p>';
  html += '<form id="zt-playlist-create-form" class="mb-4" style="max-width:400px"><div class="input-group"><input type="text" class="form-control" id="zt-playlist-name" placeholder="Playlist name" required><button type="submit" class="btn btn-primary">Create playlist</button></div><div id="zt-playlist-form-error" class="alert alert-danger mt-2" style="display:none" role="alert"></div></form>';
  html += '<div id="zt-playlists-list"></div></div>';
  self.innerHTML = html;
  var listEl = self.querySelector("#zt-playlists-list");
  var form = self.querySelector("#zt-playlist-create-form");
  var errEl = self.querySelector("#zt-playlist-form-error");
  function loadList() {
    fetch("/api/playlists", { credentials: "same-origin" })
      .then(function(r) {
        if (r.status === 401) {
          window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
          return;
        }
        return r.json();
      })
      .then(function(data) {
        if (!data) return;
        var playlists = data.playlists || [];
        if (playlists.length === 0) {
          listEl.innerHTML = '<p class="text-muted">No playlists yet. Create one above.</p>';
          return;
        }
        var table = '<table class="table table-striped"><thead><tr><th>Name</th><th>Videos</th><th>Updated</th><th></th></tr></thead><tbody>';
        playlists.forEach(function(p) {
          var name = escapeHtml(p.name || "");
          var id = p.id || "";
          var count = p.video_count != null ? p.video_count : 0;
          var updated = p.updated_at ? new Date(p.updated_at).toLocaleString() : "";
          table += '<tr><td><a href="/playlist/' + escapeHtml(id) + '">' + name + '</a></td><td>' + count + '</td><td>' + escapeHtml(updated) + '</td><td><button type="button" class="btn btn-sm btn-outline-danger zt-playlist-delete" data-id="' + escapeHtml(id) + '">Delete</button></td></tr>';
        });
        table += "</tbody></table>";
        listEl.innerHTML = table;
        listEl.querySelectorAll(".zt-playlist-delete").forEach(function(btn) {
          btn.addEventListener("click", function() {
            var pid = btn.getAttribute("data-id");
            if (!pid || !confirm("Delete this playlist?")) return;
            fetch("/api/playlists/" + encodeURIComponent(pid), { method: "DELETE", credentials: "same-origin" })
              .then(function(r) {
                if (r.status === 401) {
                  window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
                  return;
                }
                if (r.status === 204 || r.ok) loadList();
              });
          });
        });
      })
      .catch(function() {
        listEl.innerHTML = '<div class="alert alert-danger">Failed to load playlists.</div>';
        listEl.innerHTML = '<div class="alert alert-danger">Failed to load playlists.</div>';
      });
  }
  form.addEventListener("submit", function(e) {
    e.preventDefault();
    var name = self.querySelector("#zt-playlist-name").value.trim();
    errEl.style.display = "none";
    errEl.textContent = "";
    if (!name) {
      errEl.textContent = "Name is required.";
      errEl.style.display = "block";
      return;
    }
    fetch("/api/playlists", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "same-origin",
      body: JSON.stringify({ name: name })
    }).then(function(r) {
      if (r.status === 401) {
        window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
        return Promise.resolve();
      }
      return r.json().then(function(data) {
        if (r.ok) {
          form.reset();
          loadList();
          if (typeof sendToast === "function") sendToast("Playlist created", "", "bg-success", name);
        } else {
          errEl.textContent = (data && data.error) || "Failed to create playlist.";
          errEl.style.display = "block";
        }
      });
    }).catch(function() {
      errEl.textContent = "Request failed.";
      errEl.style.display = "block";
    });
  });
  loadList();
  if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
};
customElements.define("zt-profile-playlists", ZtProfilePlaylists);
})();