import SwiftUI
import Stripe
import Alamofire
import PassKit

class PaymentController: UIViewController, STPAuthenticationContext, STPApplePayContextDelegate {
    func applePayContext(_ context: STPApplePayContext, didCreatePaymentMethod paymentMethod: STPPaymentMethod, paymentInformation: PKPayment, completion: @escaping STPIntentClientSecretCompletionBlock) {
        print("here 1")
        var cs = ""
        AF.request("http://localhost:3000/intent").responseDecodable(of: PaymentIntent.self) { res in
            let a = res.value?.secret ?? ""
            if a.isEmpty {
                print("empty")
                return
            }
            cs = a
        }
        print(cs)
        completion(cs, nil)
    }
    
    func applePayContext(_ context: STPApplePayContext, didCompleteWith status: STPPaymentStatus, error: Error?) {
        print("here 2")
        switch status {
                case .success:
                    print("success")
                    break
                case .error:
                    print("error")
                    break
                case .userCancellation:
            print("user cancellation")
                    break
                @unknown default:
                    fatalError()
                }
    }
    
    func authenticationPresentingViewController() -> UIViewController {
        print("SOMETHING IS GOING ON HERE")
        return UIViewController()
    }
    func payfirst(p: STPPaymentMethodParams, s: String) {
        let paymentIntentParams = STPPaymentIntentParams(clientSecret: s)
        paymentIntentParams.paymentMethodParams = p
        let aa = STPAPIClient.shared
        aa.publishableKey = ""
        aa.confirmPaymentIntent(with: paymentIntentParams) { (paymentIntent, error) in
            dump(error)
            dump(paymentIntent)
        }
    }
    
    func apple() {
        let paymentRequest = StripeAPI.paymentRequest(withMerchantIdentifier: "MERCHANT_ID", country: "MERCHANT_COUNTRY", currency: "MERCHANT_CURRENCY")
        paymentRequest.paymentSummaryItems = [
            PKPaymentSummaryItem(label: "Domains Desc", amount: 0.00),
        ]
        if let applePayContext = STPApplePayContext(paymentRequest: paymentRequest, delegate: self) {
            applePayContext.presentApplePay()
        }
    }
    func after(s: String) {
        // Billing details
        let billingDetails = STPPaymentMethodBillingDetails()
        billingDetails.name = ""
        billingDetails.email = ""
        let billingAddress = STPPaymentMethodAddress()
        billingAddress.line1 = ""
        billingAddress.postalCode = ""
        billingAddress.country = ""
        billingDetails.address = billingAddress

        // Afterpay PaymentMethod params
        let afterpayParams = STPPaymentMethodParams(afterpayClearpay: STPPaymentMethodAfterpayClearpayParams(),
                                                                                     billingDetails: billingDetails,
                                                                                     metadata: nil)

        // Shipping details
        let shippingAddress = STPPaymentIntentShippingDetailsAddressParams(line1: "")
        shippingAddress.country = ""
        shippingAddress.postalCode = ""
        let shippingDetails = STPPaymentIntentShippingDetailsParams(address: shippingAddress, name: "")
        
        let paymentIntentParams = STPPaymentIntentParams(clientSecret: s)
        paymentIntentParams.paymentMethodParams = afterpayParams
        paymentIntentParams.shipping = shippingDetails
        
        let aa = STPAPIClient.shared
        aa.publishableKey = ""
        aa.confirmPaymentIntent(with: paymentIntentParams) { (paymentIntent, error) in
            dump(error)
            dump(paymentIntent)
        }
    }
}
