import {ElMessageBox} from 'element-plus'

let show = false

/**
 * Success
 * @param title
 * @param msg
 * @constructor
 */
function Success(title, msg) {
    MsgPop(title, msg, "iconfont icon-check-circle", "#67C23A")
}


/**
 * Error
 * @param title
 * @param msg
 * @constructor
 */
function Error(title, msg) {
    MsgPop(title, msg, "iconfont icon-times-circle", "#F56C6C")
}


/**
 * Warning
 * @param title
 * @param msg
 * @constructor
 */
function Warning(title, msg) {
    MsgPop(title, msg, "iconfont icon-exclamation-circle", "#E6A23C")
}

/**
 * Alert
 * @param title
 * @param msg
 * @param icon
 * @param iconColor
 * @constructor
 */
function MsgPop(title, msg, icon, iconColor) {
    if (show) {
        return
    }
    show = true
    ElMessageBox.alert(
        '<div>' +
        '<div style="font-size: 14px;line-height: 20px;margin-bottom: 15px;margin-top: 5px;color: #303133;font-weight: 600">' +
        '<i class="' + icon + '" style="color: ' + iconColor + '"></i> <span style="font-weight: 500">' + msg + '</span> </div>' +
        '</div>',
        title,
        {
            dangerouslyUseHTMLString: true,
            draggable: true,
            closeOnClickModal: false,
            showClose: false,
            confirmButtonText: '确认',
            callback: () => {
                show = false
            },
        }
    )
}

/**
 *
 * @param title
 * @param msg
 * @param cancel
 * @param confirm
 * @constructor
 */
function Confirm(title, msg, cancel, confirm) {
    ElMessageBox.confirm(
        '<div>' +
        '<div style="font-size: 14px;line-height: 20px;margin-bottom: 15px;margin-top: 5px;color: #303133;font-weight: 600">' +
        '<i class="iconfont icon-exclamation-circle" style="color: #E6A23C"></i> <span style="font-weight: 500">' + msg + '</span> </div>' +
        '</div>',
        title,
        {
            dangerouslyUseHTMLString: true,
            draggable: true,
            closeOnClickModal: false,
            showClose: false,
            confirmButtonText: '确认',
            cancelButtonText: "取消",
        }
    ).then(() => {
    }).catch((action) => {
        if (action === 'cancel') {
            cancel()
        } else {
            confirm()
        }
    })
}

/**
 * HandleError
 * @param error
 * @constructor
 */
function HandleError(error) {
    if (show) {
        return
    }
    show = true
    let response = error.response
    let status = ""
    let err = ""
    let msg = ""
    if (response !== undefined && response.data !== undefined) {
        err = response.data.error
        msg = response.data.msg
        status = response.status
    }
    ElMessageBox.alert(
        '<div>' +
        '<div style="font-size: 14px;line-height: 20px;margin-bottom: 15px;margin-top: 5px;color: #303133;font-weight: 600">' +
        '<i class="iconfont icon-times-circle" style="color: #F56C6C"></i> ' + msg + ' </div>' +
        '<div style="font-size: 14px;color: #606060">错误代码：<span style="color: black">' + status + '</span></div>' +
        // '<div style="font-size: 14px;color: #606060">错误信息：<span style="color: black">' + msg + '</span></div>' +
        '<div style="font-size: 14px;color: #606060">输出信息：</div>' +
        '<div style="font-size: 12px;color: #f6f6f6;margin-bottom: 10px;margin-top:3px;background-color: #606060;padding: 4px 10px;border-radius: 4px">' +
        '<span>' + err + '</span>' +
        '</div>',
        "发生错误",
        {
            dangerouslyUseHTMLString: true,
            draggable: true,
            closeOnClickModal: false,
            showClose: false,
            confirmButtonText: '确认',
            callback: () => {
                show = false
            },
        }
    )
}

/**
 *
 * @param milli
 * @param fmt
 * @returns {*}
 * @constructor
 */
function UnixMilliToDate(milli, fmt) {
    if (milli <= 0) {
        return ""
    }
    return FormatDate(fmt === "" ? "YYYY/mm/dd HH:MM:SS" : fmt, new Date(milli))
}

