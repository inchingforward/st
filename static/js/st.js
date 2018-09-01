var StoryTellers = StoryTellers || (function() {
    const ACTION_JOIN = "ACTION_JOIN";
    const ACTION_LEAVE = "ACTION_LEAVE";
    const CHAT = "CHAT";
    const STORY_ADD = "STORY_ADD";
    const STORY_CHANGE_EDITOR = "STORY_CHANGE_EDITOR";

    var ws, storyCode, authorName, chatArea, chatInput, storyDiv, storyArea, addToStoryButton;

    function Message(messageType, authorName, content) {
        this.messageType = messageType;
        this.authorName = authorName;
        this.content = content;
    }

    function initPageElements() {
        chatArea = document.getElementById("chat_area");
        chatInput = document.getElementById("chat_input");
        storyDiv = document.getElementById("story_div");
        storyArea = document.getElementById("story_area");
        addToStoryButton = document.getElementById("add_to_story_button");
    }

    function init(code, author) {
        storyCode = code;
        authorName = author;

        initPageElements();

        var url = "ws://" + window.location.host + "/ws";
        ws = new WebSocket(url);

        ws.onmessage = function (msg) {
            var message = JSON.parse(msg.data);

            if (message.messageType === ACTION_JOIN) {
                chatArea.value += "** " + message.authorName + " has joined **\n";
            } else if (message.messageType === CHAT) {
                chatArea.value += "<" + message.authorName + "> " + message.content + "\n";
            } else if (message.messageType === STORY_ADD) {
                storyDiv.innerHTML += "<p title=\"" + message.authorName + "\">" + message.content + "</p>"
                storyDiv.scrollTop = storyDiv.scrollHeight;
            } else if (message.messageType === STORY_CHANGE_EDITOR) {
                var canEdit = message.authorName === authorName;
                
                addToStoryButton.disabled = !canEdit;
            }

            chatArea.scrollTop = chatArea.scrollHeight;
        };

        chatInput.onkeydown = function (e) {
            if (e.keyCode === 13 && chatInput.value !== "") {
                var message = new Message(CHAT, authorName, chatInput.value);
                
                ws.send(JSON.stringify(message));

                chatInput.value = "";
            }
        };

        addToStoryButton.onclick = function() {
            var content = storyArea.value;
            if (!content) {
                alert("You haven't typed anything!");
                storyArea.focus();
                return;
            }

            var message = new Message(STORY_ADD, authorName, content);
            
            ws.send(JSON.stringify(message));

            storyArea.value = "";
        }

        ws.onopen = function() {
            var message = new Message(ACTION_JOIN, authorName);
            ws.send(JSON.stringify(message))
        }

        chatInput.focus();
    }

    function addToStory(text) {
        var message = new Message(STORY_ADD, authorName, text);
        ws.send(JSON.stringify(message));
    }

    return {
        init: init
    }    
})();
