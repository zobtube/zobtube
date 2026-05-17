(function() {
"use strict";

var sidebarStyles =
  ".zt-playlist-up-next{position:sticky;top:70px;max-height:calc(100vh - 90px);overflow-y:auto}" +
  ".zt-playlist-up-next-item{display:flex;gap:8px;padding:8px;border-radius:4px;margin-bottom:4px;text-decoration:none;color:inherit}" +
  ".zt-playlist-up-next-item:hover:not(.active){background:#f8f8f8;color:#167ac6}" +
  ".zt-playlist-up-next-item.active{background:#f2f2f2}" +
  ".zt-playlist-up-next-item img{width:120px;height:68px;object-fit:cover;border-radius:4px;flex-shrink:0}" +
  ".zt-playlist-up-next-item .zt-up-next-title{font-size:0.9rem;line-height:1.3}";

function escapeHtml(s) {
  return String(s).replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

function videoPlayPath(v) {
  if (!v) return "#";
  var id = v.ID || v.id;
  var t = v.Type || v.type || "v";
  return (t === "c" ? "/clip/" : "/video/") + id;
}

function thumbUrl(v) {
  if (window.ztVideoUrlThumbXS) return window.ztVideoUrlThumbXS(v);
  var id = v && (v.ID || v.id);
  return id ? "/api/video/" + id + "/thumb_xs" : "";
}

window.ztPlaylistIdFromSearch = function(search) {
  search = search != null ? search : (typeof window !== "undefined" ? window.location.search : "");
  try {
    return new URLSearchParams(search).get("playlist") || "";
  } catch (e) {
    return "";
  }
};

window.ztPlaylistPlayUrl = function(video, playlistId, opts) {
  opts = opts || {};
  if (!video || !playlistId) return videoPlayPath(video);
  var path = videoPlayPath(video);
  var q = "playlist=" + encodeURIComponent(playlistId);
  if (opts.autoplay) q += "&autoplay=1";
  return path + "?" + q;
};

function playlistVideosFromCtx(ctx) {
  if (ctx.playlist_videos && ctx.playlist_videos.length) return ctx.playlist_videos;
  var list = [];
  var seen = {};
  var cur = ctx.video;
  if (cur) {
    var curId = cur.ID || cur.id;
    if (curId) { list.push(cur); seen[curId] = 1; }
  }
  (ctx.playlist_up_next || []).forEach(function(v) {
    var vid = v.ID || v.id;
    if (vid && !seen[vid]) { list.push(v); seen[vid] = 1; }
  });
  return list;
}

window.ztPlaylistRenderUpNext = function(container, ctx, playlistId, currentVideoId) {
  if (!container || !ctx) return;
  var pl = ctx.playlist || {};
  var videos = playlistVideosFromCtx(ctx);
  var name = escapeHtml(pl.name || "Playlist");
  var plId = pl.id || playlistId || "";
  var html = "<style>" + sidebarStyles + "</style>";
  html += '<div class="zt-playlist-up-next">';
  html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-list"></i></span>';
  html += "<h3>Up next</h3>";
  if (plId) html += '<p class="mb-2" style="margin-top: 20px; font-size: 1.0rem;"><a href="/playlist/' + escapeHtml(plId) + '">' + name + "</a></p>";
  html += "<hr /></div>";
  if (videos.length === 0) {
    html += '<p class="text-muted small">No videos in this playlist.</p>';
  } else {
    videos.forEach(function(v) {
      var vid = v.ID || v.id;
      var url = window.ztPlaylistPlayUrl(v, playlistId, {});
      var title = escapeHtml(v.Name || v.name || v.Filename || v.filename || "Untitled");
      var active = currentVideoId && vid === currentVideoId ? " active" : "";
      html += '<a class="zt-playlist-up-next-item' + active + '" href="' + url + '">';
      html += '<img class="lazy" data-src="' + thumbUrl(v) + '" alt="">';
      html += '<span class="zt-up-next-title">' + title + "</span></a>";
    });
  }
  html += "</div>";
  container.innerHTML = html;
  if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(container);
  var activeEl = container.querySelector(".zt-playlist-up-next-item.active");
  if (activeEl && activeEl.scrollIntoView) activeEl.scrollIntoView({ block: "nearest" });
};

window.ztPlaylistBindAutoAdvance = function(videoEl, ctx, playlistId) {
  if (!videoEl || !ctx) return;
  var ids = ctx.playlist_video_ids || [];
  var index = ctx.playlist_index;
  if (index == null || index < 0 || !ids.length) return;
  videoEl.addEventListener("ended", function() {
    var nextIdx = index + 1;
    if (nextIdx >= ids.length) return;
    var nextId = ids[nextIdx];
    if (window.ztPlaylistNavigateToId) {
      window.ztPlaylistNavigateToId(nextId, ctx, playlistId, { autoplay: true });
    }
  });
};

window.ztPlaylistShouldAutoplay = function() {
  return typeof window !== "undefined" && window.location.search.indexOf("autoplay") !== -1;
};

window.ztPlaylistTypeForId = function(ctx, videoId) {
  if (!ctx || !videoId) return "v";
  var items = ctx.playlist_items || [];
  for (var i = 0; i < items.length; i++) {
    var it = items[i];
    if ((it.id || it.ID) === videoId) return it.type || it.Type || "v";
  }
  var upNext = ctx.playlist_up_next || [];
  for (var j = 0; j < upNext.length; j++) {
    var v = upNext[j];
    if ((v.id || v.ID) === videoId) return v.type || v.Type || "v";
  }
  var cur = ctx.video;
  if (cur && (cur.id || cur.ID) === videoId) return cur.type || cur.Type || "v";
  return "v";
};

window.ztPlaylistNavigateToId = function(videoId, ctx, playlistId, opts) {
  opts = opts || {};
  if (!videoId || !playlistId) return;
  var t = window.ztPlaylistTypeForId(ctx, videoId);
  var path = (t === "c" ? "/clip/" : "/video/") + videoId;
  var q = "playlist=" + encodeURIComponent(playlistId);
  if (opts.autoplay) q += "&autoplay=1";
  var url = path + "?" + q;
  if (window.navigate) window.navigate(url);
  else window.location.href = url;
};

})();
