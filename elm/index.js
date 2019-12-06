"use strict";

require("./index.html");
require("./src/css/reset.css");
require("./src/sass/styles.scss");
var Elm = require("./src/Main.elm");

if (navigator.mediaSession && navigator.mediaSession.setActionHandler) {
  navigator.mediaSession.setActionHandler("previoustrack", () =>
    app.ports.audioPrevPressed.send(null)
  );

  navigator.mediaSession.setActionHandler("nexttrack", () =>
    app.ports.audioNextPressed.send(null)
  );

  navigator.mediaSession.setActionHandler("seekbackward", () => {
    console.log("seekbackward");
  });

  navigator.mediaSession.setActionHandler("seekforward", () => {
    console.log("seekforward");
  });
}

// Start elm with the possible serialised model from local storage
var stored = localStorage.getItem("sound-ui-elm");
var model = stored ? JSON.parse(stored) : null;
var app = Elm.Elm.Main.init({ flags: { config: SOUND_CONFIG, model } });

// Port for serialising and storing the elm model
app.ports.setCache.subscribe(model =>
  localStorage.setItem("sound-ui-elm", JSON.stringify(model))
);

const audios = new Map();

let socket;
window.app = app;

app.ports.playAudio.subscribe(songId => {
  //console.log(`playAudio ${songId}`);
  const audio = audios.get(songId);
  if (audio) {
    audio.pause();
    audio.currentTime = 0;
    audio.play();
  }
});

app.ports.pauseAudio.subscribe(songId => {
  //console.log(`pause ${songId}`);
  const audio = audios.get(songId);
  if (audio) {
    audio.pause();
  }
});

app.ports.resumeAudio.subscribe(songId => {
  //console.log(`resume ${songId}`);
  const audio = audios.get(songId);
  if (audio) {
    audio.play();
  }
});

app.ports.loadAudio.subscribe(({ url, songId }) => {
  //console.log(`loadAudio ${songId}`);

  var a = new Audio(url);

  a.oncanplay = () => {
    //console.log(`oncanplay ${songId}`);
    app.ports.canPlayAudio.send(songId);
    a.oncanplay = null;
  };

  a.ondurationchange = () => {
    // console.log(`ondurationchange ${songId}`);
    // app.ports.audioDurationChanged.send({ songId, duration: a.duration });
  };

  a.ontimeupdate = () => {
    // console.log(`ontimeupdate ${songId}`);
    app.ports.audioTimeChanged.send({ songId, time: a.currentTime });
  };

  a.onended = () => {
    // console.log(`onended ${songId}`);
    app.ports.audioEnded.send(songId);
  };

  a.onplay = () => {
    //console.log(`onplay ${songId}`);

    if (navigator.mediaSession) {
      navigator.mediaSession.metadata = new MediaMetadata({
        title: `Unforgettable track #${songId}`,
        artist: "Nat King Cole",
        album: "The Ultimate Collection (Remastered)",
        artwork: [
          {
            src: "https://dummyimage.com/96x96",
            sizes: "96x96",
            type: "image/png"
          },
          {
            src: "https://dummyimage.com/128x128",
            sizes: "128x128",
            type: "image/png"
          },
          {
            src: "https://dummyimage.com/192x192",
            sizes: "192x192",
            type: "image/png"
          },
          {
            src: "https://dummyimage.com/256x256",
            sizes: "256x256",
            type: "image/png"
          },
          {
            src: "https://dummyimage.com/384x384",
            sizes: "384x384",
            type: "image/png"
          },
          {
            src: "https://dummyimage.com/512x512",
            sizes: "512x512",
            type: "image/png"
          }
        ]
      });
    }

    app.ports.audioPlaying.send({
      songId,
      time: a.currentTime,
      duration: a.duration
    });
  };

  a.onpause = () => {
    //console.log(`onpause ${songId}`);
    app.ports.audioPaused.send({
      songId,
      time: a.currentTime,
      duration: a.duration
    });
  };

  audios.set(songId, a);
});

app.ports.setAudioTime.subscribe(({ songId, time }) => {
  const audio = audios.get(songId);
  if (audio) {
    audio.currentTime = time;
  }
});

// Port for creating a websocket from elm
app.ports.websocketOpen.subscribe(url => {
  socket = new WebSocket(url);

  // Tell elm the socket is open
  socket.onopen = () => {
    // Let elm know when port closes
    socket.onclose = () => {
      app.ports.websocketClosed.send(null);
    };

    socket.onerror = () => {
      app.ports.websocketClosed.send(null);
    };

    app.ports.websocketOpened.send(null);
  };

  // Forward incoming messages to elm
  socket.onmessage = message => {
    app.ports.websocketIn.send(message.data);
  };
});

// Port for sending messages through the socket from elm
app.ports.websocketOut.subscribe(message => {
  if (socket && socket.readyState === 1) {
    socket.send(JSON.stringify(message));
  }
});

// Port for closing the socket from elm
app.ports.websocketClose.subscribe(() => {
  if (socket == null) {
    return;
  }

  // Don't tell elm about a close which it instigated
  socket.onclose = null;

  socket.close();
  socket = null;
});

/*
fetch(SOUND_CONFIG.root + "/api/authenticate", {
  method: "POST", // *GET, POST, PUT, DELETE, etc.
  mode: "cors", // no-cors, *cors, same-origin
  cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
  credentials: "include", // include, *same-origin, omit
  headers: {
    "Content-Type": "application/json"
    // 'Content-Type': 'application/x-www-form-urlencoded',
  },
  redirect: "follow", // manual, *follow, error
  referrer: "no-referrer", // no-referrer, *client
  body: JSON.stringify({ username: "gertrude", password: "changeme" }) // body data type must match "Content-Type" header
});
*/
