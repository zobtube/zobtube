(function() {
"use strict";
function ZtAdmConfigAuth() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmConfigAuth);
  return el;
}
ZtAdmConfigAuth.prototype = Object.create(HTMLElement.prototype);
ZtAdmConfigAuth.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/config/auth", { credentials: "same-origin" })
    .then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
    .then(function(d) {
      var enabled = d.authentication_enabled === true;
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="authentication"></zt-adm-tabs></div><div class="col-md-9 col-lg-9">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-lock"></i></span><h3>General</h3><hr /></div>';
      html += '<div class="row"><div class="col-md-12 mb-4"><h4 class="mb-4">Authentication settings</h4>';
      html += 'The authentication is currently: ';
      html += enabled ? '<span class="badge text-bg-success">Enabled</span>' : '<span class="badge text-bg-warning">Disabled</span>';
      html += '</div>';
      html += '<div class="col-md-12 mb-4"><h4 class="mb-4">Change setting</h4>';
      html += '<p>To ' + (enabled ? 'disable' : 'enable') + ' authentication, please read the following disclaimer.</p>';
      if (enabled) {
        html += '<div class="alert alert-warning" role="alert">Disabling the authentication will remove the possibility to change user. Only the first admin user created will be available and always auto-logged-in.</div>';
        html += '<button class="btn btn-xl btn-danger" id="zt-auth-disable">Disable authentication</button>';
      } else {
        html += '<div class="alert alert-warning" role="alert">Once authentication is enabled, login will be required at all time. If your account does not have a password yet, you can run the password reset command. <code>zobtube password-reset</code></div>';
        html += '<button class="btn btn-xl btn-danger" id="zt-auth-enable">Enable authentication</button>';
      }
      html += '</div></div>';
      html += '</div></div>';
      self.innerHTML = html;
      var btn = self.querySelector(enabled ? "#zt-auth-disable" : "#zt-auth-enable");
      if (btn) {
        btn.addEventListener("click", function() {
          var action = enabled ? "disable" : "enable";
          fetch("/api/adm/config/auth/" + action, { method: "GET", credentials: "same-origin" })
            .then(function(r) {
              if (r.ok) {
                if (typeof loadPage === "function") loadPage("/adm/config/auth");
                else window.location.reload();
              }
            });
        });
      }
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function(e) {
      if (e && e.message === "403") self.innerHTML = '<div class="alert alert-danger">Forbidden</div>';
      else self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-adm-config-auth", ZtAdmConfigAuth);
})();
