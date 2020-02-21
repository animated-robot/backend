const up = "up";
const down = "down";
const left = "left";
const right = "right";


let timeMs = 1;
let boundary = 100;
let stepsFromBorder = 20;
let step = 1;

let limit = () => boundary - stepsFromBorder;
let initialX = () => -boundary + 1;
let initialY = () =>  boundary - 1;


function initiateInfiniteWalking(index) {

    let sessionCode = $('#sessionCode_sendCommand_' + index).val();
    let playerId = $('#playerId_' + index).val();

    let x = $('#x_' + index).val();
    let y = $('#y_' + index).val();
    let activeActions = $('#actice_actions_' + index).val();

    let direction = right;
    let position = { x: initialX(), y: initialY()};

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
    lmt = limit()
    switch (direction) {
        case up:    return y < lmt;
        case down:  return y > -lmt;
        case left:  return x > -lmt;
        case right: return x < lmt;
        default:    return false;
    }
}