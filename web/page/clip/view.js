{{ define "clip/view.js" }}

function videoStartStop() {
  video = document.getElementById('video-clip');

  if (video.paused) {
    video.play();
    document.getElementById('play-button').style.opacity = 0;
  } else {
    video.pause();
    document.getElementById('play-button').style.opacity = 0.8;
  }
}
function videoStartStopDebounced() {
    debouncer(function() {videoStartStop();});
}

document.getElementById('video-clip').addEventListener('click', videoStartStopDebounced);
document.getElementById('video-clip').addEventListener(
  "timeupdate",
  function(event){
    onTrackedVideoFrame(this.currentTime);
  }
);



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

  // reset progress bar
  document.getElementsByClassName('ProgressBar_ProgressBar')[0].style.width='0%';

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

// wheel detection
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
  timeoutID = setTimeout( function () {
    func.apply( scope , Array.prototype.slice.call( args ) );
  } , timeout );
}

// keyboard down / up detection
document.addEventListener("keyup", event => {
  if (event.keyCode === 40) {
    nextVideo();
    return;
  }

  if (event.keyCode === 38) {
    previousVideo();
    return;
  }
});

// finger swipe detection
// from: https://stackoverflow.com/a/23230280
var yDown = null;

function handleTouchStart(evt) {
  yDown = evt.changedTouches[0].screenY;
};

function handleTouchMove(evt) {
  var yUp = evt.changedTouches[0].screenY;
  var yDiff = yDown - yUp;

  if ( Math.abs(yDiff) < 100 ) {
    // small mouvement: pause
    videoStartStopDebounced();
  } else {
    if ( yDiff > 0 ) {
      /* up swipe */
      nextVideo();
    } else {
      /* down swipe */
      previousVideo();
    }
  }
  /* reset values */
  yDown = null;
};

document.addEventListener('touchstart', handleTouchStart, false);
document.addEventListener('touchend', handleTouchMove, false);

// progression update
function onTrackedVideoFrame(currentTime){
  duration = document.getElementById('video-clip').duration;
  if (duration === NaN) {
    duration = 0;
  }
  progression = currentTime * 100 / duration;
  document.getElementsByClassName('ProgressBar_ProgressBar')[0].style.width=progression+'%';
}

// progression click
document.getElementsByClassName('ProgressBar_ClickWrapper')[0].addEventListener('click', function (e) {
  var x = e.pageX - this.offsetLeft;
  if ( x <= 0 ) {
    x = 0
  }

  progress = x/this.clientWidth;
  videoClip = document.getElementById('video-clip');
  videoClip.currentTime = videoClip.duration * progress;
});

{{ end }}
