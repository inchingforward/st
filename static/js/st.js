var StoryTellers = StoryTellers || (function() {
    const ACTION_JOIN = "ACTION_JOIN";
    const ACTION_LEAVE = "ACTION_LEAVE";
    const CHAT = "CHAT";
    const STORY_ADD = "STORY_ADD";
    const STORY_CHANGE_EDITOR = "STORY_CHANGE_EDITOR";

    function Message(messageType, authorName, content) {
        this.messageType = messageType;
        this.authorName = authorName;
        this.content = content;
    }

    function init(storyCode, authorName, chat, text) {
        console.log("would initialize", storyCode, authorName);

        var url = "ws://" + window.location.host + "/ws";
        var ws = new WebSocket(url);

        ws.onmessage = function (msg) {
            var message = JSON.parse(msg.data);

            if (message.messageType === ACTION_JOIN) {
                chat.value += "** " + message.authorName + " has joined **\n";
            } else if (message.messageType === CHAT) {
                chat.value += "<" + message.authorName + "> " + message.content + "\n";
            }

            chat.scrollTop = chat.scrollHeight;
        };

        text.onkeydown = function (e) {
            if (e.keyCode === 13 && text.value !== "") {
                var message = new Message(CHAT, authorName, text.value);
                
                ws.send(JSON.stringify(message));

                text.value = "";
            }
        };

        ws.onopen = function() {
            var message = new Message(ACTION_JOIN, authorName);
            
            ws.send(JSON.stringify(message))
        }

        text.focus();
    }

    return {
        init: init
    }    
})();
