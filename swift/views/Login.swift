import SwiftUI
import FirebaseAuth

struct Login: View {
    @StateObject var gauth = Gauth()
    var body: some View {
        if gauth.loggedin {
            ZStack {
                Image("tree").resizable().ignoresSafeArea()
                VStack(alignment: .leading, spacing: 0) {
                    HStack {
                        VStack(alignment: .leading, spacing: 0) {
                            Text("Welcome.").font(.largeTitle.bold()).foregroundColor(.white.opacity(0.9))
                            Text("Sign in/up.").font(.system(size: 20.0)).fontWeight(.medium).foregroundColor(.white.opacity(0.9))
                        }
                        .padding()
                        Spacer()
                    }
                    HStack {
                        Button(action: { gauth.signIn() }, label: {
                            Image("google").resizable().frame(width: 40.0, height: 40.0).shadow(color: Color.black.opacity(0.8), radius: 5, x: 0, y: 4)
                        }).padding(10).clipShape(Circle()).shadow(color: .white.opacity(0.3), radius: 10, x: 0, y: 1)
                    }.padding()
                }
            }
        } else {
            Settings().environmentObject(gauth)
        }
    }
}

struct Login_Previews: PreviewProvider {
    static var previews: some View {
        Login()
    }
}
