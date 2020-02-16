const socket = io.connect('http://localhost:8080/input');

socket.on('player_registered', function (playerId) {
    console.log("player registered. playerId: " + playerId);
    createSendCommandBox(playerId)
});

function createSendCommandBox(playerId) {
    let sessionCode = $('#sessionCode').val();

    var html = `
    <div>
        <h2>Player Id: ` + playerId + `</h2>
        <table id="playersTable" style="width:100%">
            <thead>
                <tr>
                    <th>X</th>
                    <th>Y</th>
                    <th>Active Actions</th>
                    <th>Send Comannd</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td><input id="x_` + playerId + `" type="text"/></td>
                    <td><input id="y_` + playerId + `" type="text"/></td>
                    <td><input id="actice_actions_` + playerId +`" type="text"/></td>
                    <td><input type="button" onclick="sendCommandOnClick('` + playerId + `')" value="Send Comand!"/></td>
                </tr>
            </tbody>
        </table>
        <div>
            <input id="playerId_` + playerId + `" type="hidden" value="` + playerId + `"/>
            <input id="sessionCode_sendCommand_` + playerId + `" type="hidden" value="` + sessionCode + `"/>
        </div>
    </div>
    
    `;

    $('#commands').append(html);
}

function registerPlayer(sessionCode, playerName) {
    socket.emit('register_player', '{"sessionCode": "' + sessionCode + '", "playerName": "' + playerName +'"}');
}

function registrePlayerOnClick() {

    let sessionCode = $('#sessionCode').val()
    let playerName = $('#playerName').val()

    registerPlayer(sessionCode, playerName)
}


function sendCommand(playerId, sessionCode, x, y, activeActions) {
    let command = '{"playerId": "' + playerId + '", "sessionCode": "' + sessionCode +'", "activeActions": [' + activeActions + '], "direction": { "x": ' + x + ', "y": ' + y + '}}';
    socket.emit('context', command);
}

function sendCommandOnClick(index) {
    let sessionCode = $('#sessionCode_sendCommand_' + index).val();
    let playerId = $('#playerId_' + index).val();

    let x = $('#x_' + index).val();
    let y = $('#y_' + index).val();
    let activeActions = $('#actice_actions_' + index).val();

    sendCommand(playerId, sessionCode, x, y, activeActions)
}

