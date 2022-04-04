const PRECISION = 5

class Overlay
{
}

class TransparentOverlay extends Overlay
{
}

class TouchpadOverlay extends Overlay
{
}

class ControllerOverlay extends Overlay
{
}

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
        this.last_x = 0;
        this.last_y = 0;
        this.delta_x = 0;
        this.delta_y = 0;
        this.buttons = [0, 0, 0]
    }

    update(e)
    {
        let coord = fromScreen(e.x, e.y)
        this.x = coord.x
        this.y = coord.y

        this.delta_x = (this.x - this.last_x).toFixed(PRECISION);
        this.delta_y = (this.y - this.last_y).toFixed(PRECISION);
        this.last_x = this.x;
        this.last_y = this.y;
    }

    broadcast(sock)
    {
        sock.send("m," + this.delta_x + "," + this.delta_y)
    }
}

class TouchState
{
    constructor()
    {
        this.ongoing = []
    }

    ongoingTouchById(id)
    {
        for (let i = 0; i < this.ongoing.length; i++)
        {
            if (this.ongoing[i].identifier == id)
            {
                return i;
            }
        }

        return -1;
    }

    copyTouch({ identifier, pageX, pageY })
    {
        return { identifier, pageX, pageY };
    }

    handleStart(e)
    {
        const touches = e.changedTouches;

        for (let i = 0; i < touches.length; i++)
        {
            this.ongoing.push(this.copyTouch(touches[i]))
        }
    }

    handleMove(e)
    {
        const touches = e.changedTouches;
        for (let i = 0; i < touches.length; i++)
        {
            const idx = this.ongoingTouchById(touches[i].identifier);
        }
    }

    handleEnd(e)
    {
        const touches = e.changedTouches;
        for (let i = 0; i < touches.length; i++)
        {
            const idx = this.ongoingTouchById(touches[i].identifier);
            if (idx >= 0)
            {
                this.ongoing.splice(idx, 1); // remove
            }
        }
    }

    handleCancel(e)
    {
        const touches = e.changedTouches;
        for (let i = 0; i < touches.length; i++)
        {
            const idx = this.ongoingTouchById(touches[i].identifier);
            if (idx >= 0)
            {
                this.ongoing.splice(idx, 1); // remove
            }
        }
    }

    broadcast(sock)
    {
        let msg = "";
        for (let i = 0; i < this.ongoing.length; i++)
        {
            msg += "(" + i + "," + "x" + "," + "y" + ")"
        }
        sock.send(msg);
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
    var touchState = new TouchState()

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

    canvas.addEventListener("touchstart", (e) =>
    {
        e.preventDefault()
        touchState.handleStart(e)
        touchState.broadcast(sock);
    }, false)
    
    canvas.addEventListener("touchend", (e) =>
    {
        e.preventDefault()
        touchState.handleEnd(e)
        touchState.broadcast(sock);
    }, false)
    
    canvas.addEventListener("touchcancel", (e) =>
    {
        e.preventDefault()
        touchState.handleCancel(e)
        touchState.broadcast(sock);
    }, false)
    
    canvas.addEventListener("touchmove", (e) =>
    {
        e.preventDefault()
        touchState.handleMove(e)
        touchState.broadcast(sock);
    }, false)

    return false
})

