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
    ChangeToolModifier : 'change-tool-modifier-event',
    PaintingChanged : 'painting-changed-event',
    PlayersChanged : 'players-changed-event',
};

// Enum:
const GameView = {
    title : 'title',
    lobby : 'lobby',
    promptselection : 'promptselection',
    artstudioEmpty : 'artstudio-empty',
    artstudioActive : 'artstudio-active',
    exhibition : 'exhibition',
    exhibitionVoting : 'exhibition-voting',
    exhibitionStickering : 'exhibition-stickering',
    showcase : 'showcase',
    gallery : 'gallery',
};

// Command:
function createCreateSessionCommand(nickName)
{
    return {
        type : CommandId.CreateSession,
        nickName : nickName, // str
    };
}

// Command:
function createJoinSessionCommand(nickName, sessionId)
{
    return {
        type : CommandId.JoinSession,
        nickName : nickName, // str
        sessionId : sessionId, // str
    };
}

// Command:
function createLeaveSessionCommand()
{
    return {
        type : CommandId.LeaveSession,
    };
}

// Command:
function createUserCommand(action)
{
    return {
        type : CommandId.User,
        action : action, // str
    };
}

// Command:
function createVoteCommand(option)
{
    return {
        type : CommandId.Vote,
        option : option, // str
    };
}

// Command:
function createPlaceStickerCommand(sticker, x, y)
{
    return {
        type : CommandId.PlaceSticker,
        sticker : sticker, // str
        x : x, // float
        y : y, // float
    };
}

// Command:
function createSetPaintingCommand(path)
{
    return {
        type : CommandId.SetPainting,
        path : path, // Any
    };
}

