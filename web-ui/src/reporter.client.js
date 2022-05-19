/* eslint-disable */

import publicStorage from "@/public.storage";

class ReporterClient {
    constructor(res, recv, close, interval) {
        this.Start = function () {
            if ("WebSocket" in window) {
                publicStorage.Load()
                let ws = new WebSocket("ws://" + window.location.hostname + ":" + publicStorage.User.reporter + "/reporter")
                this.Close = function () {
                    //send close
                    ws.send("close")
                    ws.close()
                    if (close !== null) {
                        close()
                    }
                }
                ws.onopen = function () {
                    ws.send(JSON.stringify(
                        {
                            token: publicStorage.User.token,
                            resources: res,
                            interval: interval <= 0 ? 5000 : interval
                        }
                    ))
                }
                ws.onmessage = function (e) {
                    if (e.data instanceof Blob) {
                        let blob = e.data;
                        //通过FileReader读取数据
                        let reader = new FileReader();
                        //以下这两种方式我都可以解析出来，因为Blob对象的数据可以按文本或二进制的格式进行读取
                        reader.readAsBinaryString(blob);
                        //reader.readAsText(blob, 'utf8');
                        let that = this
                        reader.onload = function () {
                            if (recv !== null) {
                                recv(reader.result)
                            }
                            //console.log(reader.result);//这个就是解析出来的数据
                        }
                    }
                }
                ws.onclose = function () {
                }
            }
        }
    }

    Start() {
    }

    Close() {
    }
}

export default ReporterClient