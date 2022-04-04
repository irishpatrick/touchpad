const PRECISION = 5

class Sprite
{
    constructor()
    {
        this.image = null;
        this.ready = false;
        this.x = 0;
        this.y = 0;
        this.w = 0;
        this.h = 0;
    }

    draw(ctx)
    {
        if (!this.ready)
        {
            return;
        }
    }
}

class MouseState
{
    constructor()
    {
        this.x = 0
        this.y = 0
        this.buttons = [0, 0, 0]
    }

    update(e)
    {
        let coord = fromScreen(e.x, e.y)
        this.x = coord.x
        this.y = coord.y
    }

    broadcast(sock)
    {
        sock.send("m," + this.x + "," + this.y)
    }
}

function fromScreen(sx, sy)
{
    adj = {x: 0, y: 0}
    adj.x = (sx / canvas.width).toFixed(PRECISION)
    adj.y = (sy / canvas.height).toFixed(PRECISION)
    return adj
}

window.addEventListener("load", (e) => 
{
    var canvas = document.getElementById("canvas")
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
    var ctx = canvas.getContext("2d")
    var sock = new WebSocket(ADDR)
    var mouseState = new MouseState()

    sock.onopen = function(e)
    {
    }

    sock.onclose = function(e)
    {
    }

    sock.onmessage = function(e)
    {
    }

    sock.onerror = function(e)
    {
    }

    canvas.addEventListener("mousedown", (e) =>
    {
        
    }, false)

    canvas.addEventListener("mouseup", (e) =>
    {
        
    }, false)

    canvas.addEventListener("mousemove", (e) =>
    {
        mouseState.update(e)
        mouseState.broadcast(sock)
    }, false)

    return false
})

