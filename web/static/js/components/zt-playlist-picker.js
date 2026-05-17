(function() {
"use strict";
function escapeHtml(s) {
  return String(s).replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
function ZtPlaylistPicker() {
  var el = Reflect.construct(HTMLElement, [], ZtPlaylistPicker);
  return el;
}
ZtPlaylistPicker.prototype = Object.create(HTMLElement.prototype);
ZtPlaylistPicker.prototype.attributeChangedCallback = function(name) {
  if (name === "data-video-id" && this._load) this._load();
};
ZtPlaylistPicker.observedAttributes = ["data-video-id"];
ZtPlaylistPicker.prototype.disconnectedCallback = function() {
  if (this._docClick) {
    document.removeEventListener("click", this._docClick, true);
    this._docClick = null;
  }
  if (this._menuEl && this._menuEl.parentNode === document.body) {
    this._menuEl.remove();
  }
  this._ztPickerInit = false;
};
ZtPlaylistPicker.prototype.connectedCallback = function() {
  var self = this;
  var videoId = this.getAttribute("data-video-id");
  if (!videoId) return;
  if (self._ztPickerInit) return;
  self._ztPickerInit = true;
  self.innerHTML =
    '<div class="zt-playlist-picker-wrap d-inline-block position-relative">' +
    '<button type="button" class="btn btn-sm btn-outline-dark zt-playlist-picker-btn">' +
    '<i class="fas fa-list text-secondary"></i> Add to playlist</button>' +
    '<ul class="dropdown-menu dropdown-menu-end p-2 zt-playlist-picker-menu" style="min-width:260px;display:none">' +
    '<li class="px-2 py-1 text-muted small">Loading...</li></ul></div>';
  var btn = self.querySelector(".zt-playlist-picker-btn");
  var menu = self.querySelector(".zt-playlist-picker-menu");
  var wrap = self.querySelector(".zt-playlist-picker-wrap");
  self._menuEl = menu;

  function positionMenu() {
    if (!btn || !menu) return;
    var rect = btn.getBoundingClientRect();
    menu.style.position = "fixed";
    menu.style.top = Math.round(rect.bottom + 4) + "px";
    menu.style.right = Math.round(window.innerWidth - rect.right) + "px";
    menu.style.left = "auto";
    menu.style.bottom = "auto";
    menu.style.zIndex = "10000";
  }

  function attachMenu() {
    if (menu && menu.parentNode !== document.body) {
      document.body.appendChild(menu);
    }
  }

  function detachMenu() {
    if (menu && wrap && menu.parentNode === document.body) {
      wrap.appendChild(menu);
      menu.style.position = "";
      menu.style.top = "";
      menu.style.right = "";
      menu.style.left = "";
      menu.style.zIndex = "";
    }
  }

  function setOpen(open) {
    if (!menu) return;
    if (open) {
      attachMenu();
      menu.style.display = "block";
      menu.classList.add("show");
      positionMenu();
      if (btn) btn.setAttribute("aria-expanded", "true");
      load();
    } else {
      menu.style.display = "none";
      menu.classList.remove("show");
      if (btn) btn.setAttribute("aria-expanded", "false");
      detachMenu();
    }
  }

  if (btn) {
    btn.addEventListener("click", function(e) {
      e.preventDefault();
      e.stopPropagation();
      setOpen(menu.style.display !== "block");
    });
  }

  if (menu) {
    menu.addEventListener("click", function(e) {
      e.stopPropagation();
    });
  }

  self._docClick = function(e) {
    var t = e.target;
    var inside = (wrap && wrap.contains(t)) || (menu && menu.contains(t));
    if (inside) return;
    setOpen(false);
  };
  document.addEventListener("click", self._docClick, true);

  function renderList(playlists) {
    videoId = self.getAttribute("data-video-id");
    if (!videoId) return;
    var html = "";
    if (!playlists.length) {
      html += '<li class="px-2 py-1 text-muted small">No playlists yet.</li>';
    } else {
      playlists.filter(function(p) { return !p.virtual; }).forEach(function(p) {
        var checked = p.contains ? " checked" : "";
        html += '<li class="dropdown-item-text px-2 py-1">' +
          '<label class="d-flex align-items-center gap-2 mb-0" style="cursor:pointer">' +
          '<input type="checkbox" class="form-check-input zt-playlist-picker-toggle" data-id="' + escapeHtml(p.id) + '"' + checked + ">" +
          "<span>" + escapeHtml(p.name || "") + "</span></label></li>";
      });
    }
    html += '<li><hr class="dropdown-divider"></li>';
    html += '<li class="px-2"><form id="zt-playlist-picker-create" class="d-flex gap-1">' +
      '<input type="text" class="form-control form-control-sm" placeholder="New playlist" required>' +
      '<button type="submit" class="btn btn-sm btn-primary">+</button></form></li>';
    menu.innerHTML = html;

    menu.querySelectorAll(".zt-playlist-picker-toggle").forEach(function(cb) {
      cb.addEventListener("change", function() {
        var pid = cb.getAttribute("data-id");
        if (!pid) return;
        if (cb.checked) {
          fetch("/api/playlists/" + encodeURIComponent(pid) + "/videos", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "same-origin",
            body: JSON.stringify({ video_id: videoId })
          }).then(function(r) {
            if (!r.ok) { cb.checked = false; return; }
            if (typeof sendToast === "function") sendToast("Added", "", "bg-success", "Video added to playlist");
          }).catch(function() { cb.checked = false; });
        } else {
          fetch("/api/playlists/" + encodeURIComponent(pid) + "/videos/" + encodeURIComponent(videoId), {
            method: "DELETE",
            credentials: "same-origin"
          }).then(function(r) {
            if (!(r.status === 204 || r.ok)) cb.checked = true;
            else if (typeof sendToast === "function") sendToast("Removed", "", "bg-secondary", "Video removed from playlist");
          }).catch(function() { cb.checked = true; });
        }
      });
    });

    var createForm = menu.querySelector("#zt-playlist-picker-create");
    if (createForm) {
      createForm.addEventListener("submit", function(e) {
        e.preventDefault();
        var input = createForm.querySelector("input");
        var name = input.value.trim();
        if (!name) return;
        fetch("/api/playlists", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "same-origin",
          body: JSON.stringify({ name: name })
        }).then(function(r) { return r.json().then(function(data) { return { ok: r.ok, data: data }; }); })
          .then(function(res) {
            if (!res.ok || !res.data || !res.data.id) return;
            var pid = res.data.id;
            return fetch("/api/playlists/" + encodeURIComponent(pid) + "/videos", {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              credentials: "same-origin",
              body: JSON.stringify({ video_id: videoId })
            }).then(function() { load(); });
          });
      });
    }
  }

  function load() {
    videoId = self.getAttribute("data-video-id");
    if (!videoId || !menu) return;
    fetch("/api/playlists?video_id=" + encodeURIComponent(videoId), { credentials: "same-origin" })
      .then(function(r) {
        if (r.status === 401) return { playlists: [] };
        return r.json();
      })
      .then(function(data) {
        renderList((data && data.playlists) || []);
        if (menu.style.display === "block") positionMenu();
      })
      .catch(function() {
        menu.innerHTML = '<li class="px-2 py-1 text-danger small">Failed to load playlists</li>';
      });
  }

  self._load = load;
};
customElements.define("zt-playlist-picker", ZtPlaylistPicker);
})();
