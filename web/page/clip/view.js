{{ define "clip/view.js" }}

function videoStartStop() {
  video = document.getElementById('video-clip');

  if (video.paused) {
    video.play();
  } else {
    video.pause();
  }
}

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
      viewCounted = true;
    },
  });
}

function nextVideo() {
  video = document.getElementById('video-clip');
  currentID = video.getAttribute('clip-id');

  var nextClipID;
  var i;
  for (var i=0; i<clipList.length; i++) {
    if (clipList[i] == currentID) {
      if (i + 1 == clipList.length) {
        console.debug('end of clip list');
        return
      }

      nextClipID = clipList[i+1];
      break;
    }
  }

  if (i + 2 == clipList.length) {
    nextButton = document.getElementById('clip-change-next');
    nextButton.classList.add('clip-change-disabled');
    nextButton.classList.remove('clip-change');
  }

  previousButton = document.getElementById('clip-change-previous');
  previousButton.classList.remove('clip-change-disabled');
  previousButton.classList.add('clip-change');

  setVideoByID(nextClipID);
}

function previousVideo() {
  video = document.getElementById('video-clip');
  currentID = video.getAttribute('clip-id');

  var nextClipID;
  var i;
  for (i=0; i<clipList.length; i++) {
    if (clipList[i] == currentID) {
      if (i - 1 < 0) {
        console.debug('beginning of clip list');
        return
      }

      nextClipID = clipList[i-1];
      break;
    }
  }

  if (i-1 == 0) {
    previousButton = document.getElementById('clip-change-previous');
    previousButton.classList.add('clip-change-disabled');
    previousButton.classList.remove('clip-change');
  }

  nextButton = document.getElementById('clip-change-next');
  nextButton.classList.remove('clip-change-disabled');
  nextButton.classList.add('clip-change');

  setVideoByID(nextClipID);
}

function setVideoByID(id) {
  var paused = video.paused;
  video.poster = '/video/'+id+'/thumb';
  video.src = '/video/'+id+'/stream';
  video.setAttribute('clip-id', id);
  video.preload = 'metadata';
  video.preload = 'auto';
  if (!paused) {
    video.play();
  }

  // get info async
  $.ajax('/api/video/'+id, {
    method: 'GET',
    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },
    success: function(data) {
      titleDiv = document.getElementById('clip-title');
      titleDiv.innerText = data.title;

      descDiv = document.getElementById('clip-description');
      description = '';
      if (data.actors != null) {
        for (const actor of data.actors) {
          description += '<b>@'+actor+'</b> ';
        }
      }
      if (data.categories != null) {
        for (const category of data.categories) {
          description += '<b>#'+category+'</b> ';
        }
      }
      descDiv.innerHTML = description;
    },
  });

}

function goHome() {
  window.location = '/';
}

function redirectToClipEdition() {
  // get video item
  video = document.getElementById('video-clip');

  // get video id
  currentID = video.getAttribute('clip-id');

  // redirect
  window.location = '/video/'+currentID+'/edit';
}

document.body.addEventListener('wheel', checkScrollDirection);

function checkScrollDirection(event) {
  if (checkScrollDirectionIsUp(event)) {
    debouncer(function() {previousVideo();});
  } else {
    debouncer(function() {nextVideo();});
  }
}

function checkScrollDirectionIsUp(event) {
  if (event.wheelDelta) {
    return event.wheelDelta > 0;
  }
  return event.deltaY < 0;
}

let timeoutID;
function debouncer(func, timeout) {
  var timeout = timeout || 200;
  var scope = this , args = arguments;
  clearTimeout( timeoutID );
  console.log("set timeout");
  timeoutID = setTimeout( function () {
    func.apply( scope , Array.prototype.slice.call( args ) );
  } , timeout );
}

{{ end }}
