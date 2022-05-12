import SwiftUI
import FirebaseAuth
import Alamofire

struct Settings: View {
    @State var showdomains = false
    @EnvironmentObject var gauth: Gauth
    @State var tlds: [Rjs] = []
    @State var showprods = false
    @State var recent = ""
    let mind: CGFloat = 50
    func signout() {
        let firebaseAuth = Auth.auth()
        do {
            try firebaseAuth.signOut()
        } catch let signOutError as NSError {
            print("Error signing out: %@", signOutError)
            return
        }
        gauth.change()
    }
    
    func loaddomains() {
        print("loading domains")
        AF.request("http://localhost:3000/tlds").responseDecodable(of: [Rjs].self) { res in
            dump(res.error)
            tlds = res.value ?? []
            dump(tlds)
            print(tlds.count)
        }
    }
    
    func loader() {
        let u = Auth.auth().currentUser
        let uid = u?.uid ?? ""
        let email = u?.email ?? ""
        let name = u?.displayName ?? ""
        let num = u?.phoneNumber ?? ""
        let did = UIDevice.current.identifierForVendor!.uuidString
        let params: [String: [String]] = [
            "uid": [uid],
            "email": [email],
            "name": [name],
            "number": [num],
            "device": [did]
        ]
        AF.request("http://localhost:3000/login", parameters: params, encoder: JSONParameterEncoder.default).responseDecodable(of: [Rjs].self) { res in
            dump(res.error)
        }
    }
    
    var body: some View {
        if showdomains {
            Domains(showdomains: $showdomains, tlds: $tlds).onAppear{
                loaddomains()
                recent = "domains"
            }
        }
        if showprods {
            Text("test")
        }
        if !showdomains && !showprods {
            ZStack {
                Image("settings").resizable().ignoresSafeArea()
                
                VStack {
                    HStack(spacing: 0) {
                        Button(action: { withAnimation{ showdomains = true }}, label: {
                                Image(systemName: "globe").font(.system(size: 25.0)).foregroundColor(.gray).opacity(0.2)
                                Text("Domains ").font(.system(size:15.0)).fontWeight(.bold).foregroundColor(.white).opacity(0.7)
                                Spacer()
                                Image(systemName: "arrow.right").font(.system(size: 12.0)).foregroundColor(.green.opacity(0.6))
                        })
                        Spacer()
                    }
                    .padding()
                    HStack(spacing: 0) {
                        Button(action: {
                            withAnimation{ showprods = true }
                            loader()
                        }, label: {
                                Image(systemName: "doc").font(.system(size: 25.0)).foregroundColor(.gray).opacity(0.2)
                                Text("Products ").font(.system(size:15.0)).fontWeight(.bold).foregroundColor(.white).opacity(0.7)
                                Spacer()
                                Image(systemName: "arrow.right").font(.system(size: 12.0)).foregroundColor(.green.opacity(0.6))
                        })
                        Spacer()
                    }
                    .padding()
                    Spacer()
                    HStack {
                        Image(systemName: "person").font(.system(size: 25.0)).foregroundColor(.white).opacity(0.2)
                        Text("\(Auth.auth().currentUser?.displayName ?? "")").font(.system(size: 15.0)).foregroundColor(.white)
                        Spacer()
                        Button(action: { signout() }, label: {
                            HStack(spacing: 2) {
                                Text("Sign out?").font(.system(size:14.0)).foregroundColor(.white.opacity(0.2))
                                Image(systemName: "xmark").font(.system(size:14.0)).foregroundColor(.red.opacity(0.8)).padding()
                            }
                        })
                    }
                }
                .padding()
            }.highPriorityGesture(DragGesture().onEnded({self.handleSwipe(translation: $0.translation.width)}))
            // Settings ZStack
        }
    }
    
    private func handleSwipe(translation: CGFloat) {
        if translation < -mind {
            print(recent)
            if recent == "" { return }
            if recent == "prods" {
                withAnimation{ showprods = true }
            }
            if recent == "domains" {
                withAnimation{ showdomains = true }
            }

        }
    }
}

struct Settings_Previews: PreviewProvider {
    static var previews: some View {
        Settings()
    }
}