/**
 *
 * @param sec
 * @param fmt
 * @returns {*}
 * @constructor
 */
function SecondToDate(sec, fmt) {
    if (sec <= 0) {
        return ""
    }
    return FormatDate(fmt === "" ? "YYYY/mm/dd HH:MM:SS" : fmt, new Date(sec * 1000))
}

/**
 *
 * @param fmt
 * @param date
 * @returns {*}
 */
function FormatDate(fmt, date) {
    let ret;
    const opt = {
        "Y+": date.getFullYear().toString(),        // 年
        "m+": (date.getMonth() + 1).toString(),     // 月
        "d+": date.getDate().toString(),            // 日
        "H+": date.getHours().toString(),           // 时
        "M+": date.getMinutes().toString(),         // 分
        "S+": date.getSeconds().toString()          // 秒
        // 有其他格式化字符需求可以继续添加，必须转化成字符串
    };
    for (let k in opt) {
        ret = new RegExp("(" + k + ")").exec(fmt);
        if (ret) {
            fmt = fmt.replace(ret[1], (ret[1].length === 1) ? (opt[k]) : (opt[k].padStart(ret[1].length, "0")))
        }

    }
    return fmt;
}

/**
 *
 * @param bytes
 * @returns {string}
 * @constructor
 */
function FormatBytesSize(bytes) {
    let b = 1024
    let k = b * 1024
    let m = k * 1024
    let g = m * 1024
    if (bytes <= 0) {
        return "0 B"
    } else if (bytes > 0 && bytes < b) {
        return bytes + " B"
    } else if (bytes >= b && bytes < k) {
        return (bytes / b).toFixed(1) + " K"
    } else if (bytes >= k && bytes < m) {
        return (bytes / k).toFixed(1) + " M"
    } else if (bytes >= m && bytes < g) {
        return (bytes / m).toFixed(2) + " G (" + (bytes / k).toFixed(1) + " M" + ")"
    } else {
        return (bytes / m).toFixed(2) + " G"
    }
}

/**
 *
 * @param bytes
 * @returns {string}
 * @constructor
 */
function FormatBytesSizeM(bytes) {
    let k = 1024 * 1024
    if (bytes <= 0) {
        return "0 M"
    } else {
        return (bytes / k).toFixed(1) + " M"
    }
}

/**
 *
 * @param bytes
 * @returns {string}
 * @constructor
 */
function FormatBytesSizeG(bytes) {
    let b = 1024
    let k = b * 1024
    let m = k * 1024
    if (bytes <= 0) {
        return "0 G"
    } else if (bytes > 0 && bytes < m) {
        return (bytes / k).toFixed(1) + " M"
    } else {
        return (bytes / m).toFixed(2) + " G"
    }
}

/**
 *
 * @param bytes
 * @returns {string}
 * @constructor
 */
function FormatBytesSpeed(bytes) {
    return FormatBytesSize(bytes) + "/s"
}

/**
 *
 * @param packet
 * @returns {string}
 * @constructor
 */
function FormatPacketSpeed(packet) {
    return FormatPacketSize(packet) + "pps"
}

/**
 *
 * @param packet
 * @returns {string}
 * @constructor
 */
function FormatPacketSize(packet) {
    let step = 1000
    let b = step
    let k = b * step
    let m = k * step
    let g = m * step
    if (packet <= 0) {
        return "0 "
    } else if (packet > 0 && packet < b) {
        return packet + " "
    } else if (packet >= b && packet < k) {
        return (packet / b).toFixed(1) + " K"
    } else if (packet >= k && packet < m) {
        return (packet / k).toFixed(1) + " M"
    } else if (packet >= m && packet < g) {
        return (packet / m).toFixed(2) + " G"
    } else {
        return (packet / m).toFixed(2) + " G"
    }
}

export default {
    HandleError,
    Warning,
    MsgPop,
    Error,
    Success,
    Confirm,
    UnixMilliToDate,
    SecondToDate,
    FormatBytesSize,
    FormatBytesSizeM,
    FormatBytesSizeG,
    FormatBytesSpeed,
    FormatPacketSpeed,
    FormatPacketSize,
    FormatDate
}