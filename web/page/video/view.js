{{ define "video/view.js" }}

viewCounted = false;
function countView() {
  if (viewCounted) {
    return
  }
  $.ajax('/api/video/{{ .Video.ID }}/count-view', {
    method: 'POST',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },
    success: function(data) {
      $('#video-view-count').animate({'opacity': 0}, 400, function(){
        viewCount = data['view-count'];
        newViewCount = '<h3>'+viewCount+' view';
        newViewCount += viewCount > 1 ? 's' : '';
        newViewCount += '</h3>';
        $(this).html(newViewCount).animate({'opacity': 1}, 400);
      });
      viewCounted = true;
    },
  });
}

{{ end }}
