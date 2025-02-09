// import { Injectable } from '@angular/core';
// import { HttpClient } from '@angular/common/http';
// import { Observable } from 'rxjs';

// @Injectable({
//   providedIn: 'root'
// })
// export class ChatService {
//   private apiUrl = 'http://localhost:8080';

//   constructor(private http: HttpClient) {}

//   sendMessage(room: string, content: string): Observable<any> {
//     return this.http.post(`${this.apiUrl}/messages`, { room, content });
//   }

//   getMessages(room: string): Observable<any[]> {
//     return this.http.get<any[]>(`${this.apiUrl}/messages?room=${room}`);
//   }
// }

import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { Message } from '../models/message.model';  // Import Message model

@Injectable({
  providedIn: 'root',
})
export class ChatService {
  private socket: WebSocket | null = null;
  private messages = new Subject<Message>();  // Change the type of messages

  connect(room: string): Observable<Message> {
    // Open WebSocket connection with the room as a query parameter
    this.socket = new WebSocket(`ws://localhost:8080/ws?room=${room}`);

    // On receiving a message, parse it as a JSON object
    this.socket.onmessage = (event) => {
      const message: Message = JSON.parse(event.data);  // Parse the message data into Message model
      this.messages.next(message);  // Emit the message to subscribers
    };

    // Return the observable to be consumed by the components
    return this.messages.asObservable();
  }

  sendMessage(message: Message) {
    // Send the message as a JSON string
    if (this.socket?.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));  // Convert message to JSON string
    } else {
      console.error('WebSocket is not open.');
    }
  }

  close() {
    // Close WebSocket connection when done
    this.socket?.close();
  }
}
