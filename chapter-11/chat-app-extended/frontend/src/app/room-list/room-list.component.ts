import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  standalone: true,
  imports: [CommonModule],
  selector: 'app-room-list',
  templateUrl: './room-list.component.html',
})
export class RoomListComponent {
  rooms = ['General', 'Tech', 'Sports'];

  constructor(private router: Router) {}

  joinRoom(room: string) {
    this.router.navigate([`/chat/${room}`]);
  }
}
