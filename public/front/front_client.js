const canvas = document.getElementById('canvas')
const ctx = canvas.getContext('2d')
const players = {}

function printSession(sessionJson) {
    let session = JSON.parse(sessionJson)
    console.log("session response: " + sessionJson)
    console.log("session code: " + session.sessionCode);
    $('#code').text(session.sessionCode)
    $('#sessionCode').text(session.sessionCode)
    $('#socketId').text(session.socketId)

    $('#playersTable tbody tr').remove()
    session.playersIds.forEach(function(element, index) {
        $('#playersTable').append("<tr><td>" + element + "</td>")
        players[element] = {
            x: canvas.width / 2 - 16,
            y: canvas.height / 2 - 16
        }
    });
}

function printInputContext(inputContextJson) {

    let inputContex = JSON.parse(inputContextJson)
    players[inputContex.playerId] = {
        ...players[inputContex.playerId],
        ...inputContex.direction
    }

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

function startSocket() {
    let socketIp = $('#socketIp').val();

    const socket = io.connect(trimUrl(socketIp) + "/front");

    socket.on('session_created', function (sessionJson) {
        printSession(sessionJson)
    });

    socket.on('session_changed', function (sessionJson) {
        printSession(sessionJson)
    });

    socket.on('input_context', function (inputContext) {
        printInputContext(inputContext)
    });

    socket.emit('create_session', 'please, create session?');
}

function trimUrl(url) {
    if (url === undefined)
        return url;

    url = url.trim();
    if (url.substring(url.length -1) != '/')
        return url;

    return url.substring(0, url.length -1);
}

const cx = canvas.width / 2 - 16
const cy = canvas.height / 2 - 16
function renderGame () {
    ctx.clearRect(0, 0, canvas.width, canvas.height)

    Object.values(players).forEach(({ x, y, color }) => {
        ctx.fillStyle = 'red'
        ctx.fillRect(cx + x * 400, cy + y * 300, 32, 32)  
    })

    data = new Date();
    console.log("renderizou: " + data.getHours() + ":" + data.getMinutes() + ":" + data.getSeconds());
    window.requestAnimationFrame(renderGame)
}

renderGame()