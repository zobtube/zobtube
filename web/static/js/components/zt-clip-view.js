(function() {
"use strict";
function ZtClipView() {
  var el = Reflect.construct(HTMLElement, [], ZtClipView);
  return el;
}
ZtClipView.prototype = Object.create(HTMLElement.prototype);
ZtClipView.prototype.disconnectedCallback = function() {
  if (this._clipKeyHandler) {
    document.removeEventListener("keyup", this._clipKeyHandler);
    this._clipKeyHandler = null;
  }
};
ZtClipView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = "Missing id"; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  var lastStatus = 0;
  fetch("/api/clip/" + encodeURIComponent(id), { credentials: "same-origin" })
    .then(function(r) { lastStatus = r.status; if (!r.ok) throw new Error(String(r.status)); return r.json(); })
    .then(function(data) {
      var v = data.video || data;
      var clipIds = data.clip_ids || [id];
      var streamUrl = data.stream_url || "/api/video/"+id+"/stream";
      var thumbUrl = "/api/video/"+id+"/thumb";
      var name = (v.Name||v.name||"").replace(/&/g,"&amp;").replace(/</g,"&lt;");
      var actors = v.Actors || v.actors || [];
      var categories = v.Categories || v.categories || [];
      var descParts = [];
      actors.forEach(function(a){ descParts.push("<b>@"+(a.Name||a.name||"")+"</b>"); });
      var seenCats = {};
      actors.forEach(function(a){ (a.Categories||a.categories||[]).forEach(function(c){ if (!seenCats[c.ID||c.id]) { seenCats[c.ID||c.id]=1; descParts.push("<b>#"+(c.Name||c.name||"")+"</b>"); } }); });
      categories.forEach(function(c){ if (!seenCats[c.ID||c.id]) { seenCats[c.ID||c.id]=1; descParts.push("<b>#"+(c.Name||c.name||"")+"</b>"); } });
      var descHtml = descParts.join(" ");
      var admin = (window.__USER__ && window.__USER__.admin);
      var editBtn = admin ? '<div style="width:40px;display:block;margin-top:25px;margin-left:auto;text-align:center"><i class="fas fa-pen clip-change" id="zt-clip-edit" style="font-size:24px;cursor:pointer"></i></div>' : '';
      var playSvg = '<svg style="transform:translate(5px,3px)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 50 50" width="50" height="50"><path fill="none" d="M0 0h50v50H0z"/><path d="m9 4.07 31.98 20.81c.02.07.02.19 0 .26L9 45.95V4.06M9.04 0C6.98 0 5 1.75 5 4.05v41.9C5 48.25 6.98 50 9.04 50c.68 0 1.36-.19 1.98-.61l32.2-20.95c2.36-1.53 2.36-5.35 0-6.88L11.02.61C10.4.19 9.71 0 9.04 0Z" fill="#fff"/></svg>';
      var arrowSvg = '<path d="M21.71 13.29a.996.996 0 0 0-1.41 0l-5.29 5.29V7c0-.55-.45-1-1-1s-1 .45-1 1v11.59L7.72 13.3a.996.996 0 1 0-1.41 1.41l7 7c.09.09.2.17.33.22.12.05.25.08.38.08s.26-.03.38-.08.23-.12.33-.22l7-7a.996.996 0 0 0 0-1.41Z" fill="currentColor"/>';
      var css = '#zt-clip-main{position:fixed;inset:0;background:#111;z-index:1000;display:flex;flex-direction:column;align-items:center}' +
        '#zt-clip-video-wrap{flex:1;width:100%;max-width:720px;display:flex;align-items:center;justify-content:center;position:relative}' +
        '#zt-clip-video{width:100%;height:100%;object-fit:contain}' +
        '.zt-clip-play{position:absolute;inset:0;opacity:.8;transition:opacity .4s;pointer-events:none;display:flex;align-items:center;justify-content:center;filter:drop-shadow(0 0 8px rgba(0,0,0,.8))}' +
        '.zt-clip-play.playing{opacity:0;pointer-events:none}' +
        '.zt-clip-play .play-inner{background:#1a1a1a;padding:20px;border-radius:50%;pointer-events:auto;cursor:pointer}' +
        '#zt-clip-nav{position:absolute;top:12px;right:0;width:100%;color:#fff}' +
        '#zt-clip-nav-inner{position:relative;margin:0 auto;max-width:720px}' +
        '#zt-clip-nav img,#zt-clip-nav svg{display:block;margin-left:auto}' +
        '.clip-change{cursor:pointer;color:#fff}.clip-change-disabled{color:#7e7e7e;cursor:default;pointer-events:none}' +
        '#zt-clip-details{position:absolute;bottom:56px;right:0;width:100%}' +
        '#zt-clip-details-inner{position:relative;margin:0 auto;max-width:720px;text-align:left;padding-left:5px}' +
        '#clip-title{color:#eee;text-shadow:1px 1px 4px #000;margin:0}' +
        '#clip-description{color:#eee;text-shadow:1px 1px 4px #000;margin:0}#clip-description b{color:#ececec}' +
        '.zt-clip-progress{position:absolute;left:0;right:0;bottom:5px;width:100%;height:5px;pointer-events:auto;padding:5px 0;cursor:pointer}' +
        '.zt-clip-progress-inner{max-width:720px;margin:0 auto}' +
        '.zt-clip-seek{width:100%;height:4px;background:hsla(0,0%,100%,.2);border-radius:4px;position:relative}' +
        '.zt-clip-seek-fill{height:100%;background:hsla(0,0%,100%,.6);border-radius:4px;transition:width .1s;position:relative}' +
        '.zt-clip-seek-fill:after{content:"";position:absolute;right:-6px;bottom:-6px;width:16px;height:16px;background:#fefefe;border-radius:20px}';
      var html = '<style>'+css+'</style>';
      html += '<div id="zt-clip-main">';
      html += '<div id="zt-clip-video-wrap">';
      html += '<video id="zt-clip-video" clip-id="'+id+'" src="'+streamUrl+'" preload="metadata" poster="'+thumbUrl+'" style="width:100%;height:100%;object-fit:contain"></video>';
      html += '<div id="play-button" class="zt-clip-play"><div class="play-inner">'+playSvg+'</div></div>';
      html += '</div>';
      html += '<div id="zt-clip-nav"><div id="zt-clip-nav-inner">';
      html += '<img src="/static/images/logo_clip.png" style="display:block;margin-left:auto;cursor:pointer" alt="Home">';
      html += '<svg id="clip-change-previous" class="clip-change-disabled" style="width:40px;display:block;margin-top:25px;margin-left:auto;transform:rotate(180deg)" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 28 28">'+arrowSvg+'</svg>';
      html += '<svg id="clip-change-next" class="clip-change" style="width:40px;display:block;margin-top:25px;margin-left:auto" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 28 28">'+arrowSvg+'</svg>';
      html += editBtn;
      html += '</div></div>';
      html += '<div id="zt-clip-details"><div id="zt-clip-details-inner">';
      html += '<div id="clip-title">'+name+'</div>';
      html += '<div id="clip-description">'+descHtml+'</div>';
      html += '</div></div>';
      html += '<div class="zt-clip-progress"><div class="zt-clip-progress-inner"><div class="zt-clip-seek"><div class="zt-clip-seek-fill" style="width:0%"></div></div></div></div>';
      html += '</div>';
      self.innerHTML = html;
      var video = self.querySelector("#zt-clip-video");
      var playBtn = self.querySelector("#play-button");
      var seekFill = self.querySelector(".zt-clip-seek-fill");
      var progressWrap = self.querySelector(".zt-clip-progress-inner");
      function togglePlay() {
        if (video.paused) { video.play(); playBtn.classList.add("playing"); }
        else { video.pause(); playBtn.classList.remove("playing"); }
      }
      if (playBtn) {
        playBtn.querySelector(".play-inner").addEventListener("click", togglePlay);
      }
      var viewCounted = false;
      var curIdx = clipIds.indexOf(id);
      var autoplay = typeof window !== "undefined" && window.location.search.indexOf("autoplay") !== -1;
      function updateNavButtons() {
        if (nextBtn) {
          nextBtn.style.pointerEvents = curIdx < clipIds.length - 1 ? "auto" : "none";
          nextBtn.setAttribute("class", curIdx < clipIds.length - 1 ? "clip-change" : "clip-change-disabled");
        }
        if (prevBtn) {
          prevBtn.style.pointerEvents = curIdx > 0 ? "auto" : "none";
          prevBtn.setAttribute("class", curIdx > 0 ? "clip-change" : "clip-change-disabled");
        }
      }
      function switchToClip(newId, shouldPlay) {
        var newIdx = clipIds.indexOf(newId);
        if (newIdx < 0) return;
        curIdx = newIdx;
        viewCounted = false;
        playBtn.classList.remove("playing");
        video.poster = "/api/video/"+newId+"/thumb";
        video.src = "/api/video/"+newId+"/stream";
        video.setAttribute("clip-id", newId);
        video.load();
        if (seekFill) seekFill.style.width = "0%";
        updateNavButtons();
        if (window.history && window.history.replaceState) history.replaceState({path: "/clip/"+newId}, "", "/clip/"+newId);
        fetch("/api/clip/"+encodeURIComponent(newId), { credentials: "same-origin" })
          .then(function(r){ return r.ok ? r.json() : Promise.reject(); })
          .then(function(data){
            var v = data.video || data;
            if (data.stream_url) {
              video.src = data.stream_url;
              video.load();
            }
            var name = (v.Name||v.name||"").replace(/&/g,"&amp;").replace(/</g,"&lt;");
            var actors = v.Actors || v.actors || [];
            var categories = v.Categories || v.categories || [];
            var descParts = [];
            actors.forEach(function(a){ descParts.push("<b>@"+(a.Name||a.name||"")+"</b>"); });
            var seenCats = {};
            actors.forEach(function(a){ (a.Categories||a.categories||[]).forEach(function(c){ if (!seenCats[c.ID||c.id]) { seenCats[c.ID||c.id]=1; descParts.push("<b>#"+(c.Name||c.name||"")+"</b>"); } }); });
            categories.forEach(function(c){ if (!seenCats[c.ID||c.id]) { seenCats[c.ID||c.id]=1; descParts.push("<b>#"+(c.Name||c.name||"")+"</b>"); } });
            var titleEl = self.querySelector("#clip-title");
            var descEl = self.querySelector("#clip-description");
            if (titleEl) titleEl.textContent = name || "";
            if (descEl) descEl.innerHTML = descParts.join(" ");
            if (editIcon) editIcon.onclick = function(){ window.navigate("/video/"+newId+"/edit"); };
            if (shouldPlay) video.play();
          });
      }
      if (video) {
        video.addEventListener("click", togglePlay);
        video.addEventListener("play", function(){
          var cid = video.getAttribute("clip-id");
          playBtn.classList.add("playing");
          if (!viewCounted) { viewCounted = true; fetch("/api/video/"+cid+"/count-view", {method:"POST",credentials:"same-origin"}); }
        });
        video.addEventListener("pause", function(){ playBtn.classList.remove("playing"); });
        video.addEventListener("ended", function(){
          if (curIdx < clipIds.length - 1) switchToClip(clipIds[curIdx+1], true);
        });
        video.addEventListener("timeupdate", function(){
          var d = video.duration;
          if (seekFill && !isNaN(d) && d > 0) seekFill.style.width = (100 * video.currentTime / d) + "%";
        });
        if (autoplay) video.play();
      }
      if (progressWrap) {
        progressWrap.addEventListener("click", function(e){
          var rect = this.getBoundingClientRect();
          var x = e.clientX - rect.left;
          if (x >= 0 && video && !isNaN(video.duration)) video.currentTime = (x / rect.width) * video.duration;
        });
      }
      var editIcon = self.querySelector("#zt-clip-edit");
      if (editIcon) editIcon.addEventListener("click", function(){ window.navigate("/video/"+id+"/edit"); });
      self.querySelector('img[alt="Home"]').addEventListener("click", function(){ window.location.href="/"; });
      var nextBtn = self.querySelector("#clip-change-next");
      var prevBtn = self.querySelector("#clip-change-previous");
      if (clipIds.length > 1) {
        updateNavButtons();
        if (nextBtn) nextBtn.addEventListener("click", function() { if (curIdx < clipIds.length - 1) switchToClip(clipIds[curIdx+1], false); });
        if (prevBtn) prevBtn.addEventListener("click", function() { if (curIdx > 0) switchToClip(clipIds[curIdx-1], false); });
        var wheelTimeout;
        function goNext() { if (curIdx < clipIds.length - 1) switchToClip(clipIds[curIdx+1], false); }
        function goPrev() { if (curIdx > 0) switchToClip(clipIds[curIdx-1], false); }
        function debounce(fn, ms) { clearTimeout(wheelTimeout); wheelTimeout = setTimeout(fn, ms || 200); }
        var mainEl = self.querySelector("#zt-clip-main");
        if (mainEl) {
          mainEl.addEventListener("wheel", function(e) {
            debounce(function() { (e.deltaY > 0 ? goNext : goPrev)(); });
            e.preventDefault();
          }, { passive: false });
          var touchY;
          mainEl.addEventListener("touchstart", function(ev){ touchY = ev.changedTouches[0].screenY; }, { passive: true });
          mainEl.addEventListener("touchend", function(ev){
            if (touchY == null) return;
            var y = ev.changedTouches[0].screenY;
            var d = touchY - y;
            touchY = null;
            if (Math.abs(d) > 100) { if (d > 0) goNext(); else goPrev(); }
          }, { passive: true });
        }
        var keyHandler = function(ev) {
          if (ev.key === "ArrowDown") { goNext(); ev.preventDefault(); }
          else if (ev.key === "ArrowUp") { goPrev(); ev.preventDefault(); }
        };
        document.addEventListener("keyup", keyHandler);
        self._clipKeyHandler = keyHandler;
      }
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function(err) {
      var status = lastStatus || (err && err.message);
      if (status === 401 || status === "401") {
        self.innerHTML = '<div class="alert alert-warning">Please <a href="/auth">sign in</a> to view this clip.</div>';
      } else {
        self.innerHTML = '<div class="alert alert-danger">Clip not found.</div>';
      }
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-clip-view", ZtClipView);
})();
