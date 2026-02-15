(function() {
"use strict";
function ZtAdmTabs() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmTabs);
  return el;
}
ZtAdmTabs.prototype = Object.create(HTMLElement.prototype);
ZtAdmTabs.prototype.connectedCallback = function() {
  var active = this.getAttribute("data-active") || "";
  function link(label, href, tab) {
    var c = tab === active ? "nav-link active" : "nav-link";
    return '<li class="nav-item"><a class="' + c + '" href="' + (tab === active ? "#" : href) + '">' + label + '</a></li>';
  }
  this.innerHTML = '<div class="col-md-12"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-hammer"></i></span><h3>Administration</h3><hr /></div></div>' +
    '<div class="col-md-12 mb-4"><ul class="nav nav-tabs">' +
    link("Overview", "/adm", "overview") +
    link("Authentication", "/adm/config/auth", "authentication") +
    link("Providers", "/adm/config/provider", "providers") +
    link("Tasks", "/adm/task/home", "tasks") +
    link("Offline mode", "/adm/config/offline", "offline") +
    '</ul></div>';
};
customElements.define("zt-adm-tabs", ZtAdmTabs);
})();
