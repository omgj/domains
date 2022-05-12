import SwiftUI
import Stripe
import UIKit

public struct StripeCard: UIViewRepresentable {
    @Binding var p: STPPaymentMethodParams?
    
    public init(p: Binding<STPPaymentMethodParams?>) {
        _p = p
    }

    public func makeCoordinator() -> Coordinator {
        return Coordinator(parent: self)
    }

    public func makeUIView(context: Context) -> STPPaymentCardTextField {
        let paymentCardField = STPPaymentCardTextField()
        paymentCardField.borderColor = nil
        paymentCardField.font = UIFont.systemFont(ofSize: 14.0)
        paymentCardField.textColor = UIColor.white
        paymentCardField.placeholderColor = UIColor(white: 1.0, alpha: 0.4)
        paymentCardField.postalCodeEntryEnabled = false
        paymentCardField.isOpaque = true
        if let cardParams = p?.card {
            paymentCardField.cardParams = cardParams
        }
        if let postalCode = p?.billingDetails?.address?.postalCode {
            paymentCardField.postalCode = postalCode
        }
        if let countryCode = p?.billingDetails?.address?.country {
            paymentCardField.countryCode = countryCode
        }
        paymentCardField.delegate = context.coordinator
        paymentCardField.setContentHuggingPriority(.required, for: .vertical)

        return paymentCardField
    }

    public func updateUIView(_ paymentCardField: STPPaymentCardTextField, context: Context) {
        if let cardParams = p?.card {
            paymentCardField.cardParams = cardParams
        }
        if let postalCode = p?.billingDetails?.address?.postalCode {
            paymentCardField.postalCode = postalCode
        }
        if let countryCode = p?.billingDetails?.address?.country {
            paymentCardField.countryCode = countryCode
        }
    }

    public class Coordinator: NSObject, STPPaymentCardTextFieldDelegate {
        var parent: StripeCard
        init(parent: StripeCard) {
            self.parent = parent
        }

        public func paymentCardTextFieldDidChange(_ cardField: STPPaymentCardTextField) {
            let paymentMethodParams = STPPaymentMethodParams(
                card: cardField.cardParams, billingDetails: nil, metadata: nil)
            if !cardField.isValid {
                parent.p = nil
                return
            }
            if let postalCode = cardField.postalCode, let countryCode = cardField.countryCode {
                let billingDetails = STPPaymentMethodBillingDetails()
                let address = STPPaymentMethodAddress()
                address.postalCode = postalCode
                address.country = countryCode
                billingDetails.address = address
                paymentMethodParams.billingDetails = billingDetails
            }
            parent.p = paymentMethodParams
        }
    }
}
