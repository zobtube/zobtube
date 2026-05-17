(function() {
"use strict";
var sidebarStyles = "#zt-profile-sidebar{position:sticky;top:70px;padding:1rem 0;background:#f7f7f7;border-radius:8px}#zt-profile-sidebar .zt-profile-sidebar-title{margin:0 0 0.75rem;padding:0 1rem;font-size:1.1rem;font-weight:600}#zt-profile-sidebar .zt-profile-sidebar-group{margin:1rem 0 0.25rem;padding:0 1rem;font-size:0.7rem;font-weight:600;text-transform:uppercase;color:#6c757d}#zt-profile-sidebar .zt-profile-sidebar-group:first-of-type{margin-top:0}#zt-profile-sidebar .nav{flex-direction:column}#zt-profile-sidebar .nav-link{padding:0.5rem 1rem;border-radius:4px;border-left:3px solid transparent}#zt-profile-sidebar .nav-link:hover{background:rgba(0,0,0,0.05)}#zt-profile-sidebar .nav-link.active{background:rgba(22,122,198,0.1);border-left-color:#167ac6;color:#167ac6}";
function ZtProfileTabs() {
  var el = Reflect.construct(HTMLElement, [], ZtProfileTabs);
  return el;
}
ZtProfileTabs.prototype = Object.create(HTMLElement.prototype);
ZtProfileTabs.prototype.connectedCallback = function() {
  var active = this.getAttribute("data-active") || "";
  function link(label, href, tab) {
    var c = tab === active ? "nav-link active" : "nav-link";
    return '<a class="' + c + '" href="' + (tab === active ? "#" : href) + '">' + label + '</a>';
  }
  function group(title) {
    return '<div class="zt-profile-sidebar-group">' + title + '</div>';
  }
  this.innerHTML = '<style>' + sidebarStyles + '</style><div id="zt-profile-sidebar"><h3 class="zt-profile-sidebar-title"><span class="heading-icon"><i class="fa fa-user"></i></span> Your account</h3><nav class="nav flex-column">' +
    group("Most viewed") +
    link("Videos", "/profile/most-viewed/videos", "most-viewed-videos") +
    link("Actors", "/profile/most-viewed/actors", "most-viewed-actors") +
    group("Security") +
    link("Password", "/profile/settings", "settings") +
    link("API tokens", "/profile/tokens", "tokens") +
    '</nav></div>';
};
customElements.define("zt-profile-tabs", ZtProfileTabs);
})();
