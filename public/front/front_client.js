const socket = io.connect('http://localhost:8080/front');

function printSession(sessionJson) {
    let session = JSON.parse(sessionJson)
    console.log("session response: " + sessionJson)
    console.log("session code: " + session.sessionCode);
    $('#code').text(session.sessionCode)
    $('#sessionCode').text(session.sessionCode)
    $('#socketId').text(session.socketId)

    $('#playersTable tbody tr').remove()
    session.players.forEach(function(element, index) {
        $('#playersTable').append("<tr><td>" + element.id + "</td><td>" + element.playerName + "</td><tr>")
    });
}

function printInputContext(inputContextJson) {

    let inputContex = JSON.parse(inputContextJson)

    var html = `
<div>
    <table style="width:100%">
        <thead>
            <tr>
                <th>Player Id</th>
                <th>Session Code</th>
                <th>X</th>
                <th>Y</th>
                <th>Active Actions</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>` + inputContex.playerId + `</td>
                <td>` + inputContex.sessionCode + `</td>
                <td>` + inputContex.direction.x + `</td>
                <td>` + inputContex.direction.y + `</td>
                <td>` + inputContex.activeActions + `</td>
            </tr>
        </tbody>
    </table>
</div>

`;

    $('#commands').append(html);
}

socket.on('session_created', function (sessionJson) {
    printSession(sessionJson)
});

socket.on('session_changed', function (sessionJson) {
    printSession(sessionJson)
});

socket.on('input_context', function (inputContext) {
    printInputContext(inputContext)
});

socket.emit('create_session', '');