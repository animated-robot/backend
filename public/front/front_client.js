const socket = io.connect('http://localhost:8080/front');

socket.on('created_session', function (sessionCode) {
    console.log(sessionCode);
});

socket.emit('create_session', 'joaozinho');