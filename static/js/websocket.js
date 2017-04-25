var socket;

var receive = true

$(document).ready(function () {
    // Create a socket
    socket = new WebSocket('ws://' + window.location.host + '/ws');
    // Message received on the socket
    socket.onmessage = function (event) {
        if (receive) {
            $('#json-renderer').jsonViewer(JSON.parse(event.data), {collapsed: false, withQuotes: false});
        };
    };
});

function wsSwitch() {
    receive = !receive;
};
