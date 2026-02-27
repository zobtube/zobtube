(function() {
"use strict";

var routes = [
  { pattern: /^\/$/, component: "zt-home" },
  { pattern: /^\/actors\/?$/, component: "zt-actor-list" },
  { pattern: /^\/actor\/new\/?$/, component: "zt-actor-create" },
  { pattern: /^\/actor\/([^\/]+)\/edit\/?$/, component: "zt-actor-edit", param: "id" },
  { pattern: /^\/actor\/([^\/]+)\/?$/, component: "zt-actor-view", param: "id" },
  { pattern: /^\/categories\/?$/, component: "zt-category-list" },
  { pattern: /^\/category\/([^\/]+)\/?$/, component: "zt-category-view", param: "id" },
  { pattern: /^\/channels\/?$/, component: "zt-channel-list" },
  { pattern: /^\/channel\/new\/?$/, component: "zt-channel-create" },
  { pattern: /^\/channel\/([^\/]+)\/edit\/?$/, component: "zt-channel-edit", param: "id" },
  { pattern: /^\/channel\/([^\/]+)\/?$/, component: "zt-channel-view", param: "id" },
  { pattern: /^\/clips\/?$/, component: "zt-clip-list" },
  { pattern: /^\/movies\/?$/, component: "zt-video-list", type: "movie" },
  { pattern: /^\/videos\/?$/, component: "zt-video-list", type: "video" },
  { pattern: /^\/video\/([^\/]+)\/edit\/?$/, component: "zt-video-edit", param: "id" },
  { pattern: /^\/video\/([^\/]+)\/?$/, component: "zt-video-view", param: "id" },
  { pattern: /^\/clip\/([^\/]+)\/?$/, component: "zt-clip-view", param: "id" },
  { pattern: /^\/profile\/tokens\/?$/, component: "zt-profile-tokens" },
  { pattern: /^\/profile\/settings\/?$/, component: "zt-profile-settings" },
  { pattern: /^\/profile\/?$/, component: "zt-profile-view" },
  { pattern: /^\/upload\/?$/, component: "zt-upload-home" },
  { pattern: /^\/adm\/?$/, component: "zt-adm-home" },
  { pattern: /^\/adm\/videos\/?$/, component: "zt-adm-object-list", object: "Video" },
  { pattern: /^\/adm\/actors\/?$/, component: "zt-adm-object-list", object: "Actor" },
  { pattern: /^\/adm\/channels\/?$/, component: "zt-adm-object-list", object: "Channel" },
  { pattern: /^\/adm\/categories\/?$/, component: "zt-adm-category" },
  { pattern: /^\/adm\/libraries\/?$/, component: "zt-adm-library-list" },
  { pattern: /^\/adm\/task\/home\/?$/, component: "zt-adm-task-home" },
  { pattern: /^\/adm\/tasks\/?$/, component: "zt-adm-task-list" },
  { pattern: /^\/adm\/task\/([^\/]+)\/?$/, component: "zt-adm-task-view", param: "id" },
  { pattern: /^\/adm\/user\/?$/, component: "zt-adm-user-new" },
  { pattern: /^\/adm\/users\/?$/, component: "zt-adm-user-list" },
  { pattern: /^\/adm\/tokens\/?$/, component: "zt-adm-token-list" },
  { pattern: /^\/adm\/config\/auth\/?$/, component: "zt-adm-config-auth" },
  { pattern: /^\/adm\/config\/provider\/?$/, component: "zt-adm-config-provider" },
  { pattern: /^\/adm\/config\/offline\/?$/, component: "zt-adm-config-offline" }
];

function matchRoute(path) {
  path = path || "/";
  if (path === "") path = "/";
  path = path.replace(/\/$/, "") || "/";
  for (var i = 0; i < routes.length; i++) {
    var r = routes[i];
    var m = path.match(r.pattern);
    if (m) {
      var params = {};
      if (r.param && m[1]) params[r.param] = m[1];
      if (r.object) params.object = r.object;
      if (r.type) params.type = r.type;
      return { component: r.component, params: params };
    }
  }
  return null;
}

var loadId = 0;
var loadTimeout = null;
var lastFinishedLoadId = 0;

function showNavProgress() {
  var bar = document.getElementById("zt-nav-progress");
  if (bar) bar.classList.add("active");
}

function hideNavProgress() {
  var bar = document.getElementById("zt-nav-progress");
  if (bar) bar.classList.remove("active");
}

function finishLoad(id, el) {
  if (id !== loadId || !el) return;
  if (id === lastFinishedLoadId) return;
  lastFinishedLoadId = id;
  var app = document.getElementById("app");
  var staging = document.getElementById("app-staging");
  if (!app || !staging) return;
  if (loadTimeout) { clearTimeout(loadTimeout); loadTimeout = null; }
  hideNavProgress();
  app.innerHTML = "";
  app.appendChild(el);
  staging.innerHTML = "";
  if (window.zt && window.zt.onload && window.zt.onload.length) {
    window.zt.onload.forEach(function(fn) { try { fn(); } catch (e) {} });
  }
  if (window.lazyLoadInstance) window.lazyLoadInstance.update();
  setTimeout(function() {
    if (window.lazyLoadInstance) window.lazyLoadInstance.update();
  }, 200);
}

function handlePageReady(el) {
  if (!el || !el.ztLoadId) return;
  if (el.ztLoadId !== loadId) return;
  var staging = document.getElementById("app-staging");
  if (!staging || !staging.contains(el)) return;
  finishLoad(el.ztLoadId, el);
}

document.addEventListener("zt-page-ready", function(e) {
  handlePageReady(e.target);
});

if (window.zt) {
  window.zt.pageReady = function(el) {
    if (el) handlePageReady(el);
  };
  window.zt.loadLazyIn = function(container) {
    if (!container) return;
    var els = container.querySelectorAll(".lazy");
    for (var i = 0; i < els.length; i++) {
      var img = els[i];
      var src = img.getAttribute("data-src");
      if (src && img.tagName === "IMG") { img.src = src; img.classList.remove("lazy"); img.classList.add("loaded"); }
    }
  };
}

function loadPage(path) {
  var app = document.getElementById("app");
  var staging = document.getElementById("app-staging");
  if (!app) return;
  path = path || window.location.pathname || "/";
  if (path === "") path = "/";

  if (loadTimeout) { clearTimeout(loadTimeout); loadTimeout = null; }

  var match = matchRoute(path);
  if (!match) {
    hideNavProgress();
    app.innerHTML = '<div class="container"><p class="text-muted">Page not found.</p><a href="/">Home</a></div>';
    if (staging) staging.innerHTML = "";
    return;
  }

  showNavProgress();
  if (staging) staging.innerHTML = "";
  loadId += 1;
  var currentLoadId = loadId;

  try {
    var el = document.createElement(match.component);
    el.ztLoadId = currentLoadId;
    for (var k in match.params) {
      if (match.params.hasOwnProperty(k)) {
        el.setAttribute("data-" + k, match.params[k]);
      }
    }
    if (staging) {
      staging.appendChild(el);
      loadTimeout = setTimeout(function() {
        if (currentLoadId === loadId && staging.contains(el)) {
          finishLoad(currentLoadId, el);
        }
      }, 15000);
    } else {
      app.innerHTML = "";
      app.appendChild(el);
      hideNavProgress();
      if (window.zt && window.zt.onload && window.zt.onload.length) {
        window.zt.onload.forEach(function(fn) { try { fn(); } catch (e) {} });
      }
    }
  } catch (e) {
    hideNavProgress();
    if (staging) staging.innerHTML = "";
    app.innerHTML = '<div class="container"><p class="text-muted">Failed to load page.</p><a href="/">Home</a></div>';
  }
}

function navigate(url) {
  var path = url;
  try {
    var a = document.createElement("a");
    a.href = url;
    path = a.pathname || "/";
  } catch (e) {}
  if (path === "") path = "/";
  history.pushState({ path: path }, "", url);
  loadPage(path);
}

window.loadPage = loadPage;
window.navigate = navigate;

document.addEventListener("click", function(e) {
  var a = e.target.closest("a[href]");
  if (!a) return;
  if (a.target === "_blank" || a.hasAttribute("download")) return;
  var href = a.getAttribute("href");
  if (!href || href === "#") return;
  if (href.indexOf("mailto:") === 0 || href.indexOf("http:") === 0 || href.indexOf("https:") === 0) {
    if (a.origin !== window.location.origin) return;
  }
  if (href.indexOf("/") === 0 && a.origin === window.location.origin) {
    if (a.hasAttribute("data-full-reload")) return;
    e.preventDefault();
    navigate(a.href);
  }
});

window.addEventListener("popstate", function() {
  loadPage(window.location.pathname || "/");
});

})();
