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

function registerPlayer(sessionCode, playerName) {
    this.window.socket.emit('register_player', '{"sessionCode": "' + sessionCode + '", "playerName": "' + playerName +'"}');
}

function registerPlayerOnClick() {

    let sessionCode = $('#sessionCode').val();
    let playerName = $('#playerName').val();

    registerPlayer(sessionCode, playerName)
}


const up = "up";
const down = "down";
const left = "left";
const right = "right";


const timeMs = 1;
const boundary = 100;
const stepsFromBorder = 20;
const step = 1;

const limit = boundary - stepsFromBorder;
const initialX = -boundary + 1;
const initialY =  boundary - 1;

function initiateInfiniteWalking(index) {

    let sessionCode = $('#sessionCode_sendCommand_' + index).val();
    let playerId = $('#playerId_' + index).val();

    let x = $('#x_' + index).val();
    let y = $('#y_' + index).val();
    let activeActions = $('#actice_actions_' + index).val();

    let direction = right;
    let position = { x: initialX, y: initialY};

    infiniteWalker(direction, position, playerId, sessionCode, activeActions)
}

function infiniteWalker(direction, position, playerId, sessionCode, activeActions) {

    setTimeout(function() {
        sendCommand(playerId, sessionCode, position.x / boundary, position.y / boundary, activeActions);

        data = new Date();
        console.log("inputou: " + data.getHours() + ":" + data.getMinutes() + ":" + data.getSeconds())

        walk = walking(direction, position);

        infiniteWalker(walk.direction, walk.position, playerId, sessionCode, activeActions);
    }, timeMs);

}

function walking(direction, position) {

    let newPosition;

    newPosition = move(direction, position.x, position.y);
    if (inRange(direction, newPosition.x, newPosition.y)) {
        return {
            direction: direction,
            position: newPosition
        }
    }

    return walking(changeDirection(direction), position);
}

function changeDirection(direction) {
    switch (direction) {
        case up:    return right;
        case right: return down;
        case down:  return left;
        case left:  return up;
        default:    return direction;
    }
}

function move(direction, x, y) {
    switch (direction) {
        case up:    return { x: x, y: y + step };
        case down:  return { x: x, y: y - step };
        case left:  return { x: x - step, y: y };
        case right: return { x: x + step, y: y };
        default:    return { x: x, y: y };
    }
}

function inRange(direction, x, y) {
    switch (direction) {
        case up:    return y < limit;
        case down:  return y > -limit;
        case left:  return x > -limit;
        case right: return x < limit;
        default:    return false;
    }
}




function sendCommand(playerId, sessionCode, x, y, activeActions) {
    let command = '{"playerId": "' + playerId + '", "sessionCode": "' + sessionCode +'", "activeActions": [' + activeActions + '], "direction": { "x": ' + x + ', "y": ' + y + '}}';
    this.window.socket.emit('context', command);
}

function sendCommandOnClick(index) {
    let sessionCode = $('#sessionCode_sendCommand_' + index).val();
    let playerId = $('#playerId_' + index).val();

    let x = $('#x_' + index).val();
    let y = $('#y_' + index).val();
    let activeActions = $('#actice_actions_' + index).val();

    sendCommand(playerId, sessionCode, x, y, activeActions)
}

function startSocket() {
    let socketIp = $('#socketIp').val();

    this.window.socket = io.connect(trimUrl(socketIp) + "/input");

    this.window.socket.on('player_registered', function (playerId) {
        console.log("player registered. playerId: " + playerId);
        createSendCommandBox(playerId)
    });
}

function trimUrl(url) {
    if (url === undefined)
        return url;

    url = url.trim();
    if (url.substring(url.length -1) != '/')
        return url;

    return url.substring(0, url.length -1);
}