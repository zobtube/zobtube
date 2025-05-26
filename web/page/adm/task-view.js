{{ define "adm/task-view.js" }}
function restartTask(taskID) {
  $.ajax('/adm/task/'+taskID+'/retry', {
    method: 'POST',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },
    success: function (result) {
      window.location.reload();
    },
  });
}
{{ end }}
