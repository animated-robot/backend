function startSocket() {

    this.window.socket = io.connect(parseSocketIp(), {
        transports: ['websocket']
    });

    this.window.socket.on('player_registered', function (playerId) {
        console.log("player registered. playerId: " + playerId);
        createSendCommandBox(playerId)
    });
}


function sendCommandOnClick(index) {
    const context = getInputContext(index);
    sendCommand(context.playerId, context.sessionCode, context.x, context.y, context.activeActions)
}

function sendCommand(playerId, sessionCode, x, y, activeActions) {
    const direction  = newDirection(x, y);
    const context = newContext(sessionCode, activeActions, playerId, direction);
    this.window.socket.emit('context', JSON.stringify(context));
}


function registerPlayer(sessionCode, playerName) {
    this.window.socket.emit('register_player', createRegisterPlayerPayload(sessionCode, playerName));
}

function createRegisterPlayerPayload(sessionCode, playerName) {
    const payload = {
        sessionCode: sessionCode,
        player: {
            name: playerName
        }
    };

    return JSON.stringify(payload)
}

function registerPlayerOnClick() {

    let sessionCode = $('#sessionCode').val();
    let playerName = $('#playerName').val();

    registerPlayer(sessionCode, playerName)
}



function createSendCommandBox(playerId) {
    const sessionCode = $('#sessionCode').val();
    const playerName = $('#playerName').val();

    var html = `
    <div>
        <h2>Player Id: ` + playerId + `</h2>
        <h2>Player Name: ` + playerName + `</h2>
        <table id="playersTable" style="width:100%">
            <thead>
                <tr>
                    <th>X</th>
                    <th>Y</th>
                    <th>Active Actions</th>
                    <th>Send Comannd</th>
                    <th>Walk</th>
                </tr>
            </thead>
            <tbody>
                <tr>
                    <td><input id="x_` + playerId + `" type="text"/></td>
                    <td><input id="y_` + playerId + `" type="text"/></td>
                    <td><input id="actice_actions_` + playerId +`" type="text"/></td>
                    <td><input type="button" onclick="sendCommandOnClick('` + playerId + `')" value="Send Comand!"/></td>
                    <td><input type="button" onclick="initiateInfiniteWalking('` + playerId + `')" value="Walk"/></td>
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
