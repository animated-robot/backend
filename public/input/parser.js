const activeActionsSeparator = ",";


function getInputContext(index) {
    const sessionCode = $('#sessionCode_sendCommand_' + index).val();
    const playerId = $('#playerId_' + index).val();

    const x = $('#x_' + index).val();
    const y = $('#y_' + index).val();
    const activeActions = $('#actice_actions_' + index).val();

    return {
        sessionCode: sessionCode,
        playerId: playerId,
        x: parseFloat(x),
        y: parseFloat(y),
        activeActions: parseActiveActions(activeActions)
    }
}

function parseActiveActions(activeActionsStr) {
    const activeActions = activeActionsStr.split(activeActionsSeparator);

    return activeActions.map(function (activeAction) {
        return activeAction.trim()
    });
}

function parseSocketIp() {
    const socketIp = $('#socketIp').val();
    return trimUrl(socketIp) + "/input";
}

function trimUrl(url) {
    if (url === undefined)
        return url;

    url = url.trim();
    if (url.substring(url.length -1) != '/')
        return url;

    return url.substring(0, url.length -1);
}
