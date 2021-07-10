export function parseOrder(order){
  let paytype, status, parsedOrder
  switch (order.paytype) {
    case "month":
      paytype = "月付"
      break
    case "year":
      paytype = "年付"
      break
  }
  switch (order.status) {
    case "waiting":
      status = "待运营师审核"
      break
    case "pass":
      status = "待工程师部署"
      break
    case "on":
      status = "运行中"
      break
    case "pause":
      status = "暂停使用"
      break
    case "canceled":
      status = "已取消"
      break
    case "suspend":
      status = "已欠费"
      break
    case "retrieve":
      status = "待回收设备"
  }
  parsedOrder = {
    device: order.device,
    id: order.id,
    comment: order.comment,
    paytype: paytype,
    status: status,
    createTime: order.create_at,
  }
  if(order.message !== '') {
    parsedOrder.message = order.message
  }

  return parsedOrder

}

export function getType(status) {
  switch (status) {
    case "待运营师审核":
      return "default"
    case "待工程师部署":
      return "warning"
    case "运行中":
      return "success"
    case "暂停使用":
      return "info"
    case "已取消":
      return "danger"
    case "已欠费":
      return "danger"
    case "待回收设备":
      return "danger"

  }
}

export function getEffect(status) {
  // use switch to modify ui easily
  switch (status) {
    case "待运营师审核":
      return "light"
    case "待工程师部署":
      return "light"
    case "运行中":
      return "light"
    case "暂停使用":
      return "light"
    case "已取消":
      return "light"
    case "已欠费":
      return "light"
    case "待回收设备":
      return "dark"
  }
}
