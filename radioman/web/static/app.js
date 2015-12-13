$(document).ready(function() {
  var player = $('#player');
  player.hide();
  $.ajax('/api/radios/default/endpoints').done(function(data) {
    for (var idx in data.endpoints) {
      var endpoint = data.endpoints[idx];
      var source = $('<source></source>').attr('src', endpoint.source);
      source.appendTo(player);
    }
    player.fadeIn(1000);
    setTimeout(function() {
      player[0].play();
    }, 100);
  });
});
