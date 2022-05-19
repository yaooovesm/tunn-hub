/* eslint-disable */

import publicStorage from "@/public.storage";

//TODO 改写为组件

window.addEventListener("beforeunload", () => {
    clear()
});


let active = {
    id: "",
    instance: null,
}

function clear() {
    if (active.instance !== null) {
        try {
            active.instance.Close()
        } catch (_) {

        }
    }
    active.instance = null
    active.id = ""
}

/**
 * reporter client
 * @param res
 * @param recv
 * @param close
 * @param error
 * @param interval
 * @constructor
 */
export default function ReporterClient(res, recv, close, error, interval) {
    clear()
    let recvFunc = recv
    let closeFunc = close
    let errorFunc = error
    let id = crypto.randomUUID();
    let started = false
    let ws = null
    this.Start = function (who) {
        if ("WebSocket" in window) {
            publicStorage.Load()
            ws = new WebSocket("ws://" + window.location.hostname + ":" + publicStorage.User.reporter + "/reporter")
            started = true
            ws.onopen = function () {
                console.log(who + " open " + id)
                try {
                    ws.send(JSON.stringify(
                        {
                            token: publicStorage.User.token,
                            resources: res,
                            interval: interval <= 0 ? 5000 : interval
                        }
                    ))
                } catch (e) {
                    if (errorFunc !== null) {
                        errorFunc(e)
                    }
                }
            }
            ws.onmessage = function (e) {
                try {
                    if (e.data instanceof Blob) {
                        let blob = e.data;
                        //通过FileReader读取数据
                        let reader = new FileReader();
                        reader.readAsBinaryString(blob);
                        reader.onload = function () {
                            if (recvFunc !== null) {
                                recvFunc(reader.result)
                            }
                        }
                    }
                } catch (err) {
                    if (errorFunc !== null) {
                        errorFunc(err)
                    }
                }
            }
            let that = this
            ws.onclose = function () {
                console.log("closed " + id)
                if (closeFunc !== null) {
                    closeFunc()
                }
                console.log(active)
            }
            active.id = id
            active.instance = this
            console.log(active)
        } else {
            if (errorFunc !== null) {
                errorFunc("浏览器不支持")
            }
        }
    }
    this.Close = function (who) {
        console.log(who + " close " + id)
        if (!started) {
            return
        }
        //send close
        try {
            ws.send("close")
            ws.close()
            if (closeFunc !== null) {
                closeFunc()
            }
        } catch (_) {
        }
        active.instance = null
        active.id = ""
        started = false
    }
}