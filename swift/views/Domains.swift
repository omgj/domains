import SwiftUI
import Stripe
import Alamofire

struct PaymentIntent: Decodable {
    var secret: String
}

struct Rjs: Decodable, Hashable {
    var domain: String = ""
    var tld: String = ""
    var price: Int64 = 0
}

struct Domains: View {
    @Binding var showdomains: Bool
    @State var p: STPPaymentMethodParams?
    @State var cart: [Rjs] = []
    @State var done = false
    @State var purchasing = false
    @State var cartopen = false
    @State var av = false
    @State var dplace = "Search for your domain ..."
    @State var count = 1
    let mind: CGFloat = 50
    @FocusState private var isFocused: Bool
    let timer = Timer.publish(every: 0.5, on: .main, in: .common).autoconnect()
    let qtimer = Timer.publish(every: 0.1, on: .main, in: .common).autoconnect()
    let e = PaymentController()
    @State var rjview: [Rjs] = []
    @State var name = ""
    @State var qnow = false
    @State var dcache: [String: [Rjs]] = [String: [Rjs]]()
    @State var ph = false
    @State var pl = false
    @State var az = false
    @State var za = false
    @State var abcnone = true
    @State var nameempty = true
    @Binding var tlds: [Rjs]
    @State var rjdview: [Rjs] = []
    @State var pop = false
    func payfirst() {
        if p == nil {
            return
        }
        AF.request("http://localhost:3000/intent").responseDecodable(of: PaymentIntent.self) { res in
            let a = res.value?.secret ?? ""
            if a.isEmpty {
                print("empty")
                return
            }
            print(a)
            e.payfirst(p: p!, s: a)
        }
    }
    func after() {
//        e.after(s: a)
    }
    
    func check() {
        if p != nil {
            withAnimation{
                done = true
            }
            return
        }
        withAnimation{
            done = false
        }
    }
    
    func back() {
        if cartopen { withAnimation{ cartopen = false; isFocused = true }; return }
        timer.upstream.connect().cancel()
        qtimer.upstream.connect().cancel()
        withAnimation{ showdomains = false }
    }
    
