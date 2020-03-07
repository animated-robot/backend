function newDirection(x, y) {
    return {
        x: x,
        y: y
    }
}

function newPlayer(id) {
    return {
        id: id,
        name: "Zak",
        color: "blue",
        height: 1.90
    }
}

function newContext(sessionCode, activeActions, playerId, direction) {
    return {
        sessionCode: sessionCode,
        activeActions: activeActions,
        playerId: playerId,
        direction: direction
    }
}