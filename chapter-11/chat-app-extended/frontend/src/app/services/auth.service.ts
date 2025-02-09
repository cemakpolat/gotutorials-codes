import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private apiUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) {}

  login(username: string, password: string): Observable<string> {
    return this.http.post<string>(`${this.apiUrl}/login`, { username, password });
  }

  register(username: string, password: string): Observable<void> {
    return this.http.post<void>(`${this.apiUrl}/register`, { username, password });
  }
}
