const up = "up";
const down = "down";
const left = "left";
const right = "right";


let timeMs = 10;
let boundary = 100;
let stepsFromBorder = 20;
let step = 1;

let limit = () => boundary - stepsFromBorder;
let initialX = () => -boundary + 1;
let initialY = () =>  boundary - 1;


function initiateInfiniteWalking(index) {

    const sessionCode = $('#sessionCode_sendCommand_' + index).val();
    const playerId = $('#playerId_' + index).val();

    const activeActions = [];
    const direction = right;
    const position = { x: initialX(), y: initialY()};

    infiniteWalker(direction, position, playerId, sessionCode, activeActions)
}

function infiniteWalker(direction, position, playerId, sessionCode, activeActions) {

    setTimeout(function() {
        sendCommand(playerId, sessionCode, position.x / boundary, position.y / boundary, activeActions);

        const data = new Date();
        console.log("inputou: " + data.getHours() + ":" + data.getMinutes() + ":" + data.getSeconds());

        const walk = walking(direction, position);

        infiniteWalker(walk.direction, walk.position, playerId, sessionCode, activeActions);
    }, timeMs);

}

function walking(direction, position) {

    const newPosition = move(direction, position.x, position.y);
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
    const lmt = limit();
    switch (direction) {
        case up:    return y < lmt;
        case down:  return y > -lmt;
        case left:  return x > -lmt;
        case right: return x < lmt;
        default:    return false;
    }
}
