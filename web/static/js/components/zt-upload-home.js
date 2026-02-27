(function() {
"use strict";

function humanFileSize(bytes, si, dp) {
  if (bytes === undefined || bytes === null) return "—";
  var thresh = si ? 1000 : 1024;
  if (Math.abs(bytes) < thresh) return bytes + " B";
  var units = si ? ["kB","MB","GB","TB"] : ["KiB","MiB","GiB","TiB"];
  var u = -1;
  var r = Math.pow(10, dp || 1);
  do {
    bytes /= thresh;
    u++;
  } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);
  return bytes.toFixed(dp || 1) + " " + units[u];
}

function ZtUploadHome() {
  var el = Reflect.construct(HTMLElement, [], ZtUploadHome);
  return el;
}
ZtUploadHome.prototype = Object.create(HTMLElement.prototype);
ZtUploadHome.prototype.connectedCallback = function() {
  var self = this;
  var currentPath = "/";
  var detailOffcanvas = null;
  var massActionModal = null;
  var massImportModal = null;
  var videoImportModal = null;
  var libraries = [];
  var selectedLibraryId = "";

  fetch("/api/adm/libraries", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(d) {
      libraries = d.items || [];
      var defaultLib = libraries.find(function(l) { return l.is_default || l.IsDefault; });
      selectedLibraryId = defaultLib ? (defaultLib.id || defaultLib.ID) : (libraries[0] && (libraries[0].id || libraries[0].ID)) || "";
      try {
        var saved = sessionStorage.getItem("zt-upload-library");
        if (saved && libraries.some(function(l) { return (l.id || l.ID) === saved; })) selectedLibraryId = saved;
      } catch (e) {}
      return Promise.resolve();
    })
    .catch(function() { libraries = []; selectedLibraryId = ""; return Promise.resolve(); })
    .then(function() {
  function updateView() {
    var sidebarStyles = "#zt-upload-sidebar{position:sticky;top:70px;padding:1rem 0;background:#f7f7f7;border-radius:8px;min-width:180px}#zt-upload-sidebar .zt-upload-sidebar-title{margin:0 0 0.75rem;padding:0 1rem;font-size:1.1rem;font-weight:600}#zt-upload-sidebar .nav{flex-direction:column}#zt-upload-sidebar .nav-link{padding:0.5rem 1rem;border-radius:4px;border-left:3px solid transparent;cursor:pointer;color:inherit;text-decoration:none;display:block}#zt-upload-sidebar .nav-link:hover{background:rgba(0,0,0,0.05)}#zt-upload-sidebar .nav-link.active{background:rgba(22,122,198,0.1);border-left-color:#167ac6;color:#167ac6}";
    var styles = ".br-link{color:#0d6efd;text-decoration:underline;cursor:pointer}.triage-listing-table{width:100%;table-layout:fixed;font-size:medium}.triage-listing-table tr{user-select:none}.triage-listing-table tr:hover td{background-color:#efefef}.checks{color:grey;cursor:pointer}";
    var html = '<style>'+sidebarStyles+'</style><style>'+styles+'</style>';
    html += '<div class="zt-upload-wrap" style="display:flex;gap:1rem;align-items:flex-start">';
    html += '<div id="zt-upload-sidebar" class="zt-upload-sidebar">';
    html += '<h3 class="zt-upload-sidebar-title"><i class="fa fa-folder-open"></i> Library</h3><nav class="nav flex-column">';
    if (libraries.length === 0) {
      html += '<span class="text-muted small px-2">No libraries</span>';
    } else {
      libraries.forEach(function(lib) {
        var id = lib.id || lib.ID;
        var name = (lib.name || lib.Name || id).replace(/&/g,"&amp;").replace(/</g,"&lt;");
        var isDefault = !!(lib.is_default || lib.IsDefault);
        var badge = isDefault ? ' <span class="badge bg-primary">default</span>' : '';
        var active = id === selectedLibraryId ? " nav-link active" : " nav-link";
        html += '<a class="' + active + '" href="#" data-library-id="' + id.replace(/"/g,"&quot;") + '">' + name + badge + '</a>';
      });
    }
    html += '</nav></div>';
    html += '<div class="zt-upload-main" style="flex:1;min-width:0">';
    html += '<div style="display:flex;justify-content:space-between;"><h4>Upload and triage folder</h4><div>';
    html += '<button class="btn btn-outline-success" id="zt-upload-file-btn">Upload file</button>';
    html += ' <button class="btn btn-outline-success" data-bs-toggle="modal" data-bs-target="#newFolderModal">New folder</button>';
    html += ' <button class="btn btn-outline-success disabled" id="zt-mass-action-btn">Mass action</button>';
    html += '</div></div><hr />';
    html += '<nav id="zt-path" style="font-size:large;margin:30px 0;--bs-breadcrumb-divider:\'>\';"></nav>';
    html += '<div class="row"><div class="col-md-12"><table class="table triage-listing-table"><colgroup><col style="width:2%"><col style="width:2%"><col style="width:71%"><col style="width:15%"><col style="width:10%"></colgroup>';
    html += '<thead><th><i id="zt-selectAllTick" class="far fa-square checks"></i></th><th colspan="2">Name</th><th>Last modified</th><th>Size</th></thead>';
    html += '<tbody id="zt-triage-listing"></tbody></table></div></div>';
    html += '<form><input type="hidden" id="zt-upload-input-path" name="path"><input type="file" id="zt-upload-input-file" name="file" style="visibility:hidden"></form>';

    html += '<div class="offcanvas offcanvas-end" id="zt-item-details" tabindex="-1"><div class="offcanvas-header"><h5 class="offcanvas-title">Details</h5><button type="button" class="btn-close" data-bs-dismiss="offcanvas"></button></div><div class="offcanvas-body" id="zt-item-details-content"></div></div>';

    html += '<div class="modal modal-lg fade" id="zt-mass-action-modal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Mass action</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><table class="table"><thead><tr><th>Path</th><th>Type</th></tr></thead><tbody id="zt-mass-action-list"></tbody></table></div><div class="modal-footer"><button type="button" class="btn btn-primary" id="zt-mass-import-btn">Import</button><button type="button" class="btn btn-danger" id="zt-mass-delete-btn">Delete</button></div></div></div></div>';

    html += '<div class="modal modal-xl fade" id="zt-mass-import-modal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Mass import</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><p class="mb-4">Select properties applied to all selected videos</p>';
    html += '<div class="form-floating mb-3"><select class="form-select" id="zt-channel-list"><option value="">None</option></select><label for="zt-channel-list">Channel</label></div>';
    html += '<div class="mb-3"><label>Actors</label><div id="zt-mass-actors" class="border rounded p-2" style="max-height:120px;overflow-y:auto"></div></div>';
    html += '<div class="mb-3"><label>Categories</label><div id="zt-mass-categories" class="border rounded p-2" style="max-height:120px;overflow-y:auto"></div></div>';
    html += '</div><div class="modal-footer"><span class="input-group-text">Import as:</span><button type="button" class="btn btn-primary" data-type="v">Video</button><button type="button" class="btn btn-primary" data-type="m">Movie</button><button type="button" class="btn btn-primary" data-type="c">Clip</button></div></div></div></div>';

    html += '<div class="modal fade modal-lg" id="zt-video-import-modal" tabindex="-1"><div class="modal-dialog modal-dialog-centered"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Import video</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><p>Filename: <code id="zt-import-filename"></code></p><p>Path: <code id="zt-import-filepath"></code></p></div><div class="modal-footer"><span class="input-group-text">Import as:</span><button type="button" class="btn btn-primary" data-type="v">Video</button><button type="button" class="btn btn-primary" data-type="m">Movie</button><button type="button" class="btn btn-primary" data-type="c">Clip</button></div></div></div></div>';

    html += '<div class="modal fade" id="newFolderModal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Create new folder</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><div class="form-floating"><input type="text" class="form-control" id="folder-new" placeholder="Name"><label for="folder-new">New folder name</label></div></div><div class="modal-footer"><button type="button" class="btn btn-success" id="zt-folder-create-btn">Create</button><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

    html += '</div></div>';
    self.innerHTML = html;

    self.querySelectorAll("#zt-upload-sidebar .nav-link[data-library-id]").forEach(function(a) {
      a.onclick = function(e) {
        e.preventDefault();
        var id = a.getAttribute("data-library-id");
        if (id) {
          selectedLibraryId = id;
          try { sessionStorage.setItem("zt-upload-library", id); } catch (e) {}
          loadListing();
        }
      };
    });

    var pathEl = self.querySelector("#zt-path");
    var pathArr = currentPath.split("/").filter(Boolean);
    var prog = "";
    var bc = '<ol class="breadcrumb"><li class="breadcrumb-item' + (currentPath === "/" ? " active" : " br-link") + '" data-path="/">Triage</li>';
    for (var i = 0; i < pathArr.length; i++) {
      prog += "/" + pathArr[i];
      bc += '<li class="breadcrumb-item' + (i === pathArr.length - 1 ? " active" : " br-link") + '" data-path="' + prog + '">' + pathArr[i] + "</li>";
    }
    bc += "</ol>";
    pathEl.innerHTML = bc;
    pathEl.querySelectorAll(".br-link, [data-path]").forEach(function(li) {
      var p = li.getAttribute("data-path");
      if (p) li.onclick = function() { currentPath = p; loadListing(); };
    });

    self.querySelector("#zt-upload-input-path").value = currentPath;

    function previewUrl(filePath) {
      var p = filePath.replace(/^\//, "");
      var q = selectedLibraryId ? "?library_id=" + encodeURIComponent(selectedLibraryId) : "";
      return "/api/upload/preview/" + encodeURIComponent(p) + q;
    }

    function getSelectedFiles() {
      var list = [];
      self.querySelectorAll("#zt-triage-listing tr[data-type=file]").forEach(function(tr) {
        var tick = tr.querySelector(".checks");
        if (tick && tick.classList.contains("fa-check-square")) list.push(tr);
      });
      return list;
    }

    function updateSelects() {
      var selected = 0;
      self.querySelectorAll("#zt-triage-listing tr[data-type=file]").forEach(function(tr) {
        var tick = tr.querySelector(".checks");
        if (tick && tick.classList.contains("fa-check-square")) selected++;
      });
      var selectAll = self.querySelector("#zt-selectAllTick");
      var massBtn = self.querySelector("#zt-mass-action-btn");
      if (selectAll) {
        var total = self.querySelectorAll("#zt-triage-listing tr[data-type=file]").length;
        selectAll.className = total > 0 && selected === total ? "far fa-check-square checks" : "far fa-square checks";
      }
      if (massBtn) {
        massBtn.classList.toggle("disabled", selected === 0);
      }
    }

    self.querySelector("#zt-selectAllTick").onclick = function() {
      var all = self.querySelectorAll("#zt-triage-listing tr[data-type=file]");
      var anyUnchecked = false;
      all.forEach(function(tr) {
        var tick = tr.querySelector(".checks");
        if (tick && tick.classList.contains("fa-square")) anyUnchecked = true;
      });
      all.forEach(function(tr) {
        var tick = tr.querySelector(".checks");
        if (tick) {
          tick.className = anyUnchecked ? "far fa-check-square checks" : "far fa-square checks";
        }
      });
      updateSelects();
    };

    detailOffcanvas = new bootstrap.Offcanvas(self.querySelector("#zt-item-details"));
    massActionModal = new bootstrap.Modal(self.querySelector("#zt-mass-action-modal"));
    massImportModal = new bootstrap.Modal(self.querySelector("#zt-mass-import-modal"));
    videoImportModal = new bootstrap.Modal(self.querySelector("#zt-video-import-modal"));

    function showFileDetails(fileName, filePath, fileType, fileSize, fileDate) {
      var niceTypes = { unknown: "Unrecognized", video: "Video", image: "Picture", archive: "Archive" };
      var content = '<h6>File details</h6><hr /><p><b>Name</b><br/>' + (fileName || "").replace(/</g,"&lt;") + '</p>';
      content += '<p><b>Type</b><br/>' + (niceTypes[fileType] || fileType) + '</p>';
      content += '<p><b>Size</b><br/>' + fileSize + '</p>';
      content += '<p><b>Last modified</b><br/>' + fileDate + '</p>';
      if (fileType === "video") {
        content += '<h6 class="mt-3">Preview</h6><hr /><video controls class="w-100" src="' + previewUrl(filePath) + '"></video>';
      } else if (fileType === "image") {
        content += '<h6 class="mt-3">Preview</h6><hr /><img class="w-100" src="' + previewUrl(filePath) + '" alt="">';
      }
      content += '<h6 class="mt-3">Actions</h6><hr />';
      if (fileType === "video") {
        content += '<button class="btn btn-outline-primary me-2" id="zt-single-import-btn"><i class="fas fa-file-import"></i> Import</button>';
      }
      content += '<a class="btn btn-outline-primary me-2" href="' + previewUrl(filePath) + '" target="_blank"><i class="fas fa-download"></i> Download</a>';
      content += '<button class="btn btn-outline-danger" id="zt-file-delete-btn"><i class="far fa-trash-alt"></i> Delete</button>';
      self.querySelector("#zt-item-details-content").innerHTML = content;

      var detailContent = self.querySelector("#zt-item-details-content");
      var singleImportBtn = detailContent.querySelector("#zt-single-import-btn");
      if (singleImportBtn) {
        singleImportBtn.onclick = function() {
          detailOffcanvas.hide();
          self.querySelector("#zt-import-filename").textContent = fileName;
          self.querySelector("#zt-import-filepath").textContent = filePath;
          self._importFilePath = filePath.replace(/^\//, "");
          videoImportModal.show();
        };
      }
      detailContent.querySelector("#zt-file-delete-btn").onclick = function() {
        var body = { File: filePath.replace(/^\//, "") };
        if (selectedLibraryId) body.library_id = selectedLibraryId;
        fetch("/api/upload/file", { method: "DELETE", credentials: "same-origin", headers: { "Content-Type": "application/json" }, body: JSON.stringify(body) })
          .then(function(r) {
            if (r.ok) { detailOffcanvas.hide(); loadListing(); if (typeof sendToast === "function") sendToast("Delete a file", "", "bg-success", "File deleted successfully"); }
            else r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Unable to delete file", "", "bg-warning", (d && d.error) || "Failed"); });
          });
      };
      detailOffcanvas.show();
    }

    self.querySelector("#zt-mass-action-btn").onclick = function() {
      if (this.classList.contains("disabled")) return;
      var rows = getSelectedFiles();
      var tbody = self.querySelector("#zt-mass-action-list");
      tbody.innerHTML = "";
      rows.forEach(function(tr) {
        tbody.innerHTML += "<tr><td>" + (tr.dataset.filepath || "").replace(/</g,"&lt;") + "</td><td>" + (tr.dataset.filetype || "") + "</td></tr>";
      });
      massActionModal.show();
    };

    self.querySelector("#zt-mass-import-btn").onclick = function() {
      massActionModal.hide();
      Promise.all([
        fetch("/api/channel", { credentials: "same-origin" }).then(function(r) { return r.json(); }),
        fetch("/api/actor", { credentials: "same-origin" }).then(function(r) { return r.json(); }),
        fetch("/api/category", { credentials: "same-origin" }).then(function(r) { return r.json(); })
      ]).then(function(arr) {
        var channels = arr[0].items || [];
        var actors = arr[1].items || [];
        var categories = arr[2].items || [];
        var chSel = self.querySelector("#zt-channel-list");
        chSel.innerHTML = '<option value="">None</option>';
        channels.forEach(function(ch) {
          chSel.innerHTML += '<option value="' + (ch.ID || ch.id) + '">' + (ch.Name || ch.name || "").replace(/</g,"&lt;") + "</option>";
        });
        var actDiv = self.querySelector("#zt-mass-actors");
        actDiv.innerHTML = "";
        actors.forEach(function(a) {
          var id = a.ID || a.id;
          actDiv.innerHTML += '<div class="form-check"><input class="form-check-input zt-mass-actor" type="checkbox" value="' + id + '"><label class="form-check-label">' + (a.Name || a.name || "").replace(/</g,"&lt;") + "</label></div>";
        });
        var catDiv = self.querySelector("#zt-mass-categories");
        catDiv.innerHTML = "";
        categories.forEach(function(c) {
          (c.Sub || c.sub || []).forEach(function(s) {
            var sid = s.ID || s.id;
            catDiv.innerHTML += '<div class="form-check"><input class="form-check-input zt-mass-cat" type="checkbox" value="' + sid + '"><label class="form-check-label">' + (s.Name || s.name || "").replace(/</g,"&lt;") + "</label></div>";
          });
        });
        massImportModal.show();
      });
    };

    self.querySelector("#zt-mass-import-modal").querySelectorAll(".modal-footer button[data-type]").forEach(function(btn) {
      btn.onclick = function() {
        var type = btn.getAttribute("data-type");
        var rows = getSelectedFiles();
        var files = rows.map(function(tr) { return tr.dataset.filepath || ""; });
        var channel = self.querySelector("#zt-channel-list").value || undefined;
        var actors = [];
        self.querySelectorAll("#zt-mass-actors .zt-mass-actor:checked").forEach(function(cb) { actors.push(cb.value); });
        var categories = [];
        self.querySelectorAll("#zt-mass-categories .zt-mass-cat:checked").forEach(function(cb) { categories.push(cb.value); });
        var payload = { files: files, type: type, channel: channel || "", actors: actors, categories: categories };
        if (selectedLibraryId) payload.library_id = selectedLibraryId;
        fetch("/api/upload/triage/mass-action", { method: "POST", credentials: "same-origin", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload) })
          .then(function(r) {
            if (r.ok) { massImportModal.hide(); loadListing(); if (typeof sendToast === "function") sendToast("Import successful", "", "bg-success", "Import tasks will run in background"); }
            else { massImportModal.hide(); loadListing(); r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Unable to import selected files", "", "bg-warning", (d && d.error) || "Failed"); }); }
          });
      };
    });

    self.querySelector("#zt-mass-delete-btn").onclick = function() {
      var rows = getSelectedFiles();
      var files = rows.map(function(tr) { return tr.dataset.filepath || ""; });
      var payload = { files: files };
      if (selectedLibraryId) payload.library_id = selectedLibraryId;
      fetch("/api/upload/triage/mass-action", { method: "DELETE", credentials: "same-origin", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload) })
        .then(function(r) {
          if (r.ok) { massActionModal.hide(); loadListing(); if (typeof sendToast === "function") sendToast("Deletion successful", "", "bg-success", "This folder should feel lighter now"); }
          else r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Unable to delete selected files", "", "bg-warning", (d && d.error) || "Failed"); });
        });
    };

    self.querySelector("#zt-video-import-modal").querySelectorAll("button[data-type]").forEach(function(btn) {
      btn.onclick = function() {
        var type = btn.getAttribute("data-type");
        var filepath = self._importFilePath;
        var filename = self.querySelector("#zt-import-filename").textContent || filepath;
        if (!filepath) return;
        var fd = new FormData();
        fd.set("name", filename);
        fd.set("filename", filepath);
        fd.set("type", type);
        if (selectedLibraryId) fd.set("library_id", selectedLibraryId);
        fetch("/api/video", { method: "POST", credentials: "same-origin", body: fd })
          .then(function(r) {
            if (r.ok) {
              return r.json().then(function(data) {
                videoImportModal.hide();
                detailOffcanvas.hide();
                loadListing();
                var vid = data && (data.video_id || data.id);
                if (typeof sendToast === "function") sendToast("Import a video", "", "bg-success", vid ? '<a href="/video/' + vid + '/edit" target="_blank">You can edit the video more here</a>' : "Import started");
              });
            }
            return r.json().catch(function(){return{};}).then(function(d) {
              if (typeof sendToast === "function") sendToast("Unable to import video", "", "bg-warning", (d && d.error) || "Import failed");
            });
          });
      };
    });

    self.querySelector("#zt-upload-file-btn").onclick = function() {
      var input = self.querySelector("#zt-upload-input-file");
      input.onchange = function(e) {
        var fd = new FormData();
        fd.set("path", currentPath);
        fd.set("file", e.target.files[0]);
        if (selectedLibraryId) fd.set("library_id", selectedLibraryId);
        fetch("/api/upload/file", { method: "POST", credentials: "same-origin", body: fd }).then(function(r) { if (r.ok) loadListing(); });
        input.value = "";
      };
      input.click();
    };

    self.querySelector("#zt-folder-create-btn").onclick = function() {
      var name = self.querySelector("#folder-new").value;
      if (!name) return;
      var fullPath = currentPath === "/" ? "/" + name : currentPath + "/" + name;
      var body = "name=" + encodeURIComponent(fullPath);
      if (selectedLibraryId) body += "&library_id=" + encodeURIComponent(selectedLibraryId);
      fetch("/api/upload/folder", { method: "POST", credentials: "same-origin", headers: { "Content-Type": "application/x-www-form-urlencoded" }, body: body })
        .then(function(r) {
          if (r.ok) {
            bootstrap.Modal.getInstance(self.querySelector("#newFolderModal")).hide();
            self.querySelector("#folder-new").value = "";
            loadListing();
            if (typeof sendToast === "function") sendToast("New folder created", "", "bg-success", "Folder " + fullPath + " created");
          } else r.json().catch(function(){return{};}).then(function(d) { if (typeof sendToast === "function") sendToast("Unable to create folder", "", "bg-warning", (d && d.error) || "Failed"); });
        });
    };

    Promise.all([
      (function() {
        var body = "path=" + encodeURIComponent(currentPath);
        if (selectedLibraryId) body += "&library_id=" + encodeURIComponent(selectedLibraryId);
        return fetch("/api/upload/triage/folder", { method: "POST", credentials: "same-origin", headers: { "Content-Type": "application/x-www-form-urlencoded" }, body: body }).then(function(r) { return r.json(); });
      })(),
      (function() {
        var body = "path=" + encodeURIComponent(currentPath);
        if (selectedLibraryId) body += "&library_id=" + encodeURIComponent(selectedLibraryId);
        return fetch("/api/upload/triage/file", { method: "POST", credentials: "same-origin", headers: { "Content-Type": "application/x-www-form-urlencoded" }, body: body }).then(function(r) { return r.json(); });
      })()
    ]).then(function(results) {
      var folderData = results[0].folders || {};
      var filesData = results[1].files || {};
      var tbody = self.querySelector("#zt-triage-listing");
      var rowsHtml = "";
      Object.keys(folderData).forEach(function(name) {
        var safe = name.replace(/&/g,"&amp;").replace(/</g,"&lt;");
        rowsHtml += '<tr data-type="folder" data-folder="' + safe + '"><td></td><td><i class="far fa-folder"></i></td><td>' + safe + '</td><td>—</td><td>' + folderData[name] + ' items</td></tr>';
      });
      Object.keys(filesData).forEach(function(name) {
          var info = filesData[name] || {};
          var sz = info.Size !== undefined ? humanFileSize(info.Size, true) : "—";
          var dt = info.LastModification ? new Date(info.LastModification).toLocaleDateString() : "—";
          var filePath = currentPath === "/" ? "/" + name : currentPath + "/" + name;
          var fileType = "unknown";
          var icon = "fas fa-file";
          if (/\.(mp4|mkv|webm)$/i.test(name)) { fileType = "video"; icon = "far fa-play-circle"; }
          else if (/\.(png|jpg|jpeg)$/i.test(name)) { fileType = "image"; icon = "far fa-image"; }
          else if (/\.(zip|tar)$/i.test(name)) { fileType = "archive"; icon = "far fa-file-archive"; }
          var safeName = name.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;");
          var safePath = filePath.replace(/&/g,"&amp;").replace(/"/g,"&quot;");
          rowsHtml += '<tr data-type="file" data-filename="' + safeName + '" data-filepath="' + safePath + '" data-filetype="' + fileType + '">';
          rowsHtml += '<td class="zt-check-cell"><i class="far fa-square checks"></i></td>';
          rowsHtml += '<td><i class="' + icon + '"></i></td><td>' + safeName + '</td><td>' + dt + '</td><td>' + sz + '</td></tr>';
        });
      tbody.innerHTML = rowsHtml;

      self.querySelectorAll("#zt-triage-listing tr[data-type=folder]").forEach(function(tr) {
        tr.ondblclick = function() {
          var folder = tr.getAttribute("data-folder");
          if (folder) { currentPath = currentPath === "/" ? "/" + folder : currentPath + "/" + folder; loadListing(); }
        };
      });

      self.querySelectorAll("#zt-triage-listing tr[data-type=file]").forEach(function(tr) {
          var checkCell = tr.querySelector(".zt-check-cell");
          if (checkCell) {
            checkCell.onclick = function(e) {
              e.stopPropagation();
              var t = tr.querySelector(".checks");
              if (t) { t.classList.toggle("fa-check-square"); t.classList.toggle("fa-square"); }
              updateSelects();
            };
          }
          tr.onclick = function(e) {
            if (e.target.closest(".zt-check-cell")) return;
            var fn = tr.dataset.filename, fp = tr.dataset.filepath, ft = tr.dataset.filetype;
            var sz = tr.cells[4] ? tr.cells[4].textContent : "—";
            var dt = tr.cells[3] ? tr.cells[3].textContent : "—";
            showFileDetails(fn, fp, ft, sz, dt);
          };
        });
      updateSelects();
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    }).catch(function(err) {
      if (typeof sendToast === "function") sendToast("Unable to show folders or files", "", "bg-warning", (err && err.message) || "Failed to load triage");
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });

    self._enterFolder = function(folderName) {
      currentPath = currentPath === "/" ? "/" + folderName : currentPath + "/" + folderName;
      loadListing();
    };
    self._updateSelects = updateSelects;
  }

  function loadListing() {
    updateView();
  }
  loadListing();
  });
};
customElements.define("zt-upload-home", ZtUploadHome);
})();
