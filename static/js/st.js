var StoryTellers = StoryTellers || (function() {
    const ACTION_JOIN = "ACTION_JOIN";
    const ACTION_LEAVE = "ACTION_LEAVE";
    const CHAT = "CHAT";
    const STORY_ADD = "STORY_ADD";
    const STORY_CHANGE_EDITOR = "STORY_CHANGE_EDITOR";

    var ws, storyCode, authorName, chatDiv, chatInput, storyDiv, storyArea, addToStoryButton;

    function Message(messageType, storyCode, authorName, content) {
        this.messageType = messageType;
        this.storyCode = storyCode;
        this.authorName = authorName;
        this.content = content;
    }

    function initPageElements() {
        chatDiv = document.getElementById("chat_div");
        chatInput = document.getElementById("chat_input");
        storyDiv = document.getElementById("story_div");
        storyArea = document.getElementById("story_area");
        addToStoryButton = document.getElementById("add_to_story_button");
    }

    function init(code, author) {
        storyCode = code;
        authorName = author;

        initPageElements();

        var url = "ws://" + window.location.host + "/ws/" + storyCode;
        ws = new WebSocket(url);

        ws.onmessage = function (msg) {
            var message = JSON.parse(msg.data);
            console.log("received this message:", message);
            
            var meClass = message.authorName && message.authorName === authorName ? "me" : "them";

            if (message.messageType === ACTION_JOIN) {
                chatDiv.innerHTML += "<p>" + "- " + message.authorName + " has joined -</p>";
            } else if (message.messageType === CHAT) {
                chatDiv.innerHTML += "<p class=\"" + meClass + "\">" + message.authorName + ": " + message.content + "</p>";
            } else if (message.messageType === STORY_ADD) {
                storyDiv.innerHTML += "<p title=\"" + message.authorName + "\" class=\"" + meClass + "\">" + message.content + "</p>"
                storyDiv.scrollTop = storyDiv.scrollHeight;
            } else if (message.messageType === STORY_CHANGE_EDITOR) {
                var canEdit = message.authorName === authorName;
                
                addToStoryButton.disabled = !canEdit;
            }

            chatDiv.scrollTop = chatDiv.scrollHeight;
        };

        chatInput.onkeydown = function (e) {
            if (e.keyCode === 13 && chatInput.value !== "") {
                var message = new Message(CHAT, storyCode, authorName, chatInput.value);
                
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

            var message = new Message(STORY_ADD, storyCode, authorName, content);
            
            ws.send(JSON.stringify(message));

            storyArea.value = "";
        }

        ws.onopen = function() {
            var message = new Message(ACTION_JOIN, storyCode, authorName);
            ws.send(JSON.stringify(message))
        }

        chatInput.focus();
    }

    function addToStory(content) {
        var message = new Message(STORY_ADD, storyCode, authorName, content);
        ws.send(JSON.stringify(message));
    }

    return {
        init: init
    }    
})();
