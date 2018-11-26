


// fuzzy find hubs

ws = new WebSocket("ws://localhost:1212/find-hubs/123");

ws.send("hello");

ws.onmessage = function (evt) {
    var messages = evt.data.split('\n');
    var msg = JSON.parse(messages[0]);
    console.log(msg);
};

// Hub Message
mainWS.send(JSON.stringify({
    type: 'hubMessage',
    body: {
        HubID: 'hubid1',
        Sender: {
            ID: 'userid1',
            Username: 'username1'
        },
        Message: 'whats up?',
    },
}));

// User Message
mainWS.send(JSON.stringify({
    type: 'userMessage',
    body: {
        Sender: {
            ID: 'senderid',
            Username: 'senderusername',
        },
        Recipient: {
            ID: 'recipientid',
            Username: 'recipientusername',
        },
        Message: 'whats up?',
    },
}));

// Friend Request
mainWS.send(JSON.stringify({
    type: 'friendRequest',
    body: {
        From: {
            ID: 'senderid',
            Username: 'senderusername',
        },
        To: {
            ID: 'recipientid',
            Username: 'recipientusername',
        },
    },
}));

// Join Invitation
mainWS.send(JSON.stringify({
    type: 'joinInvitation',
    body: {
        HubID: 'thenewhub',
        From: {
            ID: 'senderid',
            Username: 'senderusername',
        },
        To: {
            ID: 'recipientid',
            Username: 'recipientusername',
        },
    },
}));


// notifications
nws = new WebSocket("ws://localhost:1212/ws/notifications?token=");

nws.onmessage = function (evt) {
    var messages = evt.data.split('\n');
    var msg = JSON.parse(messages[0]);
    console.log(msg);
};


nws.send(""); // recipient user token