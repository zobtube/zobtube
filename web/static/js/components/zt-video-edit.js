(function() {
"use strict";
function esc(s) { return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;").replace(/>/g,"&gt;"); }
function escAttr(s) { return String(s).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;").replace(/>/g,"&gt;"); }
function niceDur(ns) {
  if (!ns) return "Unknown";
  var s = Math.floor(ns/1e9), m = Math.floor(s/60); s %= 60;
  var h = Math.floor(m/60); m %= 60;
  return h>0 ? (h<10?"0":"")+h+":"+(m<10?"0":"")+m+":"+(s<10?"0":"")+s : (m<10?"0":"")+m+":"+(s<10?"0":"")+s;
}
function fmtDate(t) {
  if (!t) return "";
  var d = new Date(t);
  return d.getFullYear()+"-"+String(d.getMonth()+1).padStart(2,0)+"-"+String(d.getDate()).padStart(2,0)+" "+String(d.getHours()).padStart(2,0)+":"+String(d.getMinutes()).padStart(2,0)+":"+String(d.getSeconds()).padStart(2,0);
}
var typeStr = { c: "clip", v: "video", m: "movie" };

function ZtVideoEdit() {
  var el = Reflect.construct(HTMLElement, [], ZtVideoEdit);
  return el;
}
ZtVideoEdit.prototype = Object.create(HTMLElement.prototype);
ZtVideoEdit.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = '<div class="alert alert-danger">Missing video ID</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  if (!(window.__USER__ && window.__USER__.admin)) {
    self.innerHTML = '<div class="alert alert-danger">Forbidden</div>';
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    return;
  }
  fetch("/api/video/" + encodeURIComponent(id) + "/edit", { credentials: "same-origin" })
    .then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
    .then(function(data) {
      var v = data.video || data.Video;
      var actors = data.actors || data.Actors || [];
      var categories = data.categories || data.Categories || [];
      if (!v) { self.innerHTML = '<div class="alert alert-danger">Video not found</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }

      var urlView = (v.Type||v.type)==="c" ? "/clip/"+id : "/video/"+id;
      var streamUrl = "/api/video/"+id+"/stream";
      var name = esc(v.Name||v.name||"");
      var filename = esc(v.Filename||v.filename||"");
      var imported = v.Imported||v.imported;
      var hasThumb = v.Thumbnail||v.thumbnail;
      var hasThumbMini = v.ThumbnailMini||v.thumbnailMini;
      var dur = niceDur(v.Duration||v.duration);
      var typeStrVal = typeStr[v.Type||v.type] || "video";
      var channel = v.Channel||v.channel;
      var channelName = channel ? (channel.Name||channel.name||"") : "None";
      var vidActors = v.Actors||v.actors||[];
      var vidCats = v.Categories||v.categories||[];

      var actorSelectedIds = vidActors.map(function(a){ return a.ID||a.id; });
      var actorSelectable = {};
      actors.forEach(function(a){ actorSelectable[a.ID||a.id] = { name: a.Name||a.name||"" }; });
      var categorySelectedIds = vidCats.map(function(c){ return c.ID||c.id; });
      var categorySelectable = {};
      categories.forEach(function(c){
        (c.Sub||c.sub||[]).forEach(function(s){ categorySelectable[s.ID||s.id] = { name: s.Name||s.name||"" }; });
      });

      window.zt.actorSelection = window.zt.actorSelection || {};
      window.zt.actorSelection.actorSelectable = actorSelectable;
      window.zt.actorSelection.actorSelected = {};
      actorSelectedIds.forEach(function(aid){ window.zt.actorSelection.actorSelected[aid]=undefined; });
      window.zt.actorSelection.onActorSelectBefore = function(aid){
        var name = (actorSelectable[aid]&&actorSelectable[aid].name)||"Actor";
        fetch("/api/video/"+id+"/actor/"+aid, { method: "PUT", credentials: "same-origin" }).then(function(r){
          if(r.ok && typeof sendToast==="function") sendToast("Actor added","","bg-success",name+" added.");
          else if(!r.ok && typeof sendToast==="function") sendToast("Actor not added","","bg-danger",name+" not added, call failed.");
        });
        return true;
      };
      window.zt.actorSelection.onActorDeselectBefore = function(aid){
        var name = (actorSelectable[aid]&&actorSelectable[aid].name)||"Actor";
        fetch("/api/video/"+id+"/actor/"+aid, { method: "DELETE", credentials: "same-origin" }).then(function(r){
          if(r.ok && typeof sendToast==="function") sendToast("Actor removed","","bg-success",name+" removed.");
          else if(!r.ok && typeof sendToast==="function") sendToast("Actor not removed","","bg-danger",name+" not removed, call failed.");
        });
        return true;
      };
      var actorSel = window.zt.actorSelection;
      actorSel._updateSelectedActors = function(){
        var chips = Array.prototype.slice.call(document.getElementsByClassName("video-actor-list"));
        for(var i=0;i<chips.length;i++){
          var aid = chips[i].getAttribute("actor-id");
          if (!(aid in actorSel.actorSelected)) chips[i].remove();
          else chips[i].style.display = "";
        }
        var addChips = document.getElementsByClassName("add-actor-list");
        for(i=0;i<addChips.length;i++){
          aid = addChips[i].getAttribute("actor-id");
          var addBtn = addChips[i].querySelector(".btn-success");
          var delBtn = addChips[i].querySelector(".btn-danger");
          if(aid in actorSel.actorSelected){ if(addBtn)addBtn.style.display="none"; if(delBtn)delBtn.style.display=""; addChips[i].style.display="none"; }
          else { if(addBtn)addBtn.style.display=""; if(delBtn)delBtn.style.display="none"; addChips[i].style.display=""; }
        }
      };
      actorSel.actorSelect = function(aid){ actorSel.actorSelected[aid]=undefined; actorSel._updateSelectedActors(); };
      actorSel.actorDeselect = function(aid){
        var name = (actorSelectable[aid]&&actorSelectable[aid].name)||"Actor";
        fetch("/api/video/"+id+"/actor/"+aid, { method: "DELETE", credentials: "same-origin" })
          .then(function(r){
            if(r.ok){ delete actorSel.actorSelected[aid]; actorSel._updateSelectedActors(); if(typeof sendToast==="function") sendToast("Actor removed","","bg-success",name+" removed."); }
            else if(typeof sendToast==="function") sendToast("Actor not removed","","bg-danger",name+" not removed, call failed.");
          });
      };

      window.zt.categorySelection = window.zt.categorySelection || {};
      window.zt.categorySelection.categorySelectable = categorySelectable;
      window.zt.categorySelection.categorySelected = {};
      categorySelectedIds.forEach(function(cid){ window.zt.categorySelection.categorySelected[cid]=undefined; });
      window.zt.categorySelection.onCategorySelectBefore = function(cid){
        var name = (categorySelectable[cid]&&categorySelectable[cid].name)||"Category";
        fetch("/api/video/"+id+"/category/"+cid, { method: "PUT", credentials: "same-origin" }).then(function(r){
          if(r.ok && typeof sendToast==="function") sendToast("Category added","","bg-success",name+" added.");
          else if(!r.ok && typeof sendToast==="function") sendToast("Category not added","","bg-danger",name+" not added, call failed.");
        });
        return true;
      };
      window.zt.categorySelection.onCategoryDeselectBefore = function(cid){
        var name = (categorySelectable[cid]&&categorySelectable[cid].name)||"Category";
        fetch("/api/video/"+id+"/category/"+cid, { method: "DELETE", credentials: "same-origin" }).then(function(r){
          if(r.ok && typeof sendToast==="function") sendToast("Category removed","","bg-success",name+" removed.");
          else if(!r.ok && typeof sendToast==="function") sendToast("Category not removed","","bg-danger",name+" not removed, call failed.");
        });
        return true;
      };
      var catSel = window.zt.categorySelection;
      catSel._updateSelectedCategories = function(){
        var chips = Array.prototype.slice.call(document.getElementsByClassName("video-category-list"));
        for(var j=0;j<chips.length;j++){
          var cid = chips[j].getAttribute("category-id");
          if (!(cid in catSel.categorySelected)) chips[j].remove();
          else chips[j].style.display = "";
        }
        var addChips = document.getElementsByClassName("add-category-list");
        for(var k=0;k<addChips.length;k++){
          cid = addChips[k].getAttribute("category-id");
          var addBtn = addChips[k].querySelector(".btn-success");
          var delBtn = addChips[k].querySelector(".btn-danger");
          if(cid in catSel.categorySelected){ if(addBtn)addBtn.style.display="none"; if(delBtn)delBtn.style.display=""; }
          else { if(addBtn)addBtn.style.display=""; if(delBtn)delBtn.style.display="none"; }
        }
      };
      catSel.categorySelect = function(cid){ catSel.categorySelected[cid]=undefined; catSel._updateSelectedCategories(); };
      catSel.categoryDeselect = function(cid){
        var name = (categorySelectable[cid]&&categorySelectable[cid].name)||"Category";
        fetch("/api/video/"+id+"/category/"+cid, { method: "DELETE", credentials: "same-origin" })
          .then(function(r){
            if(r.ok){ delete catSel.categorySelected[cid]; catSel._updateSelectedCategories(); if(typeof sendToast==="function") sendToast("Category removed","","bg-success",name+" removed."); }
            else if(typeof sendToast==="function") sendToast("Category not removed","","bg-danger",name+" not removed, call failed.");
          });
      };

      var actorChips = actors.map(function(a){
        var aid = a.ID||a.id;
        var show = actorSelectedIds.indexOf(aid)>=0 ? "" : "display:none";
        return '<div class="chip video-actor-list" actor-id="'+aid+'" style="'+show+'"><img src="/api/actor/'+encodeURIComponent(aid)+'/thumb" width="50" height="50">'+esc(a.Name||a.name)+'<button class="btn btn-danger" onclick="window.zt.actorSelection.actorDeselect(\''+aid+'\');"><i class="fa fa-trash-alt"></i></button></div>';
      }).join("");
      var categoryChips = [];
      categories.forEach(function(c){
        (c.Sub||c.sub||[]).forEach(function(s){
          var sid = s.ID||s.id;
          var show = categorySelectedIds.indexOf(sid)>=0 ? "" : "display:none";
          var thumb = (s.Thumbnail||s.thumbnail) ? '<img src="/api/category-sub/'+encodeURIComponent(sid)+'/thumb" width="100" height="50">' : "";
          categoryChips.push('<div class="chip video-category-list" category-id="'+sid+'" style="'+show+'">'+thumb+esc(s.Name||s.name)+'<button class="btn btn-danger" onclick="window.zt.categorySelection.categoryDeselect(\''+sid+'\');"><i class="fa fa-trash-alt"></i></button></div>');
        });
      });

      var actorModalChips = actors.map(function(a){
        var aid = a.ID||a.id;
        var sel = actorSelectedIds.indexOf(aid)>=0;
        return '<div class="chip add-actor-list" actor-id="'+aid+'" style="'+(sel?'display:none;':'')+'"><img src="/api/actor/'+encodeURIComponent(aid)+'/thumb" width="50" height="50">'+esc(a.Name||a.name)+'<button class="btn btn-success add-actor-add"><i class="fa fa-plus-circle"></i></button><button class="btn btn-danger add-actor-remove" style="'+(sel?'':'display:none')+'"><i class="fa fa-trash-alt"></i></button></div>';
      }).join("");

      var categoryModalHtml = "";
      categories.forEach(function(c){
        var subs = c.Sub||c.sub||[];
        if(subs.length===0) return;
        categoryModalHtml += '<h4 class="mt-3">'+esc(c.Name||c.name)+'</h4><div class="chips">';
        subs.forEach(function(s){
          var sid = s.ID||s.id;
          var sel = categorySelectedIds.indexOf(sid)>=0;
          var thumb = (s.Thumbnail||s.thumbnail) ? '<img class="lazy" data-src="/api/category-sub/'+encodeURIComponent(sid)+'/thumb" width="100" height="50">' : "";
          categoryModalHtml += '<div class="chip add-category-list" category-id="'+sid+'">'+thumb+esc(s.Name||s.name)+'<button class="btn btn-success add-category-add"><i class="fa fa-plus-circle"></i></button><button class="btn btn-danger add-category-remove" style="'+(sel?'':'display:none')+'"><i class="fa fa-trash-alt"></i></button></div>';
        });
        categoryModalHtml += '</div>';
      });

      var html = '<h2>Video editing</h2><a href="'+urlView+'">← Back to video viewer</a><br/><div class="row">';
      html += '<div class="col-md-9"><video id="zt-video-edit-player" style="width:100%;height:35vw" src="'+streamUrl+'" controls></video></div>';
      html += '<div class="col-md-3"><h5>Import status</h5><p><span id="video-import" class="badge '+(imported?'bg-success':'bg-warning')+'">'+(imported?'Imported':'In triage')+'</span></p>';
      html += '<h5>Duration</h5><p><small class="text-muted" id="video-duration">'+dur+'</small></p>';
      html += '<h5>Thumbnail</h5><p><span id="video-thumb" class="badge '+(hasThumb?'bg-success':'bg-warning')+'">'+(hasThumb?'Generated':'Missing')+'</span></p>';
      html += '<h5>Thumbnail mini</h5><p><span id="video-thumb-mini" class="badge '+(hasThumbMini?'bg-success':'bg-warning')+'">'+(hasThumbMini?'Generated':'Missing')+'</span></p><hr/>';
      html += '<h5>Actions</h5><p><button class="btn btn-primary btn-sm w-100" id="zt-gen-thumb">New thumbnail from current timecode</button></p>';
      html += '<p><button class="btn btn-primary btn-sm w-100" '+(v.Type==='c'?'disabled':'')+' data-migrate="c">Change to Clip</button></p>';
      html += '<p><button class="btn btn-primary btn-sm w-100" '+(v.Type==='m'?'disabled':'')+' data-migrate="m">Change to Movie</button></p>';
      html += '<p><button class="btn btn-primary btn-sm w-100" '+(v.Type==='v'?'disabled':'')+' data-migrate="v">Change to Video</button></p>';
      html += '<p><button class="btn btn-danger btn-sm w-100" id="zt-del-video">Delete</button></p></div>';

      html += '<div class="col-12 mb-3 mt-3"><h4>Video details</h4></div>';
      html += '<div class="col-12 mb-3"><div class="form-floating input-group"><input type="text" class="form-control" disabled id="video-title" value="'+name+'"><label for="video-title">Name</label><button class="btn btn-outline-warning" type="button" id="video-title-edit">Edit</button></div></div>';
      html += '<div class="col-6 mb-3"><div class="form-floating"><input type="text" disabled class="form-control" id="video-id" value="'+esc(id)+'"><label>ID</label></div></div>';
      html += '<div class="col-6 mb-3"><div class="form-floating"><input type="text" disabled class="form-control" value="'+filename+'"><label>Original filename</label></div></div>';
      html += '<div class="col-6 mb-3"><div class="form-floating"><input type="text" disabled class="form-control" value="'+fmtDate(v.CreatedAt||v.created_at)+'"><label>Import date</label></div></div>';
      html += '<div class="col-6 mb-3"><div class="form-floating"><input type="text" disabled class="form-control" value="'+esc(typeStrVal)+'"><label>Video Type</label></div></div>';
      html += '<div class="col-12 mb-3"><div class="form-floating input-group"><input type="text" disabled class="form-control" id="video-channel" value="'+esc(channelName)+'"><label>Channel</label><button class="btn btn-outline-warning" type="button" id="video-channel-edit">Change</button></div></div>';
      html += '<div class="col-12 mb-3"><div class="form-floating"><div class="form-control chip-selector" style="height:unset;display:flex;"><div class="chips">'+actorChips+'<div class="chip">Add an actor<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#actorSelectionModal"><i class="fa fa-plus-circle"></i></button></div></div></div><label>Actors</label></div></div>';
      html += '<div class="col-12 mb-3"><div class="form-floating"><div class="form-control chip-selector" style="height:unset;display:flex;"><div class="chips">'+categoryChips.join('')+'<div class="chip">Add a category<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#categorySelectionModal"><i class="fa fa-plus-circle"></i></button></div></div></div><label>Categories</label></div></div>';
      html += '</div>';

      html += '<div class="modal fade" id="actorSelectionModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Add actor in video</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><div class="form-floating mb-3"><input type="text" class="form-control" id="actorSelectionModalInput" autocomplete="off"><label for="actorSelectionModalInput">Actor name</label></div><div class="chips">'+actorModalChips+'<div class="chip">Add a new actor<button class="btn btn-success" onclick="window.open(\'/actor/new\',\'_blank\');"><i class="fa fa-plus-circle"></i></button></div></div></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div></div>';
      html += '<div class="modal fade" id="categorySelectionModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Add category in video</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body">'+categoryModalHtml+'</div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';
      html += '<div class="modal fade" id="editChannelModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Change video channel</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><div class="form-floating mb-3"><select class="form-select" id="channel-list"></select><label for="channel-list">Channel list</label></div></div><div class="modal-footer"><button type="button" class="btn btn-success" id="channel-send">Change</button><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

      self.innerHTML = html;

      var player = self.querySelector("#zt-video-edit-player");
      var videoTiming = "00:00:00";
      if (player) player.addEventListener("timeupdate", function(){
        var t = this.currentTime;
        var d = new Date(null); d.setMilliseconds(t*1000);
        videoTiming = d.toISOString().substr(11,8);
      });

      self.querySelector("#zt-gen-thumb").addEventListener("click", function(){
        fetch("/api/video/"+id+"/generate-thumbnail/"+videoTiming, { method: "POST", credentials: "same-origin" })
          .then(function(r){ if(r.ok){ var e=self.querySelector("#video-thumb"); e.textContent="Generated"; e.className="badge bg-success"; } });
      });
      self.querySelector("#zt-del-video").addEventListener("click", function(){
        if(!confirm("Delete this video?")) return;
        fetch("/api/video/"+id, { method: "DELETE", credentials: "same-origin" })
          .then(function(r){ if(r.ok) window.location.replace("/"); });
      });
      self.querySelectorAll("[data-migrate]").forEach(function(btn){
        btn.addEventListener("click", function(){
          var typ = this.getAttribute("data-migrate");
          var fd = new FormData(); fd.set("new_type", typ);
          fetch("/api/video/"+id+"/migrate", { method: "POST", credentials: "same-origin", body: fd })
            .then(function(r){
              if (r.ok && typeof sendToast==="function") sendToast("Change video type","","bg-success","Successful");
              else if (!r.ok) r.json().catch(function(){return{};}).then(function(d){ if(typeof sendToast==="function") sendToast("Change video type","failed","bg-danger",(d&&d.error)||"Failed"); });
            });
        });
      });

      var titleEl = self.querySelector("#video-title");
      var titleBtn = self.querySelector("#video-title-edit");
      var inEditMode = false;
      titleBtn.addEventListener("click", function(){
        if(!inEditMode){ titleEl.disabled=false; titleBtn.textContent="Send"; inEditMode=true; }
        else {
          var fd = new FormData(); fd.set("name", titleEl.value);
          fetch("/api/video/"+id+"/rename", { method: "POST", credentials: "same-origin", body: fd })
            .then(function(r){ if(r.ok){ titleEl.disabled=true; titleBtn.textContent="Edit"; inEditMode=false; } });
        }
      });

      self.querySelector("#video-channel-edit").addEventListener("click", function(){
        fetch("/api/channel/map", { credentials: "same-origin" }).then(function(r){ return r.json(); })
          .then(function(data){
            var ch = data.channels||{};
            var sel = self.querySelector("#channel-list");
            sel.innerHTML = '<option value="x">None</option>'+Object.keys(ch).map(function(k){ return '<option value="'+escAttr(k)+'">'+esc(ch[k]||"")+'</option>'; }).join("");
            new bootstrap.Modal(self.querySelector("#editChannelModal")).show();
          });
      });
      self.querySelector("#channel-send").addEventListener("click", function(){
        var cid = self.querySelector("#channel-list").value;
        var fd = new FormData(); fd.set("channelID", cid);
        fetch("/api/video/"+id+"/channel", { method: "POST", credentials: "same-origin", body: fd })
          .then(function(r){ if(r.ok) window.location.reload(); });
      });

      self.querySelectorAll(".add-actor-list .add-actor-add").forEach(function(btn){
        btn.addEventListener("click", function(){
          var aid = this.closest(".add-actor-list").getAttribute("actor-id");
          var name = (actorSelectable[aid]&&actorSelectable[aid].name)||"Actor";
          fetch("/api/video/"+id+"/actor/"+aid, { method: "PUT", credentials: "same-origin" })
            .then(function(r){
              if(r.ok){ window.zt.actorSelection.actorSelect(aid); if(typeof sendToast==="function") sendToast("Actor added","","bg-success",name+" added."); }
              else if(typeof sendToast==="function") sendToast("Actor not added","","bg-danger",name+" not added, call failed.");
            });
        });
      });
      self.querySelectorAll(".add-actor-list .add-actor-remove").forEach(function(btn){
        btn.addEventListener("click", function(){ var aid = this.closest(".add-actor-list").getAttribute("actor-id"); window.zt.actorSelection.actorDeselect(aid); });
      });
      self.querySelectorAll(".add-category-list .add-category-add").forEach(function(btn){
        btn.addEventListener("click", function(){
          var cid = this.closest(".add-category-list").getAttribute("category-id");
          var name = (categorySelectable[cid]&&categorySelectable[cid].name)||"Category";
          fetch("/api/video/"+id+"/category/"+cid, { method: "PUT", credentials: "same-origin" })
            .then(function(r){
              if(r.ok){ window.zt.categorySelection.categorySelect(cid); if(typeof sendToast==="function") sendToast("Category added","","bg-success",name+" added."); }
              else if(typeof sendToast==="function") sendToast("Category not added","","bg-danger",name+" not added, call failed.");
            });
        });
      });
      self.querySelectorAll(".add-category-list .add-category-remove").forEach(function(btn){
        btn.addEventListener("click", function(){ var cid = this.closest(".add-category-list").getAttribute("category-id"); window.zt.categorySelection.categoryDeselect(cid); });
      });

      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function(e){
      if(e && e.message==="403"){ self.innerHTML = '<div class="alert alert-danger">Forbidden</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
      if(e && e.message==="404"){ self.innerHTML = '<div class="alert alert-danger">Video not found</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
      self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-video-edit", ZtVideoEdit);
})();
