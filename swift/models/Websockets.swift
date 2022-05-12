import Foundation
import Combine

class Wcon: ConDelegate {
    var con: Con!
    init() {
        self.con = WebSocketTaskConnection(url: URL(string: "ws://localhost:3000/ws?q=test")!)
        self.con.delegate = self
        self.con.connect()
    }
    func onConnected(connection: Con) {
        print("Connected")
    }
        
    func onDisconnected(connection: Con, error: Error?) {
        if let error = error {
            print("Disconnected with error:\(error)")
        } else {
            print("Disconnected normally")
        }
    }
        
    func onError(connection: Con, error: Error) {
        print("Connection error:\(error)")
    }
    
    func onMessage(connection: Con, text: String) {
        print("Text message: \(text)")
    }
    
    func onMessage(connection: Con, data: Data) {
        print("Data message: \(data)")
    }
}

protocol Con {
    func send(text: String)
    func send(data: Data)
    func connect()
    func disconnect()
    var delegate: ConDelegate? {
        get
        set
    }
}

protocol ConDelegate {
    func onConnected(connection: Con)
    func onDisconnected(connection: Con, error: Error?)
    func onError(connection: Con, error: Error)
    func onMessage(connection: Con, text: String)
    func onMessage(connection: Con, data: Data)
}

class WebSocketTaskConnection: NSObject, Con, URLSessionWebSocketDelegate {
    var delegate: ConDelegate?
    var webSocketTask: URLSessionWebSocketTask!
    var urlSession: URLSession!
    let delegateQueue = OperationQueue()
    
    init(url: URL) {
        super.init()
        urlSession = URLSession(configuration: .default, delegate: self, delegateQueue: delegateQueue)
        webSocketTask = urlSession.webSocketTask(with: url)
    }
    
    func urlSession(_ session: URLSession, webSocketTask: URLSessionWebSocketTask, didOpenWithProtocol protocol: String?) {
        self.delegate?.onConnected(connection: self)
    }
    
    func urlSession(_ session: URLSession, webSocketTask: URLSessionWebSocketTask, didCloseWith closeCode: URLSessionWebSocketTask.CloseCode, reason: Data?) {
        self.delegate?.onDisconnected(connection: self, error: nil)
    }
    
    func connect() {
       webSocketTask.resume()
       listen()
   }
    
    func disconnect() {
        webSocketTask.cancel(with: .goingAway, reason: nil)
    }
    
    func listen()  {
        webSocketTask.receive { result in
            switch result {
            case .failure(let error):
                self.delegate?.onError(connection: self, error: error)
            case .success(let message):
                switch message {
                case .string(let text):
                    self.delegate?.onMessage(connection: self, text: text)
                case .data(let data):
                    self.delegate?.onMessage(connection: self, data: data)
                @unknown default:
                    fatalError()
                }
                
                self.listen()
            }
        }
    }
    
    func send(text: String) {
        webSocketTask.send(URLSessionWebSocketTask.Message.string(text)) { error in
            if let error = error {
                self.delegate?.onError(connection: self, error: error)
            }
        }
    }
    
    func send(data: Data) {
        webSocketTask.send(URLSessionWebSocketTask.Message.data(data)) { error in
            if let error = error {
                self.delegate?.onError(connection: self, error: error)
            }
        }
    }
}
