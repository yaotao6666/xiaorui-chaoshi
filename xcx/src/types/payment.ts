export interface XcxPayParams {
  timeStamp: string
  nonceStr: string
  package: string
  signType: string
  paySign: string
}

export interface PrepareOrderPaymentPayload {
  order_id: number
  order_no: string
  merchant_id: number
  pay_amount: number
  channel: string
  pay_params: XcxPayParams
  return_target: string
}
