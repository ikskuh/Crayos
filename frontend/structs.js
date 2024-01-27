const CommandId = {
    CreateSession : 'create-session-command',
    JoinSession : 'join-session-command',
    LeaveSession : 'leave-session-command',
    User : 'user-command',
    Vote : 'vote-command',
    PlaceSticker : 'place-sticker-command',
    SetPainting : 'set-painting-command',
};

const EventId = {
    EnterSession : 'enter-session-event',
    JoinSessionFailed : 'join-session-failed-event',
    Kicked : 'kicked-event',
    ChangeGameView : 'change-game-view-event',
    TimerChanged : 'timer-changed-event',
    ChangeToolModifier : 'change-tool-modifier-event',
    PaintingChanged : 'painting-changed-event',
    PlayersChanged : 'players-changed-event',
    PlayerReadyChanged : 'player-ready-changed-event',
};

// Enum:
const GameView = {
    title : 'title',
    lobby : 'lobby',
    promptselection : 'promptselection',
    artstudioGeneric : 'artstudio-generic',
    artstudioActive : 'artstudio-active',
    artstudioSticker : 'artstudio-sticker',
    gallery : 'gallery',
    podium : 'podium',
};

// Enum:
const Effect = {
    flashlight : 'flashlight',
    drunk : 'drunk',
    flip : 'flip',
    swap_tool : 'swap_tool',
    lock_pencil : 'lock_pencil',
};

// Enum:
const UserAction = {
    setReady : 'set-ready',
    setNotReady : 'set-not-ready',
    continueGame : 'continue',
};

// Enum:
const Backdrop = {
    arctic : 'arctic',
    graveyard : 'graveyard',
    pirateShip : 'pirate_ship',
    theaterStage1 : 'theater_stage1',
};

// Command:
function sendCreateSessionCommand(nickName)
{
    socket.send(JSON.stringify({
        type : CommandId.CreateSession,
        nickName : nickName, // str
    }));
}

// Command:
function sendJoinSessionCommand(nickName, sessionId)
{
    socket.send(JSON.stringify({
        type : CommandId.JoinSession,
        nickName : nickName, // str
        sessionId : sessionId, // str
    }));
}

// Command:
function sendLeaveSessionCommand()
{
    socket.send(JSON.stringify({
        type : CommandId.LeaveSession,
    }));
}

// Command:
function sendUserCommand(action)
{
    socket.send(JSON.stringify({
        type : CommandId.User,
        action : action, // UserAction
    }));
}

// Command:
function sendVoteCommand(option)
{
    socket.send(JSON.stringify({
        type : CommandId.Vote,
        option : option, // str
    }));
}

// Command:
function sendPlaceStickerCommand(sticker, x, y)
{
    socket.send(JSON.stringify({
        type : CommandId.PlaceSticker,
        sticker : sticker, // str
        x : x, // float
        y : y, // float
    }));
}

// Command:
function sendSetPaintingCommand(path)
{
    socket.send(JSON.stringify({
        type : CommandId.SetPainting,
        path : path, // Any
    }));
}