    func avs() {
        withAnimation {
            av.toggle()
            pop = true
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) {
                withAnimation { pop = false }
            }
        }
        if !av {
            var jk: [Rjs] = []
            for i in tlds {
                var yes = ""
                for ii in rjview {
                    if ii.tld == i.tld {
                        yes = ii.domain
                    }
                }
                if yes != "" {
                    jk.append(Rjs(domain: yes, tld: i.tld))
                }
            }
            withAnimation { rjdview = jk }
            return
        }
        withAnimation { rjdview = [] }
    }
    
    func tfield() {
        if count == dplace.count {
            count = 0
        }
        count += 1
    }
    
    func typing(word: String) {
        rjview = []
        withAnimation { qnow = true }
        let a = word.components(separatedBy: ".")
        if a.count == 1 {
            if a[0] == "" {
                withAnimation { nameempty = true }
                return
            }
            withAnimation { nameempty = false }
            let b = dcache[word] ?? []
            if b.isEmpty {
                print("New word: \(word)");
                AF.request("http://localhost:3000/cache?q=\(word)").responseDecodable(of: [Rjs].self) { res in
                    var a = res.value ?? []
                    if let error = res.error {
                        print("error with word: \(word)")
                        debugPrint(error)
                        return
                    }
                    if a.isEmpty { print("empty"); return }
                    if pl { dcache[word] = a; a.sort { $0.price < $1.price }; rjview = a; withAnimation { av = true }; return }
                    if az { dcache[word] = a; a.sort { $0.tld < $1.tld }; rjview = a; withAnimation { av = true }; return }
                    if za { dcache[word] = a; a.sort { $0.tld > $1.tld }; rjview = a; withAnimation { av = true }; return }
                    rjview = a; withAnimation { av = true }; dcache[word] = a
                }
                return
            }
            print("Old word \(word)")
            var q = dcache[word] ?? []
            if q.isEmpty { print("empty"); return }
            if pl { q.sort { $0.price < $1.price }; rjview = q; withAnimation { av = true }; return }
            if az { q.sort { $0.tld < $1.tld }; rjview = q; withAnimation { av = true }; return }
            if za { q.sort { $0.tld > $1.tld }; rjview = q; withAnimation { av = true }; return }
            q.sort { $0.price > $1.price }; rjview = q; withAnimation { av = true }
            return
        }
    }
    
    func psort() {
        if ph {
            withAnimation { ph = false; pl = true; az = false; za = false; abcnone = true; rjview.sort { $0.price < $1.price }}
            return
        }
        if pl {
            withAnimation{ ph = false; pl = false; az = false; za = false; abcnone = true }
            var t: [Rjs] = []
            var b: [Rjs] = []
            for ii in rjview {
                if cart.contains(ii) {
                    t.append(ii)
                    continue
                }
                b.append(ii)
            }
            withAnimation { rjview = t + b }
            return
        }
        withAnimation {
            ph = true; pl = false; az = false; za = false; abcnone = true;
            rjview.sort { $0.price > $1.price }
        }
    }
    
    func azsort() {
        if az {
            withAnimation { az = false; za = true; abcnone = false; ph = false; pl = false; rjview.sort { $0.tld > $1.tld }}
            return
        }
        if za {
            withAnimation{ za = false; az = false; abcnone = true; ph = false; pl = false }
            var t: [Rjs] = []
            var b: [Rjs] = []
            for i in rjview {
                if cart.contains(i) {
                    t.append(i)
                    continue
                }
                b.append(i)
            }
            withAnimation { rjview = t + b }
            return
        }
        withAnimation{
            az = true; abcnone = false; ph = false; pl = false;
            rjview.sort { $0.tld < $1.tld }
        }
    }
    
    func addtocart(rj: Rjs) {
        var removed = false
        for (i, a) in cart.enumerated() {
            if a.domain == rj.domain && rj.tld == a.tld && a.price == rj.price {
                cart.remove(at:i)
                removed = true
            }
        }
        if !removed {
            cart.append(rj)
        }
    }

    func transfer() {
        
    }

    
    var body: some View {
        // ZStack Wrapper
        ZStack {
            // Background
            Image("star").resizable().ignoresSafeArea().onReceive(timer) { _ in check() }
            //Background
            
            // Domains Wrapper
            VStack {
                
                // Domains Header
                HStack(spacing: 0) {
                    Button(action: { back() }, label: {
                        HStack {
                            Text("     ")
                            Image(systemName: "arrow.uturn.left").font(.system(size:12.0)).foregroundColor(.red.opacity(0.8)).padding(.top, cartopen ? 7 : 0)
                        }
                    })
                    Spacer()
                    VStack {
                        HStack(spacing: 0) {
                            Image(systemName: cartopen ? "cart" : "globe").font(.system(size: 12.0)).foregroundColor(.blue).scaleEffect(pop ? 1.2 : 1.0)
                            Text(cartopen ? "Cart" : " Domains").font(.system(size:18.0)).fontWeight(.bold).foregroundColor(.white.opacity(pop ? 0.8 : 1.0)).onTapGesture{ if !cartopen { avs() }}
                        }
                        if nameempty {
                            HStack {
                                Text("Search & Transfer").font(.system(size:14.0)).fontWeight(.light).foregroundColor(.gray)
                            }
                        } else {
                            Button(action: { avs() }, label: {
                                HStack(spacing: 20) {
                                    HStack(spacing: 0) {
                                        Text("\(rjview.count)").font(.system(size: 12.0)).fontWeight(av ? .bold : .regular).foregroundColor(.green)
                                        Text("Available").font(.system(size: 8.0)).fontWeight(av ? .bold : .regular).foregroundColor(.green)
                                    }
                                    HStack(spacing: 0) {
                                        Text("\(tlds.count - rjview.count)").font(.system(size: 12.0)).fontWeight(av ? .regular : .bold).foregroundColor(.red)
                                        Text("Locked").font(.system(size: 8.0)).fontWeight(av ? .regular : .bold).foregroundColor(.red)
                                    }
                                }
                            })
                        }
                    }
                    Spacer()
                    if !cartopen {
                        Button(action: { withAnimation { cartopen = true }}, label: {
                            HStack {
                                Text("\(cart.count)").font(.system(size: 12.0)).fontWeight(.light).foregroundColor(.white)
                                Image(systemName: "arrow.right").font(.system(size:12.0)).foregroundColor(.green.opacity(0.8))
                                Text("       ")
                            }
                        }).opacity(cart.isEmpty ? 0.0 : 1.0)
                    } else {
                        Text("           ")
                    }
                }
                // Domains Header
                
                // Domains Search
                if !cartopen {
                    HStack {
                        Button(action:{isFocused = true}, label: {
                            Image(systemName: "magnifyingglass").font(.system(size:20.0)).foregroundColor((name.isEmpty ? Color.white : Color.gray).opacity(0.8))
                        })
                        ZStack(alignment: .topLeading) {
                            if name.isEmpty {
                                Text(dplace.prefix(count)).font(.system(size:17.0)).tracking(0.8).foregroundColor(.gray.opacity(0.6)).padding(.top, 2).onReceive(qtimer) { _ in tfield()}}
                                TextField("", text: $name).font(.system(size:20.0)).foregroundColor(.white.opacity(0.88)).disableAutocorrection(true).autocapitalization(.none).focused($isFocused).onChange(of: name) { word in typing(word: word)}
                            }
                            Spacer()
                            if !name.isEmpty {
                                HStack(spacing: 20) {
                                    Button(action: {azsort()}, label: {
                                        HStack(spacing: 2) {
                                            Text(az || abcnone ? "a-z" : "z-a").font(.system(size:15.0)).foregroundColor(.white).opacity((az || za) && !abcnone ? 0.9 : 0.5).scaleEffect((az || za) && !abcnone ? 1.1 : 1.0)
                                        }
                                    }).opacity(name.isEmpty ? 0.0 : 1.0)
                                    
                                    Button(action: {psort()}, label: {
                                        HStack(spacing: 2) {
                                            Text("$").font(.system(size:15.0)).foregroundColor(.yellow).opacity(ph || pl ? 0.9 : 0.5).scaleEffect(ph || pl ? 1.1 : 1.0)
                                            VStack(spacing: 2) {
                                                Image(systemName: ph ? "chevron.down" : "chevron.up").font(.system(size:6.0)).foregroundColor(ph || pl ? .green : .white).opacity(ph || pl ? 0.9 : 0.5)
                                                Image(systemName: pl ? "chevron.up" : "chevron.down").font(.system(size:6.0)).foregroundColor(pl || ph ? .green : .white).opacity(ph || pl ? 0.9 : 0.5)
                                            }
                                        }
                                    }).opacity(name.isEmpty ? 0.0 : 1.0)
                                }
                            }
                    }.padding().background(.black.opacity(isFocused ? 0.6 : 0.5)).cornerRadius(10).padding(.bottom, -15).padding(10)
                }
                // Domains Search
                
                // Domains Cart
                if cartopen {
                    if cart.isEmpty {
                        HStack {
                            Text("No Items in Cart").font(.system(size:14.0)).fontWeight(.light).foregroundColor(.gray)
                        }
                    } else {
                        VStack(spacing: 0) {
                            HStack(spacing: 0) {
                                Text(" ")
                                Image(systemName: "globe").font(.system(size: 14.0)).foregroundColor(.blue)
                                Text("  Domain name").font(.system(size:12.0)).fontWeight(.light).foregroundColor(.white)
                                Spacer()
                                Text("$").font(.system(size:12.0)).foregroundColor(.yellow)
                                Text(" Yearly price   ").font(.system(size:12.0)).fontWeight(.light).foregroundColor(.white)
                            }
                            ScrollView {
                                VStack(spacing: 0) {
                                    ForEach(cart, id :\.self) { rj in
                                        Divider().padding(1)
                                        HStack(spacing: 0) {
                                            Text(" ")
                                            Image(systemName: "checkmark.circle.fill").font(.system(size: 8.0)).foregroundColor(.green).opacity(0.8)
                                            Text("**\(rj.domain).\(rj.tld)**").font(.system(size: 18.0)).foregroundColor(.white.opacity(0.9)).padding(10)
                                            Spacer()
                                            ZStack(alignment: .leading) {
                                                Text("$").font(.system(size: 10.0)).fontWeight(.light).foregroundColor(.white).padding(.bottom, 9)
                                                Text("  **\(rj.price)**").font(.system(size: 16.0)).fontWeight(.semibold).foregroundColor(.white)
                                            }
                                            Text("    ")
                                        }.padding(6).background(.black.opacity(0.5)).cornerRadius(10).contentShape(Rectangle())
                                    }
                                }
                                
                                // Stripe Card Element
                                StripeCard(p: $p).background(.black.opacity(0.5)).cornerRadius(10).contentShape(Rectangle())
                                // Stripe Card Element
                                
                                // Purchase Button
                                if done {
                                    HStack {
                                        Button(action: { payfirst() }, label: {
                                            HStack(spacing: 0) {
                                                Text("**\(cart.count)**").font(.system(size: 14.0)).fontWeight(.semibold).foregroundColor(.white)
                                                Text("\(cart.count == 1 ? " Domain" : " Domains")").font(.system(size: 12.0)).fontWeight(.light).foregroundColor(.white)
                                                Spacer()
                                                Text("\(purchasing ? "Purchasing" : "Purchase") ").font(.system(size: 14.0)).fontWeight(.semibold).foregroundColor(.white)
                                                if purchasing {
                                                    ProgressView()
                                                } else {
                                                    Image(systemName: "arrow.right").font(.system(size:10.0)).foregroundColor(.white)
                                                }
                                                Spacer()
                                                VStack(spacing: 0) {
                                                    Text("Total ").font(.system(size: 12.0)).fontWeight(.light).foregroundColor(.white)
                                                    ZStack(alignment: .leading) {
                                                        Text("$").font(.system(size: 10.0)).fontWeight(.light).foregroundColor(.white).padding(.bottom, 9)
                                                        Text("  **1000**").font(.system(size: 16.0)).fontWeight(.semibold).foregroundColor(.white)
                                                    }
                                                }
                                            }.padding().background(.green.opacity(0.5)).cornerRadius(15.0).overlay(RoundedRectangle(cornerRadius: 15).stroke(Color.green.opacity(0.9), lineWidth: 1))
                                        })
                                    }.padding()
                                }
                                // Purchase Button
                            }
                            Spacer()
                        }.padding()
                    }
                }
                // Domains Cart
                
                // Domains Results
                if qnow && !cartopen {
                    if rjview.isEmpty && !name.isEmpty {
                        VStack {
                            ProgressView().tint(Color.white.opacity(0.8)).padding()
                            Spacer()
                        }
                    } else {
                        if rjdview.isEmpty {
                            VStack(spacing: 0) {
                                ScrollView {
                                    VStack(spacing: 0) {
                                        ForEach(rjview, id: \.self) { rj in
                                            Divider().padding(1)
                                            HStack(spacing: 0) {
                                                Text(" ")
                                                Image(systemName: cart.contains(rj) ? "checkmark.circle.fill" : "checkmark").font(.system(size: 8.0)).foregroundColor(.green).opacity(0.8)
                                                Text("\(rj.domain)**.\(rj.tld)**").font(.system(size: 18.0)).foregroundColor(.white.opacity(0.9)).padding(10)
                                                Spacer()
                                                Text("\(rj.price)").font(.system(size: 16.0)).fontWeight(.light).foregroundColor((cart.contains(rj) ? Color.white : Color.gray).opacity(0.9)).tracking(1.0)
                                                Text(" ")
                                            }.padding(6).background(.black.opacity(0.5)).cornerRadius(10).contentShape(Rectangle()).onTapGesture{
                                                addtocart(rj : rj)
                                            }
                                        }
                                    }.padding()
                                }
                            }
                        } else {
                            VStack(spacing: 0) {
                                ScrollView {
                                    VStack(spacing: 0) {
                                        ForEach(rjdview, id: \.self) { rj in
                                            Divider().padding(1)
                                            HStack(spacing: 0) {
                                                Image(systemName: "xmark.circle.fill").font(.system(size: 8.0)).foregroundColor(.red).opacity(0.8)
                                                Text("\(rj.domain)**.\(rj.tld)**").font(.system(size: 18.0)).foregroundColor(.white.opacity(0.9)).padding(10)
                                                Spacer()
                                                Image(systemName: "globe.badge.chevron.backward").font(.system(size: 18.0)).foregroundColor(.gray)
                                            }.padding(6).background(.black.opacity(0.5)).cornerRadius(10).contentShape(Rectangle()).onTapGesture{ transfer() }
                                        }
                                    }.padding()
                                }
                            }
                        }
                        
                    }
                }
                // Domains Results
                
                Spacer()
                
            }.padding()
            // Domains Wrapper
            
        }.highPriorityGesture(DragGesture().onEnded({self.handleSwipe(translation: $0.translation.width)}))
        // Domains ZStack
    }
    
    private func handleSwipe(translation: CGFloat) {
        if translation > mind && !cartopen { back() }
        if translation > mind && cartopen {
            withAnimation { cartopen = false }
        } else if translation < -mind && !cartopen {
            withAnimation{ cartopen = true }
        }
    }
}

struct Domains_Previews: PreviewProvider {
    @State static var s: Bool = false
    @State static var t: [Rjs] = []
    static var previews: some View {
        Domains(showdomains: $s, tlds: $t)
    }
}
