import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  standalone: true,
  imports: [CommonModule, FormsModule],
  selector: 'app-auth',
  templateUrl: './auth.component.html',
})
export class AuthComponent {
  username = 'testuser';
  password = 'password123';
  isRegister = false;

  constructor(private authService: AuthService, private router: Router) {}

  submit() {
    if (this.isRegister) {
      this.authService.register(this.username, this.password).subscribe(() => {
        alert('Registration successful');
        this.isRegister = false;
      });
    } else {
      this.authService.login(this.username, this.password).subscribe(
        (token) => {
          localStorage.setItem('token', token);
          this.router.navigate(['/rooms']);
        },
        () => alert('Login failed')
      );
    }
  }
}
