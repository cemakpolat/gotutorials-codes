import { Routes } from '@angular/router';
import { AuthComponent } from './auth/auth.component';
import { ChatRoomComponent } from './chat-room/chat-room.component';
import { RoomListComponent } from './room-list/room-list.component';

export const routes: Routes = [
  { path: '', component: AuthComponent },
  { path: 'rooms', component: RoomListComponent },
  { path: 'chat/:room', component: ChatRoomComponent },
];
