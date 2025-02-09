// import { Component, OnInit } from '@angular/core';
// import { ChatService } from '../services/chat.service';

// @Component({
//   standalone: true,
//   selector: 'app-chat-room',
//   imports: [],
//   templateUrl: './chat-room.component.html',
//   styleUrls: ['./chat-room.component.css']
// })
// export class ChatRoomComponent implements OnInit {
//   messages: any[] = [];
//   newMessage: string = '';
//   room: string = 'general';

//   constructor(private chatService: ChatService) {}

//   ngOnInit(): void {
//     this.chatService.getMessages(this.room).subscribe(messages => {
//       this.messages = messages;
//     });
//   }

//   sendMessage() {
//     if (this.newMessage.trim()) {
//       this.chatService.sendMessage(this.room, this.newMessage).subscribe(message => {
//         this.messages.push(message);
//         this.newMessage = '';
//       });
//     }
//   }
// }
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ChatService } from '../services/chat.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Message } from '../models/message.model';

@Component({
  standalone: true,
  imports: [CommonModule, FormsModule],
  selector: 'app-chat-room',
  templateUrl: './chat-room.component.html',
})
export class ChatRoomComponent implements OnInit {
  room = 'general';  // Example room
  messages: Message[] = [];
  newMessage: string = '';

  constructor(private chatService: ChatService) {}

  ngOnInit(): void {
    this.chatService.connect(this.room).subscribe((message: Message) => {
      this.messages.push(message);  // Add received message to the list
    });
  }

  sendMessage(): void {
    const message: Message = {
      content: this.newMessage,
      receiver: 'user2',  // Hardcoded for now, replace with dynamic logic
      sender: 'user1',    // Hardcoded for now, replace with dynamic logic
      room: this.room,
      timestamp: new Date().toISOString()
    };

    this.chatService.sendMessage(message);  // Send message through WebSocket
    this.newMessage = '';  // Clear input after sending
  }
}