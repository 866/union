var socket;

$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws');
    // Message received on the socket
    socket.onmessage = function (event) {
        $('#json-renderer').jsonViewer(JSON.parse(event.data), {collapsed: false, withQuotes: false});
    };
});
