(function() {
"use strict";
var sidebarStyles = "#zt-adm-sidebar{position:sticky;top:70px;padding:1rem 0;background:#f7f7f7;border-radius:8px}#zt-adm-sidebar .zt-adm-sidebar-title{margin:0 0 0.75rem;padding:0 1rem;font-size:1.1rem;font-weight:600}#zt-adm-sidebar .zt-adm-sidebar-group{margin:1rem 0 0.25rem;padding:0 1rem;font-size:0.7rem;font-weight:600;text-transform:uppercase;color:#6c757d}#zt-adm-sidebar .zt-adm-sidebar-group:first-of-type{margin-top:0}#zt-adm-sidebar .nav{flex-direction:column}#zt-adm-sidebar .nav-link{padding:0.5rem 1rem;border-radius:4px;border-left:3px solid transparent}#zt-adm-sidebar .nav-link:hover{background:rgba(0,0,0,0.05)}#zt-adm-sidebar .nav-link.active{background:rgba(22,122,198,0.1);border-left-color:#167ac6;color:#167ac6}";
function ZtAdmTabs() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmTabs);
  return el;
}
ZtAdmTabs.prototype = Object.create(HTMLElement.prototype);
ZtAdmTabs.prototype.connectedCallback = function() {
  var active = this.getAttribute("data-active") || "";
  function link(label, href, tab) {
    var c = tab === active ? "nav-link active" : "nav-link";
    return '<a class="' + c + '" href="' + (tab === active ? "#" : href) + '">' + label + '</a>';
  }
  function group(title) {
    return '<div class="zt-adm-sidebar-group">' + title + '</div>';
  }
  this.innerHTML = '<style>' + sidebarStyles + '</style><div id="zt-adm-sidebar"><h3 class="zt-adm-sidebar-title"><span class="heading-icon"><i class="fa fa-hammer"></i></span> Administration</h3><nav class="nav flex-column">' +
    group("Main") +
    link("Overview", "/adm", "overview") +
    link("Videos", "/adm/videos", "videos") +
    link("Actors", "/adm/actors", "actors") +
    link("Channels", "/adm/channels", "channels") +
    link("Categories", "/adm/categories", "categories") +
    group("Authentication") +
    link("General", "/adm/config/auth", "authentication") +
    link("Users", "/adm/users", "users") +
    link("API tokens", "/adm/tokens", "tokens") +
    group("External") +
    link("Offline mode", "/adm/config/offline", "offline") +
    link("Providers", "/adm/config/provider", "providers") +
    group("Internal") +
    link("Tasks", "/adm/task/home", "tasks") +
    '</nav></div>';
};
customElements.define("zt-adm-tabs", ZtAdmTabs);
})();
