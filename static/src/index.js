const axios = require("axios").default

const PRECISION = 5

var HOST = undefined
var WEBSOCK_ADDR = undefined

function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function aliveRequest()
{
}

function renewRequest()
{
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
        this.start = []
        this.ongoing = []
        this.last = []
        this.leftClicks = 0
        this.rightClicks = 0
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

        return -1
    }

    lastTouchById(id)
    {
        for (let i = 0; i < this.last.length; i++)
        {
            if (this.last[i].identifier == id)
            {
                return i
            }
        }

        return -1
    }

    copyTouch({ identifier, pageX, pageY })
    {
        return { identifier, pageX, pageY }
    }

    handleStart(e)
    {
        const touches = e.changedTouches;

        for (let i = 0; i < touches.length; i++)
        {
            this.start.push(this.copyTouch(touches[i]))
            this.ongoing.push(this.copyTouch(touches[i]))
            this.last.push(this.copyTouch(touches[i]))
        }
    }

    handleMove(e)
    {
        const touches = e.changedTouches;
        for (let i = 0; i < touches.length; i++)
        {
            const idx = this.ongoingTouchById(touches[i].identifier)
            this.last[idx] = this.copyTouch(this.ongoing[idx])
            this.ongoing[idx] = this.copyTouch(touches[i])
        }
    }

    handleEnd(e)
    {
        let fingerClicks = 0
        const touches = e.changedTouches;
        for (let i = 0; i < touches.length; i++)
        {
            const idx = this.ongoingTouchById(touches[i].identifier)
            if (idx >= 0)
            {
                let fingerDelta = {
                    x: this.last[i].pageX - this.start[i].pageX,
                    y: this.last[i].pageY - this.start[i].pageY
                }
                if (fingerDelta.x * fingerDelta.x + fingerDelta.y * fingerDelta.y < 4)
                {
                    fingerClicks++
                }

                this.start.splice(idx, 1) // remove
                this.ongoing.splice(idx, 1) // remove
                this.last.splice(idx, 1) // remove
            }
        }

        if (fingerClicks > 0)
        {
            if (fingerClicks == 1)
            {
                this.leftClicks++
            }
            else if (fingerClicks == 2)
            {
                this.rightClicks++
            }
        }
    }

    handleCancel(e)
    {
        const touches = e.changedTouches;
        for (let i = 0; i < touches.length; i++)
        {
            const idx = this.ongoingTouchById(touches[i].identifier)
            if (idx >= 0)
            {
                this.ongoing.splice(idx, 1) // remove
                this.last.splice(idx, 1) // remove
            }
        }
    }

    broadcast(sock)
    {
        if (!isSockOpen(sock))
        {
            return;
        }

        let msg = ""
        let delta_x = 0
        let delta_y = 0

        if (this.ongoing.length == 0)
        {
            return
        }

        for (let i = 0; i < this.ongoing.length; i++)
        {
            delta_x = this.ongoing[i].pageX - this.last[i].pageX
            delta_y = this.ongoing[i].pageY - this.last[i].pageY
            msg += "(" + i + "," + delta_x + "," + delta_y + ")"
        }

        if (this.leftClicks > 0)
        {
            --this.leftClicks
        }

        if (this.rightClicks > 0)
        {
            --this.rightClicks
        }

        sock.send(msg)
    }
}

function fromScreen(sx, sy)
{
    adj = {x: 0, y: 0}
    adj.x = (sx / canvas.width).toFixed(PRECISION)
    adj.y = (sy / canvas.height).toFixed(PRECISION)
    return adj
}

function init()
{
    if (getCookie("token") === null || getCookie("token") === "") {
        window.location.href = "/auth.html"
        return
    }

    var aliveIntervalID = setInterval(aliveRequest, 60 * 1000)
    var renewIntervalID = setInterval(renewRequest, 240 * 1000)
    var canvas = document.getElementById("canvas")
    canvas.width = window.innerWidth
    canvas.height = window.innerHeight
    var ctx = canvas.getContext("2d")
    var sock = new WebSocket(WEBSOCK_ADDR)
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

    canvas.addEventListener("touchstart", (e) =>
    {
        e.preventDefault()
        touchState.handleStart(e)
        touchState.broadcast(sock)
    }, false)
    
    canvas.addEventListener("touchend", (e) =>
    {
        e.preventDefault()
        touchState.handleEnd(e)
        touchState.broadcast(sock)
    }, false)
    
    canvas.addEventListener("touchcancel", (e) =>
    {
        e.preventDefault()
        touchState.handleCancel(e)
        touchState.broadcast(sock)
    }, false)
    
    canvas.addEventListener("touchmove", (e) =>
    {
        e.preventDefault()
        touchState.handleMove(e)
        touchState.broadcast(sock)
    }, false)
}

window.onload = (event) => {
    let getHostPromise = new Promise((resolve, reject) => {
        var xhr = new XMLHttpRequest()
        xhr.onreadystatechange = () => {
            if (xhr.readyState == 1) {
                xhr.send()
            }
            else if (xhr.readyState == 4) {
                resolve(xhr.response)
            }
        }
        xhr.open("GET", "/api/url")
    })

    getHostPromise.then((value) => {
        HOST = value
        WEBSOCK_ADDR = "ws://" + HOST + "/api/echo"
        init()
    })

    return false
}

window.onclose = () => {
    document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
}

