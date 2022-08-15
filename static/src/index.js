const PRECISION = 5
/** endpoint constants found in index.html **/

function showLogin(visible)
{
    let elem = document.getElementById("login")
    if (visible)
    {
        element.classList.remove("hidden")
    }
    else
    {
        element.classList.append("hidden")
    }
}

function aliveRequest()
{
    var xhr = new XMLHttpRequest()
    xhr.open("POST", ALIVE_ENDPOINT, false)
    xhr.send()
}

function renewRequest()
{
    var xhr = new XMLHttpRequest()
    xhr.open("POST", RENEW_ENDPOINT, false)
    xhr.send()
}

function bindRequest(formElem)
{

    var xhr = new XMLHttpRequest()
    xhr.open("POST", BIND_ENDPOINT, false)
    var formData = new FormData(formElem)
    formData.append("test", "hello world")
    //xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded')
    xhr.send(formData)
}

function isSockOpen(sock)
{
    return sock.readyState === WebSocket.OPEN
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
        if (!isSockOpen(sock))
        {
            return;
        }

        sock.send("m," + this.delta_x + "," + this.delta_y)
    }
}

class TouchState
{
    constructor()
    {
        this.ongoing = []
        this.last = []
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

    lastTouchById(id)
    {
        for (let i = 0; i < this.last.length; i++)
        {
            if (this.last[i].identifier == id)
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
            this.last.push(this.copyTouch(touches[i]))
        }
    }

    handleMove(e)
    {
        const touches = e.changedTouches;
        for (let i = 0; i < touches.length; i++)
        {
            const idx = this.ongoingTouchById(touches[i].identifier);
            this.last[idx] = this.copyTouch(this.ongoing[idx]);
            this.ongoing[idx] = this.copyTouch(touches[i]);
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
                this.last.splice(idx, 1); // remove
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
                this.last.splice(idx, 1); // remove
            }
        }
    }

    broadcast(sock)
    {
        if (!isSockOpen(sock))
        {
            return;
        }

        let msg = "";
        let delta_x = 0;
        let delta_y = 0;
        for (let i = 0; i < this.ongoing.length; i++)
        {
            delta_x = this.ongoing[i].pageX - this.last[i].pageX;
            delta_y = this.ongoing[i].pageY - this.last[i].pageY;
            msg += "(" + i + "," + delta_x + "," + delta_y + ")"
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
    document.getElementById("login").addEventListener("submit", (e) => {
        e.preventDefault() // don't submit
        console.log("bind")
        bindRequest(document.getElementById("login"))
    })
    var aliveIntervalID = setInterval(aliveRequest, 60 * 1000)
    var renewIntervalID = setInterval(renewRequest, 240 * 1000)
    var canvas = document.getElementById("canvas")
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
    var ctx = canvas.getContext("2d")
    var sock = new WebSocket(ADDR)
    //var mouseState = new MouseState()
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
        //mouseState.update(e)
        //mouseState.broadcast(sock)
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

