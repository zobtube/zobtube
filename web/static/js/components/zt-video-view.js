(function() {
"use strict";
function niceDurationShort(ns) {
  if (!ns) return "";
  var s = Math.floor(ns/1e9), m = Math.floor(s/60); s %= 60;
  var h = Math.floor(m/60); m %= 60;
  return h > 0 ? h + "h" + String(m).padStart(2,0) : m > 0 ? m + " min" : s + " sec";
}
function ZtVideoView() {
  var el = Reflect.construct(HTMLElement, [], ZtVideoView);
  return el;
}
ZtVideoView.prototype = Object.create(HTMLElement.prototype);
ZtVideoView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = "Missing id"; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  var pageSearch = typeof window !== "undefined"
    ? (window.location.search || window.__ZT_PAGE_SEARCH__ || "")
    : "";
  var playlistId = window.ztPlaylistIdFromSearch ? window.ztPlaylistIdFromSearch(pageSearch) : "";
  if (playlistId) self.setAttribute("data-playlist-id", playlistId);
  var shouldAutoplay = pageSearch.indexOf("autoplay") !== -1;
  var apiUrl = "/api/video/" + encodeURIComponent(id);
  if (playlistId) apiUrl += "?playlist=" + encodeURIComponent(playlistId);
  var fetchOpts = { credentials: "same-origin" };
  if (playlistId) fetchOpts.cache = "no-store";
  fetch(apiUrl, fetchOpts)
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(data) {
      var v = data.video || data;
      var viewCount = data.view_count || 0;
      var cats = data.categories || {};
      var randVideos = data.random_videos || [];
      var playlistCtx = data.playlist_video_ids ? data : null;
      var admin = (window.__USER__ && window.__USER__.admin);
      var esc = window.ztEscHtml;
      var name = esc(v.Name||v.name||v.Filename||v.filename||"Untitled");
      var streamUrl = data.stream_url || "/api/video/"+encodeURIComponent(id)+"/stream";
      var thumbUrl = "/api/video/"+encodeURIComponent(id)+"/thumb";
      var durShort = niceDurationShort(v.Duration||v.duration);
      var catsHtml = Object.keys(cats).map(function(cid){ return '<a class="btn btn-sm btn-secondary" href="/category/'+esc(cid)+'">'+esc(cats[cid]||"")+'</a>'; }).join("");
      var channel = v.Channel || v.channel;
      var channelHtml = channel ? '<a class="btn btn-sm btn-dark" href="/channel/'+esc(channel.ID||channel.id)+'"><i class="fas fa-tv"></i> '+esc(channel.Name||channel.name||"")+'</a>' : "";
      var actors = v.Actors || v.actors || [];
      var actorsHtml = actors.map(function(a){
        var sexIcon = (a.Sex||a.sex)==="f" ? "fa-venus" : (a.Sex||a.sex)==="m" ? "fa-mars" : (a.Sex||a.sex)==="tw" ? "fa-mars-and-venus" : "fa-person-circle-question";
        return '<a class="btn btn-sm btn-danger" href="/actor/'+esc(a.ID||a.id)+'"><span class="badge text-bg-light"><i class="fa '+sexIcon+'"></i></span> '+esc(a.Name||a.name||"")+'</a>';
      }).join("");
      var viewText = viewCount ? (viewCount > 1 ? viewCount + " views" : viewCount + " view") : "Not viewed yet!";
      var durBadgeStyle = 'height:fit-content;vertical-align:super;font-size:0.85rem;margin-left:8px';
      var plName = playlistCtx && playlistCtx.playlist ? (playlistCtx.playlist.name || "Playlist") : "";
      var mainCol = playlistCtx ? "col-lg-9 col-md-8" : "col-lg-12 col-md-12";
      var html = '<div class="row"><div class="' + mainCol + '"><div class="video-post-wrapper">';
      html += '<div><h3 id="page_view_video_title" class="post-title mt-3 d-inline-block">'+name+'</h3><span class="badge text-bg-secondary" style="'+durBadgeStyle+'">'+durShort+'</span></div>';
      html += '<div style="margin-top:0.5rem;display:flex;flex-wrap:wrap;gap:5px;align-items:center">'+channelHtml+' '+actorsHtml+' '+catsHtml+'</div>';
      var streamHref = esc(window.ztSafeUrl(streamUrl));
      html += '<div class="video-posts-video"><hr /><div class="ratio ratio-16x9"><video id="zt-main-video" style="width:100%" src="'+streamHref+'" preload="metadata" poster="'+esc(thumbUrl)+'" controls></video></div></div>';
      html += '<div class="video-posts-data"><div class="video-post-title"><div class="video-post-info"><h5 id="video-view-count"><i class="far fa-eye text-secondary"></i><span>'+esc(viewText)+'</span></h5></div></div>';
      html += '<div class="video-post-counter"><a download="'+esc(id)+'.mp4" href="'+streamHref+'" class="btn btn-sm btn-outline-dark"><i class="fas fa-download text-secondary"></i> Download</a> <zt-playlist-picker data-video-id="'+esc(id)+'"></zt-playlist-picker>'+(admin ? ' <a class="btn btn-sm btn-outline-dark" href="/video/'+esc(id)+'/edit"><i class="fa fa-edit text-secondary"></i> Edit</a>' : '')+'</div></div></div>';
      html += '</div>'; // close main column
      if (playlistCtx) {
        html += '<div class="col-lg-3 col-md-4"><div id="zt-playlist-up-next-sidebar"></div></div>';
      }
      html += '</div>';
      if (!playlistCtx) {
        html += '<div class="popular-videos"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-play"></i></span><h3>More Random Videos</h3></div><div class="row">';
        randVideos.forEach(function(rv){ html += '<div class="col-md-3"><zt-video-tile data-video="'+esc(JSON.stringify(rv))+'"></zt-video-tile></div>'; });
        html += "</div></div></div>";
      }
      self.innerHTML = html;
      if (playlistCtx && window.ztPlaylistRenderUpNext) {
        var sidebar = self.querySelector("#zt-playlist-up-next-sidebar");
        window.ztPlaylistRenderUpNext(sidebar, playlistCtx, playlistId, id);
      }
      var videoEl = self.querySelector("#zt-main-video") || self.querySelector("video");
      if (videoEl) {
        var viewCounted = false;
        videoEl.addEventListener("play", function() {
          if (viewCounted) return;
          viewCounted = true;
          fetch("/api/video/"+id+"/count-view", {method:"POST",credentials:"same-origin"}).then(function() {
            var s = self.querySelector("#video-view-count span");
            if (s) {
              var n = viewCount + 1;
              viewCount = n;
              s.textContent = n > 1 ? n + " views" : n + " view";
            }
          });
        });
        if (playlistCtx && window.ztPlaylistBindAutoAdvance) {
          window.ztPlaylistBindAutoAdvance(videoEl, playlistCtx, playlistId);
        }
        if (shouldAutoplay) {
          videoEl.play();
        }
      }
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Not found.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-video-view", ZtVideoView);
})();
